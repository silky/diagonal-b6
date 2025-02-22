package ui

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"diagonal.works/b6"
	"diagonal.works/b6/api"
	"diagonal.works/b6/api/functions"
	"diagonal.works/b6/geojson"
	"diagonal.works/b6/ingest"
	pb "diagonal.works/b6/proto"
	"diagonal.works/b6/renderer"
	"github.com/golang/geo/s2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// See https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Number/MAX_SAFE_INTEGER
const MaxSafeJavaScriptInteger = (1 << 53) - 1

type Options struct {
	StaticPath        string
	JavaScriptPath    string
	StaticV2Path      string
	StorybookPath     string
	EnableV2UI        bool
	EnableVite        bool
	EnableStorybook   bool
	BasemapRules      renderer.RenderRules
	UI                UI
	Worlds            ingest.Worlds
	APIOptions        api.Options
	InstrumentHandler func(handler http.Handler, name string) http.Handler
	Lock              *sync.RWMutex
}

type DropPrefixFilesystem struct {
	Prefix string
	Next   http.FileSystem
}

func (d *DropPrefixFilesystem) Open(filename string) (http.File, error) {
	if strings.HasPrefix(filename, d.Prefix) {
		return d.Next.Open(filename[len(d.Prefix):])
	}
	return nil, fs.ErrNotExist
}

type MergedFilesystem []string

func (m MergedFilesystem) Open(filename string) (http.File, error) {
	for _, path := range m {
		full := filepath.Join(path, filename)
		if _, err := os.Stat(full); err == nil {
			return os.Open(full)
		}
	}
	return nil, fs.ErrNotExist
}

func RegisterWebInterface(root *http.ServeMux, options *Options) error {
	staticPaths := strings.Split(options.StaticPath, ",")

	v1Path := "/"
	v2Path := "/v2.html"
	if options.EnableV2UI {
		v1Path = "/v1.html"
		v2Path = "/"
	}

	if len(staticPaths) > 0 {
		if v1Path == "/" {
			root.Handle("/", http.FileServer(MergedFilesystem(staticPaths)))
		} else {
			root.Handle(v1Path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, filepath.Join(staticPaths[0], "index.html"))
			}))
		}
	}

	root.Handle("/b6.css", http.FileServer(MergedFilesystem(staticPaths)))

	root.Handle("/bundle.js", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(options.JavaScriptPath, "bundle.js"))
	}))

	staticV2Paths := strings.Split(options.StaticV2Path, ",")
	root.Handle("/assets/", http.FileServer(MergedFilesystem(staticV2Paths)))
	if len(staticV2Paths) > 0 {
		if options.EnableVite {
			root.Handle(v2Path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, filepath.Join(staticV2Paths[0], "index-vite.html"))
			}))
		} else {
			root.Handle(v2Path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, filepath.Join(staticV2Paths[0], "index.html"))
			}))
		}
	}

	if options.EnableStorybook {
		storybookPaths := strings.Split(options.StorybookPath, ",")
		root.Handle("/storybook/", http.FileServer(&DropPrefixFilesystem{
			Prefix: "/storybook",
			Next:   MergedFilesystem(storybookPaths),
		}))
	}

	evaluator := api.Evaluator{
		Worlds:          options.Worlds,
		Options:         options.APIOptions,
		FunctionSymbols: functions.Functions(),
		Adaptors:        functions.Adaptors(),
		Lock:            options.Lock,
	}

	var ui UI
	if options.UI != nil {
		ui = options.UI
	} else {
		ui = &OpenSourceUI{
			Worlds:          options.Worlds,
			Evaluator:       evaluator,
			BasemapRules:    renderer.BasemapRenderRules,
			FunctionSymbols: functions.Functions(),
			Lock:            options.Lock,
		}
	}
	startup := http.Handler(&StartupHandler{UI: ui})
	if options.InstrumentHandler != nil {
		startup = options.InstrumentHandler(startup, "startup")
	}
	root.Handle("/startup", startup)
	stack := http.Handler(&StackHandler{UI: ui})
	if options.InstrumentHandler != nil {
		stack = options.InstrumentHandler(stack, "ui")
	}
	root.Handle("/stack", stack)
	root.Handle("/evaluate", &EvaluateHandler{
		Evaluator: evaluator,
	})
	root.Handle("/compare", &CompareHandler{
		Evaluator: evaluator,
		Worlds:    options.Worlds,
	})

	return nil
}

