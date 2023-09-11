package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"diagonal.works/b6"
	"diagonal.works/b6/api"
	"diagonal.works/b6/api/functions"
	"diagonal.works/b6/geojson"
	"diagonal.works/b6/ingest"
	pb "diagonal.works/b6/proto"
	"diagonal.works/b6/renderer"
	"github.com/golang/geo/s2"
	"google.golang.org/protobuf/encoding/protojson"
)

func getStringExpression(f b6.Feature, key string) *pb.NodeProto {
	return &pb.NodeProto{
		Node: &pb.NodeProto_Call{
			Call: &pb.CallNodeProto{
				Function: &pb.NodeProto{
					Node: &pb.NodeProto_Symbol{
						Symbol: "get-string",
					},
				},
				Args: []*pb.NodeProto{
					&pb.NodeProto{
						Node: &pb.NodeProto_Literal{
							Literal: &pb.LiteralNodeProto{
								Value: &pb.LiteralNodeProto_FeatureIDValue{
									FeatureIDValue: b6.NewProtoFromFeatureID(f.FeatureID()),
								},
							},
						},
					},
					&pb.NodeProto{
						Node: &pb.NodeProto_Literal{
							Literal: &pb.LiteralNodeProto{
								Value: &pb.LiteralNodeProto_StringValue{
									StringValue: key,
								},
							},
						},
					},
				},
			},
		},
	}
}

const CollectionLineLimit = 200
const CollectionHighlightLimit = 10000

func fillKeyValues(c api.Collection, keys []interface{}, values []interface{}) ([]interface{}, []interface{}, error) {
	i := c.Begin()
	var err error
	for {
		var ok bool
		ok, err = i.Next()
		if !ok || err != nil {
			break
		}
		keys = append(keys, i.Key())
		values = append(values, i.Value())
		if len(keys) >= CollectionHighlightLimit {
			break
		}
	}
	return keys, values, err
}

func isFeatureCollection(keys []interface{}, values []interface{}) bool {
	if len(keys) > 0 {
		if id, ok := keys[0].(b6.FeatureID); ok {
			if f, ok := values[0].(b6.Feature); ok {
				return id == f.FeatureID()
			}
		}
	}
	return false
}

func isArrayCollection(keys []interface{}, values []interface{}) bool {
	if len(keys) > 0 {
		for i, k := range keys {
			if ii, ok := k.(int); !ok || i != ii {
				return false
			}
		}
	} else {
		return false
	}
	return true
}

func fillMatchingFunctionSymbols(symbols []string, result interface{}, functions api.FunctionSymbols) []string {
	t := reflect.TypeOf(result)
	for symbol, f := range functions {
		tt := reflect.TypeOf(f)
		if tt.Kind() == reflect.Func && tt.NumIn() > 1 {
			if api.CanUseAsArg(t, tt.In(1)) {
				symbols = append(symbols, symbol)
			}
		}
	}
	return symbols
}

func NewUIHandler(w ingest.MutableWorld, cores int) *UIHandler {
	local := make(api.FunctionSymbols)
	for name, f := range functions.Functions() {
		local[name] = f
	}
	return &UIHandler{
		World:            w,
		RenderRules:      renderer.BasemapRenderRules,
		Cores:            cores,
		FunctionSymbols:  local,
		FunctionWrappers: functions.Wrappers(),
	}
}

type UIResponseProtoJSON pb.UIResponseProto

func (b *UIResponseProtoJSON) MarshalJSON() ([]byte, error) {
	return protojson.Marshal((*pb.UIResponseProto)(b))
}

func (b *UIResponseProtoJSON) UnmarshalJSON(buffer []byte) error {
	return protojson.Unmarshal(buffer, (*pb.UIResponseProto)(b))
}

type UIResponseJSON struct {
	Proto   *UIResponseProtoJSON `json:"proto,omitempty"`
	GeoJSON geojson.GeoJSON      `json:"geojson,omitempty"`
}

