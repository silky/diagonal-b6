package functions

// Code generated by b6-api. DO NOT EDIT.

var functionDocs = map[string]Doc{
	"accessible-all": Doc{Doc: "Return the a collection of the features reachable from the given origins, within the given duration in seconds, that match the given query.\nKeys of the collection are origins, values are reachable destinations.\nOptions are passed as tags containing the mode, and mode specific values. Examples include:\nWalking, with the default speed of 4.5km/h:\nmode=walk\nWalking, a speed of 3km/h:\nmode=walk, walking speed=3.0\nTransit at peak times:\nmode=transit\nTransit at off-peak times:\nmode=transit, peak=no\nWalking, accounting for elevation:\nelevation=true (optional: uphill=hard downhill=hard)\nWalking, with the resulting collection flipped such that keys are\ndestinations and values are origins. Useful for efficiency if you assume\nsymmetry, and the number of destinations is considerably smaller than the\nnumber of origins:\nmode=walk, flip=yes\n", ArgNames: []string{"origins","destinations","duration","options"}},
	"accessible-routes": Doc{Doc: "", ArgNames: []string{"origin","destinations","duration","options"}},
	"add": Doc{Doc: "Return a added to b.\n", ArgNames: []string{"a","b"}},
	"add-collection": Doc{Doc: "Add a collection feature with the given id, tags and items.\n", ArgNames: []string{"id","tags","collection"}},
	"add-ints": Doc{Doc: "Deprecated.\n", ArgNames: []string{"a","b"}},
	"add-point": Doc{Doc: "Adds a point feature with the given id, tags and members.\n", ArgNames: []string{"point","id","tags"}},
	"add-relation": Doc{Doc: "Add a relation feature with the given id, tags and members.\n", ArgNames: []string{"id","tags","members"}},
	"add-tag": Doc{Doc: "Add the given tag to the given feature.\n", ArgNames: []string{"id","tag"}},
	"add-tags": Doc{Doc: "Add the given tags to the given features.\nThe keys of the given collection specify the features to change, the\nvalues provide the tag to be added.\n", ArgNames: []string{"collection"}},
	"all": Doc{Doc: "Return a query that will match any feature.\n", ArgNames: []string{}},
	"all-tags": Doc{Doc: "Return a collection of all the tags on the given feature.\nKeys are ordered integers from 0, values are tags.\n", ArgNames: []string{"id"}},
	"and": Doc{Doc: "Return a query that will match features that match both given queries.\n", ArgNames: []string{"a","b"}},
	"apply-to-area": Doc{Doc: "Wrap the given function such that it will only be called when passed an area.\n", ArgNames: []string{"f"}},
	"apply-to-path": Doc{Doc: "Wrap the given function such that it will only be called when passed a path.\n", ArgNames: []string{"f"}},
	"apply-to-point": Doc{Doc: "Wrap the given function such that it will only be called when passed a point.\n", ArgNames: []string{"f"}},
	"area": Doc{Doc: "Return the area of the given polygon in m².\n", ArgNames: []string{"area"}},
	"building-access": Doc{Doc: "Deprecated. Use accessible.\n", ArgNames: []string{"origins","limit","mode"}},
	"cap-polygon": Doc{Doc: "Return a polygon approximating a spherical cap with the given center and radius in meters.\n", ArgNames: []string{"center","radius"}},
	"centroid": Doc{Doc: "Return the centroid of the given geometry.\nFor multipolygons, we return the centroid of the convex hull formed from\nthe points of those polygons.\n", ArgNames: []string{"geometry"}},
	"changes-from-file": Doc{Doc: "Return the changes contained in the given file.\nAs the file is read by the b6 server process, the filename it relative\nto the filesystems it sees. Reading from files on cloud storage is\nsupported.\n", ArgNames: []string{"filename"}},
	"changes-to-file": Doc{Doc: "Export the changes that have been applied to the world to the given filename as yaml.\nAs the file is written by the b6 server process, the filename it relative\nto the filesystems it sees. Writing files to cloud storage is\nsupported.\n", ArgNames: []string{"filename"}},
	"clamp": Doc{Doc: "Return the given value, unless it falls outside the given inclusive bounds, in which case return the boundary.\n", ArgNames: []string{"v","low","high"}},
	"closest": Doc{Doc: "Return the closest feature from the given origin via the given mode, within the given distance in meters, matching the given query.\nSee accessible-all for options values.\n", ArgNames: []string{"origin","options","distance","query"}},
	"closest-distance": Doc{Doc: "Return the distance through the graph of the closest feature from the given origin via the given mode, within the given distance in meters, matching the given query.\nSee accessible-all for options values.\n", ArgNames: []string{"origin","options","distance","query"}},
	"collect-areas": Doc{Doc: "Return a single area containing all areas from the given collection.\nIf areas in the collection overlap, loops within the returned area\nwill overlap, which will likely cause undefined behaviour in many\nfunctions.\n", ArgNames: []string{"areas"}},
	"collection": Doc{Doc: "Return a collection of the given key value pairs.\n", ArgNames: []string{"pairs"}},
	"connect": Doc{Doc: "Add a path that connects the two given points, if they're not already directly connected.\n", ArgNames: []string{"a","b"}},
	"connect-to-network": Doc{Doc: "Add a path and point to connect given feature to the street network.\nThe street network is defined at the set of paths tagged #highway that\nallow traversal of more than 500m. A point is added to the closest\nnetwork path at the projection of the origin point on that path, unless\nthat point is within 4m of an existing path point.\n", ArgNames: []string{"feature"}},
	"connect-to-network-all": Doc{Doc: "Add paths and points to connect the given collection of features to the\nnetwork. See connect-to-network for connection details.\nMore efficient than using map with connect-to-network, as the street\nnetwork is only computed once.\n", ArgNames: []string{"features"}},
	"containing-areas": Doc{Doc: "", ArgNames: []string{"points","q"}},
	"convex-hull": Doc{Doc: "Return the convex hull of the given geometries.\n", ArgNames: []string{"c"}},
	"count": Doc{Doc: "Return the number of items in the given collection.\nThe function will not evaluate and traverse the entire collection if it's possible to count\nthe collection efficiently.\n", ArgNames: []string{"collection"}},
	"count-tag-value": Doc{Doc: "Deprecated.\n", ArgNames: []string{"id","key"}},
	"count-valid-ids": Doc{Doc: "Return the number of valid feature IDs in the given collection\n", ArgNames: []string{"collection"}},
	"count-values": Doc{Doc: "Return a collection of the number of occurances of each value in the given collection.\n", ArgNames: []string{"collection"}},
	"debug-all-query": Doc{Doc: "Deprecated.\n", ArgNames: []string{"token"}},
	"debug-tokens": Doc{Doc: "Return the search index tokens generated for the given feature.\nIntended for debugging use only.\n", ArgNames: []string{"id"}},
	"degree": Doc{Doc: "Return the number of paths connected to the given point.\nA single path will be counted twice if the point isn't at one of its\ntwo ends - once in one direction, and once in the other.\n", ArgNames: []string{"point"}},
	"distance-meters": Doc{Doc: "Return the distance in meters between the given points.\n", ArgNames: []string{"a","b"}},
	"distance-to-point-meters": Doc{Doc: "Return the distance in meters between the given path, and the project of the give point onto it.\n", ArgNames: []string{"path","point"}},
	"divide": Doc{Doc: "Return a divided by b.\n", ArgNames: []string{"a","b"}},
	"divide-int": Doc{Doc: "Deprecated.\n", ArgNames: []string{"a","b"}},
	"entrance-approach": Doc{Doc: "", ArgNames: []string{"area"}},
	"export-world": Doc{Doc: "Write the current world to the given filename in the b6 compact index format.\nAs the file is written by the b6 server process, the filename it relative\nto the filesystems it sees. Writing files to cloud storage is\nsupported.\n", ArgNames: []string{"filename"}},
	"filter": Doc{Doc: "Return a collection of the items of the given collection for which the value of the given function applied to each value is true.\n", ArgNames: []string{"collection","function"}},
	"filter-accessible": Doc{Doc: "Return a collection containing only the values of the given collection that match the given query.\nIf no values for a key match the query, emit a single invalid feature ID\nfor that key, allowing callers to count the number of keys with no valid\nvalues.\nKeys are taken from the given collection.\n", ArgNames: []string{"collection","filter"}},
	"find": Doc{Doc: "Return a collection of the features present in the world that match the given query.\nKeys are IDs, and values are features.\n", ArgNames: []string{"query"}},
	"find-area": Doc{Doc: "Return the area feature with the given ID.\n", ArgNames: []string{"id"}},
	"find-areas": Doc{Doc: "Return a collection of the area features present in the world that match the given query.\nKeys are IDs, and values are features.\n", ArgNames: []string{"query"}},
	"find-collection": Doc{Doc: "Return the collection feature with the given ID.\n", ArgNames: []string{"id"}},
	"find-expression": Doc{Doc: "Return the expression feature with the given ID.\n", ArgNames: []string{"id"}},
	"find-feature": Doc{Doc: "Return the feature with the given ID.\n", ArgNames: []string{"id"}},
	"find-relation": Doc{Doc: "Return the relation feature with the given ID.\n", ArgNames: []string{"id"}},
	"find-relations": Doc{Doc: "Return a collection of the relation features present in the world that match the given query.\nKeys are IDs, and values are features.\n", ArgNames: []string{"query"}},
	"first": Doc{Doc: "Return the first value of the given pair.\n", ArgNames: []string{"pair"}},
	"flatten": Doc{Doc: "Return a collection with keys and values taken from the collections that form the values of the given collection.\n", ArgNames: []string{"collection"}},
	"float-value": Doc{Doc: "Return the value of the given tag as a float.\nPropagates error if the value isn't a valid float.\n", ArgNames: []string{"tag"}},
	"geojson-areas": Doc{Doc: "Return the areas present in the given geojson.\n", ArgNames: []string{"g"}},
	"get": Doc{Doc: "Return the tag with the given key on the given feature.\nReturns a tag. To return the string value of a tag, use get-string.\n", ArgNames: []string{"id","key"}},
	"get-float": Doc{Doc: "Return the value of tag with the given key on the given feature as a float.\nReturns error if there isn't a feature with that id, a tag with that key, or if the value isn't a valid float.\n", ArgNames: []string{"id","key"}},
	"get-int": Doc{Doc: "Return the value of tag with the given key on the given feature as an integer.\nReturns error if there isn't a feature with that id, a tag with that key, or if the value isn't a valid integer.\n", ArgNames: []string{"id","key"}},
	"get-string": Doc{Doc: "Return the value of tag with the given key on the given feature as a string.\nReturns an empty string if there isn't a tag with that key.\n", ArgNames: []string{"id","key"}},
	"gt": Doc{Doc: "Return true if a is greater than b.\n", ArgNames: []string{"a","b"}},
	"histogram": Doc{Doc: "Return a change that adds a histogram for the given collection.\n", ArgNames: []string{"collection"}},
	"histogram-swatch": Doc{Doc: "Return a change that adds a histogram with only colour swatches for the given collection.\n", ArgNames: []string{"collection"}},
	"id-to-relation-id": Doc{Doc: "Deprecated.\n", ArgNames: []string{"namespace","id"}},
	"import-geojson": Doc{Doc: "Add features from the given geojson to the world.\nIDs are formed from the given namespace, and the index of the feature\nwithin the geojson collection (or 0, if a single feature is used).\n", ArgNames: []string{"features","namespace"}},
	"import-geojson-file": Doc{Doc: "Add features from the given geojson file to the world.\nIDs are formed from the given namespace, and the index of the feature\nwithin the geojson collection (or 0, if a single feature is used).\nAs the file is read by the b6 server process, the filename it relative\nto the filesystems it sees. Reading from files on cloud storage is\nsupported.\n", ArgNames: []string{"filename","namespace"}},
	"int-value": Doc{Doc: "Return the value of the given tag as an integer.\nPropagates error if the value isn't a valid integer.\n", ArgNames: []string{"tag"}},
	"interpolate": Doc{Doc: "Return the point at the given fraction along the given path.\n", ArgNames: []string{"path","fraction"}},
	"intersecting": Doc{Doc: "Return a query that will match features that intersect the given geometry.\n", ArgNames: []string{"geometry"}},
	"intersecting-cap": Doc{Doc: "Return a query that will match features that intersect a spherical cap centred on the given point, with the given radius in meters.\n", ArgNames: []string{"center","radius"}},
	"join": Doc{Doc: "Return a path formed from the points of the two given paths, in the order they occur in those paths.\n", ArgNames: []string{"pathA","pathB"}},
	"keyed": Doc{Doc: "Return a query that will match features tagged with the given key independent of value.\n", ArgNames: []string{"key"}},
	"length": Doc{Doc: "Return the length of the given path in meters.\n", ArgNames: []string{"path"}},
	"ll": Doc{Doc: "Return a point at the given latitude and longitude, specified in degrees.\n", ArgNames: []string{"lat","lng"}},
	"map": Doc{Doc: "Return a collection with the result of applying the given function to each value.\nKeys are unmodified.\n", ArgNames: []string{"collection","function"}},
	"map-geometries": Doc{Doc: "Return a geojson representing the result of applying the given function to each geometry in the given geojson.\n", ArgNames: []string{"g","f"}},
	"map-items": Doc{Doc: "Return a collection of the result of applying the given function to each pair(key, value).\nKeys are unmodified.\n", ArgNames: []string{"collection","function"}},
	"map-parallel": Doc{Doc: "Return a collection with the result of applying the given function to each value.\nKeys are unmodified, and function application occurs in parallel, bounded\nby the number of CPU cores allocated to b6.\n", ArgNames: []string{"collection","function"}},
	"matches": Doc{Doc: "Return true if the given feature matches the given query.\n", ArgNames: []string{"id","query"}},
	"materialise": Doc{Doc: "Return a change that adds a collection feature to the world with the given ID, containing the result of calling the given function.\nThe given function isn't passed any arguments.\nAlso adds an expression feature (with the same namespace and value)\nrepresenting the given function.\n", ArgNames: []string{"id","function"}},
	"materialise-map": Doc{Doc: "", ArgNames: []string{"collection","id","function"}},
	"merge-changes": Doc{Doc: "Return a change that will apply all the changes in the given collection.\nChanges are applied transactionally. If the application of one change\nfails (for example, because it includes a path that references a missing\npoint), then no changes will be applied.\n", ArgNames: []string{"collection"}},
	"or": Doc{Doc: "Return a query that will match features that match either of the given queries.\n", ArgNames: []string{"a","b"}},
	"ordered-join": Doc{Doc: "Returns a path formed by joining the two given paths.\nIf necessary to maintain consistency, the order of points is reversed,\ndetermined by which points are shared between the paths. Returns an error\nif no endpoints are shared.\n", ArgNames: []string{"pathA","pathB"}},
	"pair": Doc{Doc: "Return a pair containing the given values.\n", ArgNames: []string{"first","second"}},
	"parse-geojson": Doc{Doc: "Return the geojson represented by the given string.\n", ArgNames: []string{"s"}},
	"parse-geojson-file": Doc{Doc: "Return the geojson contained in the given file.\nAs the file is read by the b6 server process, the filename it relative\nto the filesystems it sees. Reading from files on cloud storage is\nsupported.\n", ArgNames: []string{"filename"}},
	"paths-to-reach": Doc{Doc: "Return a collection of the paths used to reach all features matching the given query from the given origin via the given mode, within the given distance in meters.\nKeys are the paths used, values are the number of times that path was used during traversal.\nSee accessible-all for options values.\n", ArgNames: []string{"origin","options","distance","query"}},
	"percentiles": Doc{Doc: "Return a collection where values represent the perentile of the corresponding value in the given collection.\nThe returned collection is ordered by percentile, with keys drawn from the\ngiven collection.\n", ArgNames: []string{"collection"}},
	"point-features": Doc{Doc: "Return a collection of the point features referenced by the given feature.\nKeys are ids of the respective value, values are point features. Area\nfeatures return the points referenced by their path features.\n", ArgNames: []string{"f"}},
	"point-paths": Doc{Doc: "Return a collection of the path features referencing the given point.\nKeys are the ids of the respective paths.\n", ArgNames: []string{"id"}},
	"points": Doc{Doc: "Return a collection of the points of the given geometry.\nKeys are ordered integers from 0, values are points.\n", ArgNames: []string{"geometry"}},
	"reachable": Doc{Doc: "Return the a collection of the features reachable from the given origin via the given mode, within the given distance in meters, that match the given query.\nSee accessible-all for options values.\nDeprecated. Use accessible-all.\n", ArgNames: []string{"origin","options","distance","query"}},
	"reachable-area": Doc{Doc: "Return the area formed by the convex hull of the features matching the given query reachable from the given origin via the given mode specified in options, within the given distance in meters.\nSee accessible-all for options values.\n", ArgNames: []string{"origin","options","distance"}},
	"rectangle-polygon": Doc{Doc: "Return a rectangle polygon with the given top left and bottom right points.\n", ArgNames: []string{"a","b"}},
	"remove-tag": Doc{Doc: "Remove the tag with the given key from the given feature.\n", ArgNames: []string{"id","key"}},
	"remove-tags": Doc{Doc: "Remove the given tags from the given features.\nThe keys of the given collection specify the features to change, the\nvalues provide the key of the tag to be removed.\n", ArgNames: []string{"collection"}},
	"s2-center": Doc{Doc: "Return a collection the center of the s2 cell with the given token.\n", ArgNames: []string{"token"}},
	"s2-covering": Doc{Doc: "Return a collection of of s2 cells tokens that cover the given area at the given level.\n", ArgNames: []string{"area","minLevel","maxLevel"}},
	"s2-grid": Doc{Doc: "Return a collection of points representing the centroids of s2 cells that cover the given area at the given level.\n", ArgNames: []string{"area","level"}},
	"s2-points": Doc{Doc: "Return a collection of points representing the centroids of s2 cells that cover the given area between the given levels.\n", ArgNames: []string{"area","minLevel","maxLevel"}},
	"s2-polygon": Doc{Doc: "Return the bounding area of the s2 cell with the given token.\n", ArgNames: []string{"token"}},
	"sample-points": Doc{Doc: "Return a collection of points along the given path, with the given distance in meters between them.\nKeys are ordered integers from 0, values are points.\n", ArgNames: []string{"path","distanceMeters"}},
	"sample-points-along-paths": Doc{Doc: "Return a collection of points along the given paths, with the given distance in meters between them.\nKeys are the id of the respective path, values are points.\n", ArgNames: []string{"paths","distanceMeters"}},
	"second": Doc{Doc: "Return the second value of the given pair.\n", ArgNames: []string{"pair"}},
	"sightline": Doc{Doc: "", ArgNames: []string{"from","radius"}},
	"snap-area-edges": Doc{Doc: "Return an area formed by projecting the edges of the given polygon onto the paths present in the world matching the given query.\nPaths beyond the given threshold in meters are ignored.\n", ArgNames: []string{"area","query","threshold"}},
	"sum-by-key": Doc{Doc: "Return a collection of the result of summing the values of each item with the same key.\nRequires values to be integers.\n", ArgNames: []string{"c"}},
	"tag": Doc{Doc: "Return a tag with the given key and value.\n", ArgNames: []string{"key","value"}},
	"tagged": Doc{Doc: "Return a query that will match features tagged with the given key and value.\n", ArgNames: []string{"key","value"}},
	"take": Doc{Doc: "Return a collection with the first n entries of the given collection.\n", ArgNames: []string{"collection","n"}},
	"tile-ids": Doc{Doc: "Deprecated\n", ArgNames: []string{"feature"}},
	"tile-ids-hex": Doc{Doc: "Deprecated\n", ArgNames: []string{"feature"}},
	"tile-paths": Doc{Doc: "Return the URL paths for the tiles containing the given geometry at the given zoom level.\n", ArgNames: []string{"geometry","zoom"}},
	"to-geojson": Doc{Doc: "", ArgNames: []string{"renderable"}},
	"to-geojson-collection": Doc{Doc: "", ArgNames: []string{"renderables"}},
	"to-str": Doc{Doc: "", ArgNames: []string{"a"}},
	"top": Doc{Doc: "Return a collection with the n entries from the given collection with the greatest values.\nRequires the values of the given collection to be integers or floats.\n", ArgNames: []string{"collection","n"}},
	"type-area": Doc{Doc: "Return a query that will match area features.\n", ArgNames: []string{}},
	"type-path": Doc{Doc: "Return a query that will match path features.\n", ArgNames: []string{}},
	"type-point": Doc{Doc: "Return a query that will match point features.\n", ArgNames: []string{}},
	"typed": Doc{Doc: "Wrap a query to only match features with the given feature type.\n", ArgNames: []string{"typ","q"}},
	"value": Doc{Doc: "Return the value of the given tag as a string.\n", ArgNames: []string{"tag"}},
	"with-change": Doc{Doc: "Return the result of calling the given function in a world in which the given change has been applied.\nThe underlying world used by the server is not modified.\n", ArgNames: []string{"change","function"}},
	"within": Doc{Doc: "Return a query that will match features that intersect the given area.\nDeprecated. Use intersecting.\n", ArgNames: []string{"a"}},
	"within-cap": Doc{Doc: "Return a query that will match features that intersect a spherical cap centred on the given point, with the given radius in meters.\nDeprecated. Use intersecting-cap.\n", ArgNames: []string{"point","radius"}},
}