type lockedHandler struct {
	handler http.Handler
	lock    *sync.RWMutex
}

func (l *lockedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	l.handler.ServeHTTP(w, r)
}

func lockHandler(handler http.Handler, lock *sync.RWMutex) http.Handler {
	return &lockedHandler{handler: handler, lock: lock}
}

func RegisterTiles(root *http.ServeMux, options *Options) {
	rules := renderer.BasemapRenderRules
	if options.BasemapRules != nil {
		rules = options.BasemapRules
	}
	base := http.Handler(lockHandler(&renderer.TileHandler{Renderer: &renderer.BasemapRenderer{RenderRules: rules, Worlds: options.Worlds}}, options.Lock))
	if options.InstrumentHandler != nil {
		base = options.InstrumentHandler(base, "tiles_base")
	}
	root.Handle("/tiles/base/", base)
	query := http.Handler(lockHandler(&renderer.TileHandler{Renderer: renderer.NewQueryRenderer(options.Worlds, options.APIOptions.Cores)}, options.Lock))
	if options.InstrumentHandler != nil {
		query = options.InstrumentHandler(query, "tiles_query")
	}
	root.Handle("/tiles/query/", query)
	histogram := http.Handler(lockHandler(&renderer.TileHandler{Renderer: renderer.NewHistogramRenderer(rules, options.Worlds)}, options.Lock))
	if options.InstrumentHandler != nil {
		histogram = options.InstrumentHandler(histogram, "tiles_histogram")
	}
	root.Handle("/tiles/histogram/", histogram)
	collection := http.Handler(lockHandler(&renderer.TileHandler{Renderer: renderer.NewCollectionRenderer(rules, options.Worlds)}, options.Lock))
	if options.InstrumentHandler != nil {
		histogram = options.InstrumentHandler(histogram, "tiles_collection")
	}
	root.Handle("/tiles/collection/", collection)

}

type StartupRequest struct {
	Root          b6.CollectionID
	MapCenter     *LatLngJSON
	MapZoom       *int
	OpenDockIndex *int
	Expression    string
}

func (s *StartupRequest) FillFromURL(url *url.URL) {
	if r := url.Query().Get("r"); len(r) > 0 {
		if id := b6.FeatureIDFromString(r[1:]); id.IsValid() && id.Type == b6.FeatureTypeCollection {
			s.Root = id.ToCollectionID()
		}
	}

	if ll := url.Query().Get("ll"); len(ll) > 0 {
		if lll, err := b6.LatLngFromString(ll); err == nil {
			s.MapCenter = &LatLngJSON{
				LatE7: int(lll.Lat.E7()),
				LngE7: int(lll.Lng.E7()),
			}
		}
	}

	if z := url.Query().Get("z"); len(z) > 0 {
		if zi, err := strconv.ParseInt(z, 10, 64); err == nil {
			s.MapZoom = new(int)
			*s.MapZoom = int(zi)
		}
	}

	if d := url.Query().Get("d"); len(d) > 0 {
		if di, err := strconv.ParseInt(d, 10, 64); err == nil {
			s.OpenDockIndex = new(int)
			*s.OpenDockIndex = int(di)
		}
	}

	if e := url.Query().Get("e"); len(e) > 0 {
		s.Expression = e
	}
}

type StartupResponseJSON struct {
	Version       string              `json:"version,omitempty"`
	Docked        []*UIResponseJSON   `json:"docked,omitempty"`
	OpenDockIndex *int                `json:"openDockIndex,omitempty"`
	MapCenter     *LatLngJSON         `json:"mapCenter,omitempty"`
	MapZoom       int                 `json:"mapZoom,omitempty"`
	Root          *FeatureIDProtoJSON `json:"root,omitempty"`
	Expression    string              `json:"expression,omitempty"`
	Error         string              `json:"error,omitempty"`
	Session       uint64              `json:"session,omitempty"`
	Locked        bool                `json:"locked,omitempty"`
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
	GeoJSON []geojson.GeoJSON    `json:"geoJSON,omitempty"`
}

