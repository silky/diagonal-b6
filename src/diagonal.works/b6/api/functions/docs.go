package functions

// Code generated by b6-api. DO NOT EDIT.

var functionDocs = map[string]Doc{
	"accessible": Doc{Doc: "", ArgNames: []string{"origins","destinations","distance","options"}},
	"add": Doc{Doc: "", ArgNames: []string{"a","b"}},
	"add-collection": Doc{Doc: "", ArgNames: []string{"id","tags","collection"}},
	"add-ints": Doc{Doc: "", ArgNames: []string{"a","b"}},
	"add-point": Doc{Doc: "", ArgNames: []string{"p","c"}},
	"add-relation": Doc{Doc: "", ArgNames: []string{"id","tags","members"}},
	"add-tag": Doc{Doc: "", ArgNames: []string{"id","tag"}},
	"add-tags": Doc{Doc: "", ArgNames: []string{"collection"}},
	"all": Doc{Doc: "", ArgNames: []string{}},
	"all-tags": Doc{Doc: "", ArgNames: []string{"id"}},
	"and": Doc{Doc: "", ArgNames: []string{"a","b"}},
	"apply-to-area": Doc{Doc: "", ArgNames: []string{"f"}},
	"apply-to-path": Doc{Doc: "", ArgNames: []string{"f"}},
	"apply-to-point": Doc{Doc: "", ArgNames: []string{"f"}},
	"area": Doc{Doc: "", ArgNames: []string{"area"}},
	"building-access": Doc{Doc: "", ArgNames: []string{"origins","limit","mode"}},
	"cap-polygon": Doc{Doc: "", ArgNames: []string{"center","radius"}},
	"centroid": Doc{Doc: "", ArgNames: []string{"geometry"}},
	"changes-from-file": Doc{Doc: "", ArgNames: []string{"filename"}},
	"changes-to-file": Doc{Doc: "", ArgNames: []string{"filename"}},
	"clamp": Doc{Doc: "", ArgNames: []string{"v","low","high"}},
	"closest": Doc{Doc: "", ArgNames: []string{"origin","mode","distance","query"}},
	"closest-distance": Doc{Doc: "findClosestFeatureDistance returns the distance to the closest matching feature.\nIdeally, we'd either return the distance along with the feature as a pair from, or\nreturn a new primitive Route instance that described the route to that feature,\nallowing distance to be derived. Neither are possible right now, so this is a\nstopgap. TODO: Improve this API.\n", ArgNames: []string{"origin","mode","distance","query"}},
	"collect-areas": Doc{Doc: "", ArgNames: []string{"c"}},
	"collection": Doc{Doc: "", ArgNames: []string{"pairs"}},
	"connect": Doc{Doc: "", ArgNames: []string{"a","b"}},
	"connect-all-to-network": Doc{Doc: "", ArgNames: []string{"features"}},
	"connect-to-network": Doc{Doc: "", ArgNames: []string{"feature"}},
	"containing-areas": Doc{Doc: "", ArgNames: []string{"points","q"}},
	"convex-hull": Doc{Doc: "", ArgNames: []string{"c"}},
	"count": Doc{Doc: "", ArgNames: []string{"collection"}},
	"count-tag-value": Doc{Doc: "", ArgNames: []string{"id","key"}},
	"count-values": Doc{Doc: "", ArgNames: []string{"c"}},
	"debug-all-query": Doc{Doc: "", ArgNames: []string{"token"}},
	"debug-tokens": Doc{Doc: "", ArgNames: []string{"id"}},
	"degree": Doc{Doc: "", ArgNames: []string{"point"}},
	"distance-meters": Doc{Doc: "", ArgNames: []string{"a","b"}},
	"distance-to-point-meters": Doc{Doc: "", ArgNames: []string{"path","point"}},
	"divide": Doc{Doc: "", ArgNames: []string{"a","b"}},
	"divide-int": Doc{Doc: "", ArgNames: []string{"a","b"}},
	"empty-points": Doc{Doc: "", ArgNames: []string{}},
	"export-world": Doc{Doc: "", ArgNames: []string{"filename"}},
	"filter": Doc{Doc: "", ArgNames: []string{"c","f"}},
	"filter-accessible": Doc{Doc: "", ArgNames: []string{"accessible","filter"}},
	"find": Doc{Doc: "", ArgNames: []string{"query"}},
	"find-area": Doc{Doc: "", ArgNames: []string{"id"}},
	"find-areas": Doc{Doc: "", ArgNames: []string{"query"}},
	"find-feature": Doc{Doc: "", ArgNames: []string{"id"}},
	"find-path": Doc{Doc: "", ArgNames: []string{"id"}},
	"find-paths": Doc{Doc: "", ArgNames: []string{"query"}},
	"find-point": Doc{Doc: "", ArgNames: []string{"id"}},
	"find-points": Doc{Doc: "", ArgNames: []string{"query"}},
	"find-relation": Doc{Doc: "", ArgNames: []string{"id"}},
	"find-relations": Doc{Doc: "", ArgNames: []string{"query"}},
	"first": Doc{Doc: "", ArgNames: []string{"pair"}},
	"flatten": Doc{Doc: "", ArgNames: []string{"c"}},
	"float-value": Doc{Doc: "", ArgNames: []string{"tag"}},
	"geojson-areas": Doc{Doc: "", ArgNames: []string{"g"}},
	"get": Doc{Doc: "", ArgNames: []string{"id","key"}},
	"get-float": Doc{Doc: "", ArgNames: []string{"id","key"}},
	"get-int": Doc{Doc: "", ArgNames: []string{"id","key"}},
	"get-string": Doc{Doc: "", ArgNames: []string{"id","key"}},
	"gt": Doc{Doc: "", ArgNames: []string{"a","b"}},
	"histogram": Doc{Doc: "", ArgNames: []string{"c"}},
	"id-to-relation-id": Doc{Doc: "", ArgNames: []string{"namespace","id"}},
	"import-geojson": Doc{Doc: "", ArgNames: []string{"g","namespace"}},
	"import-geojson-file": Doc{Doc: "", ArgNames: []string{"filename","namespace"}},
	"int-value": Doc{Doc: "", ArgNames: []string{"tag"}},
	"interpolate": Doc{Doc: "", ArgNames: []string{"p","fraction"}},
	"intersecting": Doc{Doc: "", ArgNames: []string{"g"}},
	"intersecting-cap": Doc{Doc: "", ArgNames: []string{"center","radius"}},
	"join": Doc{Doc: "", ArgNames: []string{"a","b"}},
	"keyed": Doc{Doc: "", ArgNames: []string{"key"}},
	"length": Doc{Doc: "", ArgNames: []string{"path"}},
	"ll": Doc{Doc: "", ArgNames: []string{"lat","lng"}},
	"map": Doc{Doc: "", ArgNames: []string{"collection","f"}},
	"map-geometries": Doc{Doc: "", ArgNames: []string{"g","f"}},
	"map-items": Doc{Doc: "", ArgNames: []string{"collection","f"}},
	"map-parallel": Doc{Doc: "", ArgNames: []string{"collection","f"}},
	"matches": Doc{Doc: "", ArgNames: []string{"id","query"}},
	"materialise": Doc{Doc: "", ArgNames: []string{"id","c"}},
	"merge-changes": Doc{Doc: "", ArgNames: []string{"collection"}},
	"or": Doc{Doc: "", ArgNames: []string{"a","b"}},
	"ordered-join": Doc{Doc: "orderedJoinPaths returns a new path formed by joining a and b, in that order, reversing\nthe order of the points to maintain a consistent order, determined by which points of\nthe paths are shared. Returns an error if the paths don't share an end point.\n", ArgNames: []string{"a","b"}},
	"pair": Doc{Doc: "", ArgNames: []string{"first","second"}},
	"parse-geojson": Doc{Doc: "", ArgNames: []string{"s"}},
	"parse-geojson-file": Doc{Doc: "", ArgNames: []string{"filename"}},
	"paths-to-reach": Doc{Doc: "", ArgNames: []string{"origin","mode","distance","query"}},
	"percentiles": Doc{Doc: "TODO: percentiles inefficiently calculates the exact percentile by sorting the entire\ncollection. We could use a histogram sketch instead, maybe constructed in the\nbackground with Collection\n", ArgNames: []string{"collection"}},
	"point-features": Doc{Doc: "", ArgNames: []string{"f"}},
	"point-paths": Doc{Doc: "", ArgNames: []string{"id"}},
	"points": Doc{Doc: "", ArgNames: []string{"g"}},
	"reachable": Doc{Doc: "", ArgNames: []string{"origin","mode","distance","query"}},
	"reachable-area": Doc{Doc: "", ArgNames: []string{"origin","mode","distance"}},
	"reachable-points": Doc{Doc: "", ArgNames: []string{"origin","mode","distance","query"}},
	"rectangle-polygon": Doc{Doc: "", ArgNames: []string{"p0","p1"}},
	"remove-tag": Doc{Doc: "", ArgNames: []string{"id","key"}},
	"remove-tags": Doc{Doc: "", ArgNames: []string{"collection"}},
	"s2-center": Doc{Doc: "", ArgNames: []string{"token"}},
	"s2-covering": Doc{Doc: "", ArgNames: []string{"area","minLevel","maxLevel"}},
	"s2-grid": Doc{Doc: "", ArgNames: []string{"area","level"}},
	"s2-points": Doc{Doc: "", ArgNames: []string{"area","minLevel","maxLevel"}},
	"s2-polygon": Doc{Doc: "", ArgNames: []string{"token"}},
	"sample-points": Doc{Doc: "", ArgNames: []string{"path","distanceMeters"}},
	"sample-points-along-paths": Doc{Doc: "", ArgNames: []string{"paths","distanceMeters"}},
	"second": Doc{Doc: "", ArgNames: []string{"pair"}},
	"sightline": Doc{Doc: "", ArgNames: []string{"from","radius"}},
	"snap-area-edges": Doc{Doc: "", ArgNames: []string{"g","query","threshold"}},
	"sum-by-key": Doc{Doc: "", ArgNames: []string{"c"}},
	"tag": Doc{Doc: "", ArgNames: []string{"key","value"}},
	"tagged": Doc{Doc: "", ArgNames: []string{"key","value"}},
	"take": Doc{Doc: "", ArgNames: []string{"c","n"}},
	"tile-ids": Doc{Doc: "", ArgNames: []string{"feature"}},
	"tile-ids-hex": Doc{Doc: "", ArgNames: []string{"feature"}},
	"tile-paths": Doc{Doc: "", ArgNames: []string{"g","zoom"}},
	"to-geojson": Doc{Doc: "", ArgNames: []string{"renderable"}},
	"to-geojson-collection": Doc{Doc: "", ArgNames: []string{"renderables"}},
	"top": Doc{Doc: "", ArgNames: []string{"c","n"}},
	"type-area": Doc{Doc: "", ArgNames: []string{}},
	"type-path": Doc{Doc: "", ArgNames: []string{}},
	"type-point": Doc{Doc: "", ArgNames: []string{}},
	"value": Doc{Doc: "", ArgNames: []string{"tag"}},
	"with-change": Doc{Doc: "", ArgNames: []string{"change","f"}},
	"within": Doc{Doc: "", ArgNames: []string{"a"}},
	"within-cap": Doc{Doc: "", ArgNames: []string{"p","radius"}},
}
