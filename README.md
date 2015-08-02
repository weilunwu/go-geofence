go-geofence is a library to perform point-in-geofence searches in Golang.

Advantages compared with golang-geo: 
1. go-geofence supports searching in polygons with holes inside
2. go-geofence uses a tiled cache to store pre-computed search results so it can determine inclusion very efficiently. Therefore the library is tailored for create once, query many times uses.

Benchmark results:
BenchmarkGeofence	10000000	       109 ns/op
BenchmarkGeoContains	 3000000	       475 ns/op
Detailed benchmark tests can be found in geofence_test.go