type UIHandler struct {
	World            ingest.MutableWorld
	RenderRules      renderer.RenderRules
	Cores            int
	FunctionSymbols  api.FunctionSymbols
	FunctionWrappers api.FunctionWrappers
}

func (b *UIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := &pb.UIRequestProto{}
	response := &UIResponseJSON{
		Proto: &UIResponseProtoJSON{
			Stack: &pb.StackProto{},
		},
	}

	if r.Method == "GET" {
		request.Expression = r.URL.Query().Get("e")
	} else if r.Method == "POST" {
		var err error
		var body []byte
		if body, err = io.ReadAll(r.Body); err == nil {
			r.Body.Close()
			err = protojson.Unmarshal(body, request)
		}
		if err != nil {
			http.Error(w, "Bad request body", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return
	}

	if request.Expression == "" && request.Node == nil {
		http.Error(w, "No expression", http.StatusBadRequest)
		return
	}

	if request.Expression != "" {
		var err error
		if request.Node == nil {
			response.Proto.Node, err = api.ParseExpression(request.Expression)
		} else {
			response.Proto.Node, err = api.ParseExpressionWithLHS(request.Expression, request.Node)
		}
		if err != nil {
			response.Proto.Stack.Substacks = fillSubstacksFromError(response.Proto.Stack.Substacks, err)
			sendUIResponse(response, w)
			return
		}
	} else {
		response.Proto.Node = request.Node
	}
	response.Proto.Node = api.Simplify(response.Proto.Node, b.FunctionSymbols)

	substack := &pb.SubstackProto{}
	fillSubstackFromExpression(substack, response.Proto.Node, true)
	if len(substack.Lines) > 0 {
		response.Proto.Stack.Substacks = append(response.Proto.Stack.Substacks, substack)
	}

	context := api.Context{
		World:            b.World,
		FunctionSymbols:  b.FunctionSymbols,
		FunctionWrappers: b.FunctionWrappers,
		Cores:            b.Cores,
		Context:          context.Background(),
	}
	result, err := api.Evaluate(response.Proto.Node, &context)
	if err == nil {
		if err = fillResponseFromResult(response, result, b.RenderRules, b.World); err == nil {
			shell := &pb.ShellLineProto{
				Functions: make([]string, 0),
			}
			shell.Functions = fillMatchingFunctionSymbols(shell.Functions, result, b.FunctionSymbols)
			response.Proto.Stack.Substacks = append(response.Proto.Stack.Substacks, &pb.SubstackProto{
				Lines: []*pb.LineProto{{Line: &pb.LineProto_Shell{Shell: shell}}},
			})
		}
	} else {
		response.Proto.Stack.Substacks = fillSubstacksFromError(response.Proto.Stack.Substacks, err)
	}
	sendUIResponse(response, w)
}

func sendUIResponse(response *UIResponseJSON, w http.ResponseWriter) {
	output, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(output))
}

func highlightInResponse(response *pb.UIResponseProto, id b6.FeatureID) {
	if response.Highlighted == nil {
		response.Highlighted = &pb.FeatureIDsProto{}
	}
	n := fmt.Sprintf("/%s/%s", id.Type.String(), id.Namespace.String())
	ids := -1
	for i, nn := range response.Highlighted.Namespaces {
		if n == nn {
			ids = i
			break
		}
	}
	if ids < 0 {
		ids = len(response.Highlighted.Ids)
		response.Highlighted.Namespaces = append(response.Highlighted.Namespaces, n)
		response.Highlighted.Ids = append(response.Highlighted.Ids, &pb.IDsProto{})
	}
	response.Highlighted.Ids[ids].Ids = append(response.Highlighted.Ids[ids].Ids, id.Value)
}

func fillSubstackFromExpression(lines *pb.SubstackProto, expression *pb.NodeProto, root bool) {
	if call, ok := expression.Node.(*pb.NodeProto_Call); ok {
		if call.Call.Pipelined {
			left := call.Call.Args[0]
			right := &pb.NodeProto{
				Node: &pb.NodeProto_Call{
					Call: &pb.CallNodeProto{
						Function: call.Call.Function,
						Args:     call.Call.Args[1:],
					},
				},
			}
			fillSubstackFromExpression(lines, left, false)
			fillSubstackFromExpression(lines, right, false)
			return
		}
	}
	_, isLiteral := expression.Node.(*pb.NodeProto_Literal)
	if expression, ok := api.UnparseNode(expression); ok {
		if !isLiteral || !root {
			lines.Lines = append(lines.Lines, &pb.LineProto{
				Line: &pb.LineProto_Expression{
					Expression: &pb.ExpressionLineProto{
						Expression: expression,
					},
				},
			})
		}
	} else {
		lines.Lines = append(lines.Lines, &pb.LineProto{
			Line: &pb.LineProto_Error{
				Error: &pb.ErrorLineProto{
					Error: "can't convert function",
				},
			},
		})
	}
}

func fillSubstackFromCollection(substack *pb.SubstackProto, c api.Collection, response *pb.UIResponseProto) error {
	// TODO: Set collection title based on collection contents
	if countable, ok := c.(api.Countable); ok {
		substack.Lines = append(substack.Lines, &pb.LineProto{
			Line: &pb.LineProto_ValuePair{
				ValuePair: &pb.ValuePairLineProto{
					First: &pb.ClickableAtomProto{
						Atom: atomFromString("Collection"),
					},
					Second: &pb.ClickableAtomProto{
						Atom: atomFromString(strconv.Itoa(countable.Count())),
					},
				},
			},
		})
	} else {
		substack.Lines = append(substack.Lines, &pb.LineProto{
			Line: &pb.LineProto_Value{
				Value: &pb.ValueLineProto{
					Atom: atomFromString("Collection"),
				},
			},
		})
	}

	keys := make([]interface{}, 0, 8)
	values := make([]interface{}, 0, 8)
	var err error
	keys, values, err = fillKeyValues(c, keys, values)
	if err != nil {
		return err
	}
	if isFeatureCollection(keys, values) || isArrayCollection(keys, values) {
		for i := range values {
			if i < CollectionLineLimit {
				line := valueLineFromValue(values[i])
				substack.Lines = append(substack.Lines, line)
			} else {
				break
			}
		}
	} else {
		for i := range keys {
			if i < CollectionLineLimit {
				line := valuePairLineFromValues(keys[i], values[i])
				substack.Lines = append(substack.Lines, line)
			} else {
				break
			}
		}
	}

	for i := range keys {
		if id, ok := keys[i].(b6.Identifiable); ok {
			highlightInResponse(response, id.FeatureID())
		}
		if id, ok := values[i].(b6.Identifiable); ok {
			highlightInResponse(response, id.FeatureID())
		}
	}
	return nil
}

func fillSubstackFromHistogram(substack *pb.SubstackProto, c *api.HistogramCollection) error {
	keys, values, err := fillKeyValues(c, nil, nil)
	if err != nil {
		return err
	}

	total := 0
	begin := len(substack.Lines)
	for i, key := range keys {
		// TODO: Factor out w/ collection-key-or-value fill when we rework blocks
		var bucketRange string
		switch v := key.(type) {
		case int:
			bucketRange = strconv.Itoa(v)
		case float64:
			bucketRange = fmt.Sprintf("%f", v)
		case string:
			bucketRange = v
		case b6.Tag:
			bucketRange = api.UnparseTag(v)
		default:
			bucketRange = fmt.Sprintf("%+v", v)
		}
		substack.Lines = append(substack.Lines, &pb.LineProto{
			Line: &pb.LineProto_HistogramBar{
				HistogramBar: &pb.HistogramBarLineProto{
					Range: bucketRange,
					Value: int32(values[i].(int)),
					Index: int32(i),
				},
			},
		})
		total += values[i].(int)
	}
	for i := begin; i < len(substack.Lines); i++ {
		substack.Lines[i].GetHistogramBar().Total = int32(total)
	}
	return nil
}

func lineFromTags(f b6.Feature) *pb.LineProto {
	tags := f.AllTags()
	tl := &pb.TagsLineProto{
		Tags: make([]*pb.TagAtomProto, len(tags)),
	}
	for i, tag := range tags {
		if strings.HasPrefix(tag.Key, "#") || strings.HasPrefix(tag.Key, "#") {
			tl.Tags[i] = &pb.TagAtomProto{Prefix: tag.Key[0:1], Key: tag.Key[1:], Value: tag.Value}
		} else {
			tl.Tags[i] = &pb.TagAtomProto{Prefix: "", Key: tag.Key, Value: tag.Value}
		}
		tl.Tags[i].ClickExpression = getStringExpression(f, tag.Key)
	}
	return &pb.LineProto{
		Line: &pb.LineProto_Tags{
			Tags: tl,
		},
	}
}

func atomFromString(value string) *pb.AtomProto {
	return &pb.AtomProto{
		Atom: &pb.AtomProto_Value{
			Value: value,
		},
	}
}

func featureLabel(f b6.Feature) string {
	if name := f.Get("name"); name.IsValid() {
		return name.Value
	} else if code := f.Get("code"); code.IsValid() {
		return code.Value
	} else if ref := f.Get("ref"); ref.IsValid() {
		return ref.Value
	} else {
		switch f.FeatureID().Namespace {
		case b6.NamespaceGBCodePoint:
			if postcode, ok := b6.PostcodeFromPointID(f.FeatureID().ToPointID()); ok {
				return postcode
			}
		case b6.NamespaceUKONSBoundaries:
			if code, _, ok := b6.UKONSCodeFromFeatureID(f.FeatureID()); ok {
				return code
			}
		}
	}
	return LabelForFeature(f).Singular
}

func atomFromValue(value interface{}) *pb.AtomProto {
	if i, ok := api.ToInt(value); ok {
		return atomFromString(strconv.Itoa(i))
	} else if f, err := api.ToFloat64(value); err == nil {
		return atomFromString(fmt.Sprintf("%f", f))
	} else {
		switch v := value.(type) {
		case string:
			return atomFromString(v)
		case b6.Feature:
			return &pb.AtomProto{
				Atom: &pb.AtomProto_LabelledIcon{
					LabelledIcon: &pb.LabelledIconProto{
						Icon:  v.FeatureID().Type.String(),
						Label: featureLabel(v),
					},
				},
			}
		case b6.FeatureID:
			return &pb.AtomProto{
				Atom: &pb.AtomProto_LabelledIcon{
					LabelledIcon: &pb.LabelledIconProto{
						Icon:  v.Type.String(),
						Label: strings.Title(v.Type.String()),
					},
				},
			}
		case b6.Tag:
			return atomFromString(api.UnparseTag(v))
		case b6.Point:
			ll := s2.LatLngFromPoint(v.Point())
			return atomFromString(fmt.Sprintf("%f, %f", ll.Lat.Degrees(), ll.Lng.Degrees()))
		default:
			return atomFromString(fmt.Sprintf("%v", v))
		}
	}
}

func clickExpressionFromIdentifiable(f b6.Identifiable) *pb.NodeProto {
	return &pb.NodeProto{
		Node: &pb.NodeProto_Call{
			Call: &pb.CallNodeProto{
				Function: &pb.NodeProto{
					Node: &pb.NodeProto_Symbol{
						Symbol: "find-feature",
					},
				},
				Args: []*pb.NodeProto{
					&pb.NodeProto{
						Node: &pb.NodeProto_Literal{
							Literal: &pb.LiteralNodeProto{
								Value: &pb.LiteralNodeProto_FeatureIDValue{
									FeatureIDValue: b6.NewProtoFromFeatureID(f.FeatureID()),
								},
							},
						},
					},
				},
			},
		},
	}
}

func clickExpressionFromValue(value interface{}) *pb.NodeProto {
	switch v := value.(type) {
	case b6.Identifiable:
		return clickExpressionFromIdentifiable(v)
	}
	return nil
}

func valueLineFromValue(value interface{}) *pb.LineProto {
	return &pb.LineProto{
		Line: &pb.LineProto_Value{
			Value: &pb.ValueLineProto{
				Atom:            atomFromValue(value),
				ClickExpression: clickExpressionFromValue(value),
			},
		},
	}
}

func valuePairLineFromValues(first interface{}, second interface{}) *pb.LineProto {
	return &pb.LineProto{
		Line: &pb.LineProto_ValuePair{
			ValuePair: &pb.ValuePairLineProto{
				First: &pb.ClickableAtomProto{
					Atom:            atomFromValue(first),
					ClickExpression: clickExpressionFromValue(first),
				},
				Second: &pb.ClickableAtomProto{
					Atom:            atomFromValue(second),
					ClickExpression: clickExpressionFromValue(second),
				},
			},
		},
	}
}

func fillSubstacksFromAtom(substacks []*pb.SubstackProto, atom *pb.AtomProto) []*pb.SubstackProto {
	return append(substacks, &pb.SubstackProto{
		Lines: []*pb.LineProto{
			{
				Line: &pb.LineProto_Value{
					Value: &pb.ValueLineProto{
						Atom: atom,
					},
				},
			},
		},
	})
}

func fillSubstacksFromString(substacks []*pb.SubstackProto, value string) []*pb.SubstackProto {
	return fillSubstacksFromAtom(substacks, atomFromString(value))
}

func fillSubstacksFromError(substacks []*pb.SubstackProto, err error) []*pb.SubstackProto {
	return append(substacks, &pb.SubstackProto{
		Lines: []*pb.LineProto{
			{
				Line: &pb.LineProto_Error{
					Error: &pb.ErrorLineProto{
						Error: err.Error(),
					},
				},
			},
		},
	})
}

func fillSubstacksFromFeature(substacks []*pb.SubstackProto, f b6.Feature, w b6.World) []*pb.SubstackProto {
	substack := &pb.SubstackProto{}
	substack.Lines = append(substack.Lines, valueLineFromValue(f))
	substack.Lines = append(substack.Lines, lineFromTags(f))
	substacks = append(substacks, substack)
	if path, ok := f.(b6.PathFeature); ok {
		substack := &pb.SubstackProto{Collapsable: true}
		line := valuePairLineFromValues("Points", path.Len())
		substack.Lines = append(substack.Lines, line)
		for i := 0; i < path.Len(); i++ {
			if point := path.Feature(i); point != nil {
				substack.Lines = append(substack.Lines, valueLineFromValue(point))
			} else {
				substack.Lines = append(substack.Lines, valueLineFromValue(path.Point(i)))
			}
		}
		substacks = append(substacks, substack)
	}
	relations := b6.AllRelations(w.FindRelationsByFeature(f.FeatureID()))
	if len(relations) > 0 {
		substack := &pb.SubstackProto{Collapsable: true}
		line := valuePairLineFromValues("Relations", len(relations))
		substack.Lines = append(substack.Lines, line)
		for _, r := range relations {
			substack.Lines = append(substack.Lines, valueLineFromValue(r))
		}
		substacks = append(substacks, substack)
	}
	return substacks
}

func fillResponseFromResult(response *UIResponseJSON, result interface{}, rules renderer.RenderRules, w b6.World) error {
	p := (*pb.UIResponseProto)(response.Proto)
	switch r := result.(type) {
	case string:
		p.Stack.Substacks = fillSubstacksFromString(p.Stack.Substacks, r)
	case b6.Feature:
		p.Stack.Substacks = fillSubstacksFromFeature(p.Stack.Substacks, r, w)
		highlightInResponse(p, r.FeatureID())
	case b6.Query:
		if q, ok := api.UnparseQuery(r); ok {
			p.Stack.Substacks = fillSubstacksFromString(p.Stack.Substacks, q)
		} else {
			// TODO: Improve the rendering of queries
			p.Stack.Substacks = fillSubstacksFromString(p.Stack.Substacks, "Query")
		}
	case b6.Tag:
		p.Stack.Substacks = fillSubstacksFromAtom(p.Stack.Substacks, atomFromValue(r))
		if !rules.IsRendered(r) {
			if q, ok := api.UnparseQuery(b6.Tagged(r)); ok {
				p.QueryLayers = append(p.QueryLayers, q)
			}
		}
	case *api.HistogramCollection:
		substack := &pb.SubstackProto{}
		if err := fillSubstackFromHistogram(substack, r); err == nil {
			p.Stack.Substacks = append(p.Stack.Substacks, substack)
		} else {
			return err
		}
	case api.Collection:
		substack := &pb.SubstackProto{}
		if err := fillSubstackFromCollection(substack, r, p); err == nil {
			p.Stack.Substacks = append(p.Stack.Substacks, substack)
		} else {
			return err
		}
	case b6.Area:
		dimension := 0.0
		for i := 0; i < r.Len(); i++ {
			dimension += b6.AreaToMeters2(r.Polygon(i).Area())
		}
		atom := &pb.AtomProto{
			Atom: &pb.AtomProto_Download{
				Download: fmt.Sprintf("%.2fm² area", dimension),
			},
		}
		p.Stack.Substacks = fillSubstacksFromAtom(p.Stack.Substacks, atom)
		response.GeoJSON = r.ToGeoJSON()
	case b6.Path:
		dimension := b6.AngleToMeters(r.Polyline().Length())
		atom := &pb.AtomProto{
			Atom: &pb.AtomProto_Download{
				Download: fmt.Sprintf("%.2fm path", dimension),
			},
		}
		p.Stack.Substacks = fillSubstacksFromAtom(p.Stack.Substacks, atom)
		response.GeoJSON = r.ToGeoJSON()
	case *geojson.FeatureCollection:
		var label string
		if n := len(r.Features); n == 1 {
			label = "1 GeoJSON feature"
		} else {
			label = fmt.Sprintf("%d GeoJSON features", n)
		}
		atom := &pb.AtomProto{
			Atom: &pb.AtomProto_Download{
				Download: fmt.Sprintf(label),
			},
		}
		p.Stack.Substacks = fillSubstacksFromAtom(p.Stack.Substacks, atom)
		response.GeoJSON = r
	case *geojson.Feature:
		atom := &pb.AtomProto{
			Atom: &pb.AtomProto_Download{
				Download: "GeoJSON feature",
			},
		}
		p.Stack.Substacks = fillSubstacksFromAtom(p.Stack.Substacks, atom)
		response.GeoJSON = r
	case *geojson.Geometry:
		atom := &pb.AtomProto{
			Atom: &pb.AtomProto_Download{
				Download: "GeoJSON geometry",
			},
		}
		p.Stack.Substacks = fillSubstacksFromAtom(p.Stack.Substacks, atom)
		response.GeoJSON = geojson.NewFeatureWithGeometry(*r)
	default:
		substack := &pb.SubstackProto{
			Lines: []*pb.LineProto{{
				Line: &pb.LineProto_Value{
					Value: &pb.ValueLineProto{
						Atom: atomFromValue(r),
					},
				},
			}},
		}
		p.Stack.Substacks = append(p.Stack.Substacks, substack)
	}
	return nil
}