func NewUIResponseJSON() *UIResponseJSON {
	return &UIResponseJSON{
		Proto: &UIResponseProtoJSON{
			Stack: &pb.StackProto{},
		},
	}
}

func (u *UIResponseJSON) AddGeoJSON(g geojson.GeoJSON) {
	u.GeoJSON = append(u.GeoJSON, g)
	u.Proto.GeoJSON = append(u.Proto.GeoJSON, &pb.GeoJSONProto{
		Index: int32(len(u.GeoJSON) - 1),
	})
}

type UI interface {
	ServeStartup(request *StartupRequest, response *StartupResponseJSON, ui UI) error
	ServeStack(request *pb.UIRequestProto, response *UIResponseJSON, ui UI) error
	Render(response *UIResponseJSON, value interface{}, root b6.CollectionID, locked bool, ui UI, closeable bool) error
}

type FeatureIDProtoJSON pb.FeatureIDProto

func (b *FeatureIDProtoJSON) MarshalJSON() ([]byte, error) {
	return protojson.Marshal((*pb.FeatureIDProto)(b))
}

func (b *FeatureIDProtoJSON) UnmarshalJSON(buffer []byte) error {
	return protojson.Unmarshal(buffer, (*pb.FeatureIDProto)(b))
}

type LatLngJSON struct {
	LatE7 int `json:"latE7"`
	LngE7 int `json:"lngE7"`
}

const DefaultMapZoom = 16

type StartupHandler struct {
	UI UI
}

func (s *StartupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := StartupRequest{}
	request.FillFromURL(r.URL)
	response := StartupResponseJSON{
		Version: b6.BackendVersion,
	}
	if err := s.UI.ServeStartup(&request, &response, s.UI); err != nil {
		response.Error = err.Error()
	}
	SendJSON(response, w, r)
}

type StackHandler struct {
	UI UI
}

func (s *StackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request := &pb.UIRequestProto{}
	if !FillStackRequest(request, w, r) {
		return
	}

	response := NewUIResponseJSON()

	if err := s.UI.ServeStack(request, response, s.UI); err == nil {
		SendJSON(response, w, r)
	} else {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func FillStackRequest(request *pb.UIRequestProto, w http.ResponseWriter, r *http.Request) bool {
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
			log.Println(err.Error())
			log.Println("Bad request body")
			http.Error(w, "Bad request body", http.StatusBadRequest)
			return false
		}
	} else {
		http.Error(w, "Bad method", http.StatusMethodNotAllowed)
		return false
	}

	if request.Expression == "" && request.Node == nil {
		http.Error(w, "No expression", http.StatusBadRequest)
		return false
	}

	return true
}

func SendJSON(value interface{}, w http.ResponseWriter, r *http.Request) {
	var output bytes.Buffer
	var encoder *json.Encoder
	var toClose io.Closer
	if strings.Index(r.Header.Get("Accept-Encoding"), "gzip") >= 0 {
		compresor := gzip.NewWriter(&output)
		encoder = json.NewEncoder(compresor)
		w.Header().Set("Content-Encoding", "gzip")
		toClose = compresor
	} else {
		encoder = json.NewEncoder(&output)
	}
	err := encoder.Encode(value)
	if err == nil && toClose != nil {
		err = toClose.Close()
	}
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output.Bytes())
}

type OpenSourceUI struct {
	Worlds          ingest.Worlds
	BasemapRules    renderer.RenderRules
	FunctionSymbols api.FunctionSymbols // For function name completion
	Evaluator       api.Evaluator
	Lock            *sync.RWMutex
}

func (o *OpenSourceUI) ServeStartup(request *StartupRequest, response *StartupResponseJSON, ui UI) error {
	o.Lock.RLock()
	defer o.Lock.RUnlock()
	w := o.Worlds.FindOrCreateWorld(request.Root.FeatureID())
	if root := b6.FindCollectionByID(request.Root, w); root != nil {
		response.Locked = root.Get("locked").String() == "yes"
		c := b6.AdaptCollection[string, any](root)
		i := c.Begin()
		for {
			ok, err := i.Next()
			if err != nil {
				return fmt.Errorf("%s: %w", request.Root, err)
			} else if !ok {
				break
			}
			if i.Key() == "centroid" {
				switch v := i.Value().(type) {
				// Read a raw point from the yaml:
				//
				// collection:
				//  - - centroid
				//		- point: 55.9480999,-3.2000552
				//
				// Can be constructed, in Python, from the function `b6.ll`, like:
				// >>> b6.ll(55.948, -3.2)
				//
				case b6.Geo:
					ll := s2.LatLngFromPoint(v.Point())
					response.MapCenter = &LatLngJSON{
						LatE7: int(ll.Lat.E7()),
						LngE7: int(ll.Lng.E7()),
					}
					response.MapZoom = DefaultMapZoom

				// It was a FeatureID (of a point), so look it up.
				// Note: This fails silently for non-point features.
				case b6.FeatureID:
					if centroid := w.FindFeatureByID(v); centroid != nil {
						if p, ok := centroid.(b6.PhysicalFeature); ok {
							ll := s2.LatLngFromPoint(p.Point())
							response.MapCenter = &LatLngJSON{
								LatE7: int(ll.Lat.E7()),
								LngE7: int(ll.Lng.E7()),
							}
							response.MapZoom = DefaultMapZoom
						}
					}
				default:
					return fmt.Errorf("Couldn't interpret centroid in world %s of type %T", request.Root, i.Value())
				}
			} else if i.Key() == "docked" {
				if featureId, ok := i.Value().(b6.FeatureID); ok {
					if docked := w.FindFeatureByID(featureId); docked != nil {
						uiResponse := NewUIResponseJSON()
						if err := ui.Render(uiResponse, docked, request.Root, true, ui, false); err == nil {
							stripShellLinesFromResponse(uiResponse)
							response.Docked = append(response.Docked, uiResponse)
						} else {
							return fmt.Errorf("%s: %w", i.Value(), err)
						}
					}
				}
			}
		}
	}

	if request.Root.IsValid() {
		id := b6.NewProtoFromFeatureID(request.Root.FeatureID())
		response.Root = (*FeatureIDProtoJSON)(id)
	}

	if request.MapCenter != nil {
		response.MapCenter = request.MapCenter
	}
	if request.MapZoom != nil {
		response.MapZoom = *request.MapZoom
	}
	if request.OpenDockIndex != nil {
		response.OpenDockIndex = request.OpenDockIndex
	}
	response.Expression = request.Expression
	response.Session = uint64(rand.Int63n(MaxSafeJavaScriptInteger))
	return nil
}

func (o *OpenSourceUI) ServeStack(request *pb.UIRequestProto, response *UIResponseJSON, ui UI) error {
	o.Lock.RLock()
	defer o.Lock.RUnlock()
	root := b6.NewFeatureIDFromProto(request.Root)

	var expression b6.Expression
	var err error
	if request.Expression != "" {
		if request.Node == nil {
			expression, err = api.ParseExpression(request.Expression)
		} else {
			if lhs, err := b6.ExpressionFromProto(request.Node); err == nil {
				expression, err = api.ParseExpressionWithLHS(request.Expression, lhs)
			}
		}

	} else {
		expression, err = b6.ExpressionFromProto(request.Node)
	}
	if err != nil {
		ui.Render(response, err, root.ToCollectionID(), request.Locked, ui, true)
		var substack pb.SubstackProto
		fillSubstackFromError(&substack, err)
		response.Proto.Stack.Substacks = append(response.Proto.Stack.Substacks, &substack)
		return nil
	}

	if !request.Locked {
		substack := &pb.SubstackProto{}
		fillSubstackFromExpression(substack, expression, true)
		if len(substack.Lines) > 0 {
			response.Proto.Stack.Substacks = append(response.Proto.Stack.Substacks, substack)
		}
	}

	if unparsed, ok := api.UnparseExpression(expression); ok {
		response.Proto.Expression = unparsed
	}

	if response.Proto.Node, err = expression.ToProto(); err != nil {
		ui.Render(response, err, root.ToCollectionID(), request.Locked, ui, true)
		return nil
	}

	var result interface{}
	result, err = o.Evaluator.EvaluateExpression(expression, root)

	if err == nil {
		if a, ok := result.(*api.AppliedChange); ok {
			if f, ok := o.uiFeature(a.Modified, root); ok {
				err = ui.Render(response, f, root.ToCollectionID(), request.Locked, ui, true)
			} else {
				err = ui.Render(response, a.Modified, root.ToCollectionID(), request.Locked, ui, true)
			}
			response.Proto.TilesChanged = true
		} else {
			err = ui.Render(response, result, root.ToCollectionID(), request.Locked, ui, true)
		}
	}
	if err != nil {
		ui.Render(response, err, root.ToCollectionID(), request.Locked, ui, true)
	}
	return nil
}

func (o *OpenSourceUI) uiFeature(c b6.UntypedCollection, root b6.FeatureID) (b6.Feature, bool) {
	w := o.Worlds.FindOrCreateWorld(root)
	i := c.BeginUntyped()
	for {
		ok, err := i.Next()
		if !ok || err != nil {
			return nil, false
		}
		if id, ok := i.Value().(b6.Identifiable); ok {
			if f := w.FindFeatureByID(id.FeatureID()); f != nil {
				if t := f.Get("b6"); t.IsValid() && t.Value.String() == "histogram" {
					return f, true
				}
			}
		}
	}
}

func (o *OpenSourceUI) Render(response *UIResponseJSON, value interface{}, root b6.CollectionID, locked bool, ui UI, closeable bool) error {
	if err := o.fillResponseFromResult(response, value, o.Worlds.FindOrCreateWorld(root.FeatureID()), closeable); err == nil {
		shell := &pb.ShellLineProto{
			Functions: make([]string, 0),
		}
		shell.Functions = fillMatchingFunctionSymbols(shell.Functions, value, o.FunctionSymbols)
		response.Proto.Stack.Substacks = append(response.Proto.Stack.Substacks, &pb.SubstackProto{
			Lines: []*pb.LineProto{{Line: &pb.LineProto_Shell{Shell: shell}}},
		})
		return nil
	} else {
		return o.fillResponseFromResult(response, err, o.Worlds.FindOrCreateWorld(root.FeatureID()), closeable)
	}
}

func (o *OpenSourceUI) fillResponseFromResult(response *UIResponseJSON, result interface{}, w b6.World, closeable bool) error {
	p := (*pb.UIResponseProto)(response.Proto)
	switch r := result.(type) {
	case error:
		var substack pb.SubstackProto
		fillSubstackFromError(&substack, r)
		p.Stack.Substacks = append(p.Stack.Substacks, &substack)
	case string:
		var substack pb.SubstackProto
		fillSubstackFromAtom(&substack, AtomFromString(r))
		p.Stack.Substacks = append(p.Stack.Substacks, &substack)
	case b6.Feature:
		if title := r.Get("b6:title"); title.IsValid() {
			var substack pb.SubstackProto
			fillSubstackFromAtom(&substack, AtomFromString(title.Value.String()))
			p.Stack.Substacks = append(p.Stack.Substacks, &substack)
		}
		switch r.FeatureID().Type {
		case b6.FeatureTypeExpression:
			// This is not perfect, as it makes original expression that
			// returned the ExpressionFeature, and the expression from the
			// feature itself, look like part of the same stack.
			// TODO: improve the UX for expression features
			substack := &pb.SubstackProto{}
			expression := api.AddPipelines(api.Simplify(b6.NewCallExpression(r.Get(b6.ExpressionTag).Value, []b6.Expression{}), o.FunctionSymbols))
			fillSubstackFromExpression(substack, expression, true)
			if len(substack.Lines) > 0 {
				response.Proto.Stack.Substacks = append(response.Proto.Stack.Substacks, substack)
			}
			id := b6.MakeCollectionID(r.FeatureID().Namespace, r.FeatureID().Value)
			if c := b6.FindCollectionByID(id, w); c != nil {
				substack := &pb.SubstackProto{}
				if err := fillSubstackFromCollection(substack, c, p, w); err == nil {
					p.Stack.Substacks = append(p.Stack.Substacks, substack)
				} else {
					return err
				}
			}
		default:
			p.Stack.Id = b6.NewProtoFromFeatureID(r.FeatureID())
			if c, ok := r.(b6.CollectionFeature); ok {
				if b6 := c.Get("b6"); b6.Value.String() == "histogram" {
					return fillResponseFromHistogramFeature(response, c, w)
				}
			}
			p.Stack.Substacks = fillSubstacksFromFeature(response, p.Stack.Substacks, r, w, closeable)
			highlightInResponse(p, r.FeatureID())
			if p, ok := r.(b6.PhysicalFeature); ok {
				// Note: We explicitly do _not_ allow this to be evaluated on paths. I
				// think there's a few reasons why:
				//
				//	1. We probably don't want that anyway
				//
				//	2. It gives a `panic` like `Expected a latlng` somewhere; so
				//	there's some assumption that gets broken for paths that can be
				//	investigated more deeply later.
				//
				// Note: This means the "Toggle Visibility" button doesn't actually
				// work properly; i.e. it doesn't show/hide the GeoJSON layer.
				if r.FeatureID().Type == b6.FeatureTypePoint ||
					r.FeatureID().Type == b6.FeatureTypeArea ||
					r.FeatureID().Type == b6.FeatureTypeRelation {
					response.AddGeoJSON(p.ToGeoJSON())
				}
			}
		}
	case b6.FeatureID:
		if f := w.FindFeatureByID(r); f != nil {
			return o.fillResponseFromResult(response, f, w, closeable)
		} else {
			return o.fillResponseFromResult(response, r.String(), w, closeable)
		}
	case b6.Query:
		if q, ok := api.UnparseQuery(r); ok {
			var substack pb.SubstackProto
			fillSubstackFromAtom(&substack, AtomFromString(q))
			p.Stack.Substacks = append(p.Stack.Substacks, &substack)
			p.Layers = append(p.Layers, &pb.MapLayerProto{
				Path:   "query",
				Q:      q,
				Before: pb.MapLayerPosition_MapLayerPositionEnd,
			})
		} else {
			// TODO: Improve the rendering of queries
			var substack pb.SubstackProto
			fillSubstackFromAtom(&substack, AtomFromString("Query"))
			p.Stack.Substacks = append(p.Stack.Substacks, &substack)
		}
	case b6.Tag:
		var substack pb.SubstackProto
		fillSubstackFromAtom(&substack, AtomFromValue(r, w))
		p.Stack.Substacks = append(p.Stack.Substacks, &substack)
		if !o.BasemapRules.IsRendered(r) {
			if q, ok := api.UnparseQuery(b6.Tagged(r)); ok {
				before := pb.MapLayerPosition_MapLayerPositionEnd
				if r.Key == "#boundary" {
					before = pb.MapLayerPosition_MapLayerPositionBuildings
				}
				p.Layers = append(p.Layers, &pb.MapLayerProto{
					Path:   "query",
					Q:      q,
					Before: before,
				})
			}
		}
	case b6.UntypedCollection:
		substack := &pb.SubstackProto{}
		if err := fillSubstackFromCollection(substack, r, p, w); err == nil {
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
		var substack pb.SubstackProto
		fillSubstackFromAtom(&substack, atom)
		p.Stack.Substacks = append(p.Stack.Substacks, &substack)
		response.AddGeoJSON(r.ToGeoJSON())
	case b6.Geometry:
		switch r.GeometryType() {
		case b6.GeometryTypePoint:
			ll := s2.LatLngFromPoint(r.Point())
			atom := &pb.AtomProto{
				Atom: &pb.AtomProto_Value{
					Value: fmt.Sprintf("%f, %f", ll.Lat.Degrees(), ll.Lng.Degrees()),
				},
			}
			var substack1 pb.SubstackProto
			fillSubstackFromAtom(&substack1, atom)
			p.Stack.Substacks = append(p.Stack.Substacks, &substack1)
			response.AddGeoJSON(r.ToGeoJSON())
		case b6.GeometryTypePath:
			dimension := b6.AngleToMeters(r.Polyline().Length())
			atom := &pb.AtomProto{
				Atom: &pb.AtomProto_Download{
					Download: fmt.Sprintf("%.2fm path", dimension),
				},
			}
			var substack pb.SubstackProto
			fillSubstackFromAtom(&substack, atom)
			p.Stack.Substacks = append(p.Stack.Substacks, &substack)
			response.AddGeoJSON(r.ToGeoJSON())
		default:
			return o.fillResponseFromResult(response, r.ToGeoJSON(), w, closeable)
		}
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
		var substack pb.SubstackProto
		fillSubstackFromAtom(&substack, atom)
		p.Stack.Substacks = append(p.Stack.Substacks, &substack)
		response.AddGeoJSON(r)
	case *geojson.Feature:
		atom := &pb.AtomProto{
			Atom: &pb.AtomProto_Download{
				Download: "GeoJSON feature",
			},
		}
		var substack pb.SubstackProto
		fillSubstackFromAtom(&substack, atom)
		p.Stack.Substacks = append(p.Stack.Substacks, &substack)
		response.AddGeoJSON(r)
	case *geojson.Geometry:
		atom := &pb.AtomProto{
			Atom: &pb.AtomProto_Download{
				Download: "GeoJSON geometry",
			},
		}
		var substack pb.SubstackProto
		fillSubstackFromAtom(&substack, atom)
		p.Stack.Substacks = append(p.Stack.Substacks, &substack)
		response.AddGeoJSON(geojson.NewFeatureWithGeometry(*r))
	default:
		substack := &pb.SubstackProto{
			Lines: []*pb.LineProto{{
				Line: &pb.LineProto_Value{
					Value: &pb.ValueLineProto{
						Atom: AtomFromValue(r, w),
					},
				},
			}},
		}
		p.Stack.Substacks = append(p.Stack.Substacks, substack)
	}
	switch r := result.(type) {
	case b6.Geometry:
		if centroid, ok := b6.Centroid(r); ok {
			response.Proto.MapCenter = b6.NewPointProtoFromS2Point(centroid)
		}
	case geojson.GeoJSON:
		response.Proto.MapCenter = b6.NewPointProtoFromS2Point(r.Centroid().ToS2Point())
	}
	return nil
}

type EvaluateResponseProtoJSON pb.EvaluateResponseProto

func (e *EvaluateResponseProtoJSON) MarshalJSON() ([]byte, error) {
	return protojson.Marshal((*pb.EvaluateResponseProto)(e))
}

func (e *EvaluateResponseProtoJSON) UnmarshalJSON(buffer []byte) error {
	return protojson.Unmarshal(buffer, (*pb.EvaluateResponseProto)(e))
}

type EvaluateHandler struct {
	Evaluator api.Evaluator
}

func (e *EvaluateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var result interface{}
	if r.Method == "GET" {
		q := r.URL.Query()
		root := b6.FeatureIDFromString(q.Get("r"))
		if result, err = e.Evaluator.EvaluateString(q.Get("e"), root); err == nil {
			if a, ok := result.(*api.AppliedChange); ok {
				result = a.Modified
			}
		}
	} else if r.Method == "POST" {
		var body []byte
		if body, err = io.ReadAll(r.Body); err == nil {
			r.Body.Close()
			var request pb.EvaluateRequestProto
			if err = protojson.Unmarshal(body, &request); err == nil {
				result, err = e.Evaluator.EvaluateProto(&request)
			}
		}
	} else {
		err = fmt.Errorf("Bad HTTP method")
	}
	var literal b6.Literal
	if err == nil {
		literal, err = b6.FromLiteral(result)
	}
	var node *pb.NodeProto
	if err == nil {
		node, err = literal.ToProto()
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := &pb.EvaluateResponseProto{
		Result: node,
	}
	SendJSON((*EvaluateResponseProtoJSON)(response), w, r)
}

type ProtoJSON[Proto proto.Message] struct {
	m Proto
}

func (p ProtoJSON[_]) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(p.m)
}

func (p ProtoJSON[_]) UnmarshalJSON(buffer []byte) error {
	return protojson.Unmarshal(buffer, p.m)
}

func WrapProtoForJSON[Proto proto.Message](m Proto) ProtoJSON[Proto] {
	return ProtoJSON[Proto]{m: m}
}

type CompareHandler struct {
	Evaluator api.Evaluator
	Worlds    ingest.Worlds
}

func (c *CompareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var request pb.ComparisonRequestProto
	var body []byte
	var err error
	if body, err = io.ReadAll(r.Body); err == nil {
		err = protojson.Unmarshal(body, &request)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	baseline := c.Worlds.FindOrCreateWorld(b6.NewFeatureIDFromProto(request.Baseline))
	scenarios := make([]b6.FeatureID, len(request.Scenarios))
	for i := range scenarios {
		scenarios[i] = b6.NewFeatureIDFromProto(request.Scenarios[i])
	}

	var analysis b6.CollectionFeature
	var expression b6.Feature
	id := b6.NewFeatureIDFromProto(request.Analysis)
	if id.Type == b6.FeatureTypeCollection {
		if analysis = b6.FindCollectionByID(id.ToCollectionID(), baseline); analysis != nil {
			id.Type = b6.FeatureTypeExpression
			if f := baseline.FindFeatureByID(id); f == nil {
				err = fmt.Errorf("no expression with ID %s", id)
			} else {
				expression = f
			}
		} else {
			err = fmt.Errorf("no collection with ID %s", id)
		}
	} else {
		err = fmt.Errorf("analysis is not a collection")
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var response pb.ComparisonLineProto
	response.Baseline = newComparisonHistogram(analysis, baseline)

	c.Evaluator.Lock.RLock()
	defer c.Evaluator.Lock.RUnlock()
	for _, scenario := range scenarios {
		if _, err = c.Evaluator.EvaluateExpression(expression.Get(b6.ExpressionTag).Value, scenario); err == nil {
			w := c.Worlds.FindOrCreateWorld(scenario)
			if comparison := b6.FindCollectionByID(analysis.CollectionID(), w); comparison != nil {
				response.Scenarios = append(response.Scenarios, newComparisonHistogram(comparison, w))
			} else {
				err = fmt.Errorf("expression didn't produce required analysis")
			}
		}
		if err != nil {
			break
		}
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	allHistograms := make([]*pb.ComparisonHistogramProto, 1, 1+len(request.Scenarios))
	allHistograms[0] = response.Baseline
	allHistograms = append(allHistograms, response.Scenarios...)
	equaliseBars(allHistograms)

	SendJSON(WrapProtoForJSON(&response), w, r)
}

func newComparisonHistogram(c b6.CollectionFeature, w b6.World) *pb.ComparisonHistogramProto {
	var comparison pb.ComparisonHistogramProto
	response := NewUIResponseJSON()
	fillResponseFromHistogramFeature(response, c, w)
	p := (*pb.UIResponseProto)(response.Proto)
	for _, s := range p.Stack.Substacks {
		for _, line := range s.Lines {
			if bar, ok := line.Line.(*pb.LineProto_HistogramBar); ok {
				comparison.Bars = append(comparison.Bars, bar.HistogramBar)
			}
		}
	}
	return &comparison
}

func equaliseBars(histograms []*pb.ComparisonHistogramProto) {
	byKey := make(map[string]*pb.AtomProto)
	for _, h := range histograms {
		for _, bar := range h.Bars {
			byKey[SortableKeyForAtom(bar.Range)] = bar.Range
		}
	}
	keys := make([]string, 0, len(byKey))
	for key := range byKey {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, h := range histograms {
		existing := make(map[string]*pb.HistogramBarLineProto)
		for _, bar := range h.Bars {
			existing[SortableKeyForAtom(bar.Range)] = bar
		}
		total := 0
		equalised := make([]*pb.HistogramBarLineProto, 0, len(keys))
		for _, key := range keys {
			if bar, ok := existing[key]; ok {
				total = int(bar.Total)
				equalised = append(equalised, bar)
			} else {
				added := &pb.HistogramBarLineProto{
					Range: byKey[key],
					Value: 0,
				}
				equalised = append(equalised, added)
			}
		}
		for i, bar := range equalised {
			bar.Index = int32(i)
			bar.Total = int32(total)
		}
		h.Bars = equalised
	}
}

func stripShellLinesFromResponse(response *UIResponseJSON) {
	p := (*pb.UIResponseProto)(response.Proto)
	stripped := make([]*pb.SubstackProto, 0, len(p.Stack.Substacks))
	for _, s := range p.Stack.Substacks {
		lines := make([]*pb.LineProto, 0, len(s.Lines))
		for _, l := range s.Lines {
			if _, ok := l.Line.(*pb.LineProto_Shell); !ok {
				lines = append(lines, l)
			}
		}
		if len(lines) > 0 {
			s.Lines = lines
			stripped = append(stripped, s)
		}
	}
	p.Stack.Substacks = stripped
}
