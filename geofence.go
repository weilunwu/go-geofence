package geofence

import "github.com/kellydunn/golang-geo"

// Geofence is a struct for efficient search whether a point is in polygon
type Geofence struct {
	vertices    []*geo.Point
	holes       [][]*geo.Point
	tiles       map[float64]string
	granularity int64
	minX        float64
	maxX        float64
	minY        float64
	maxY        float64
	tileWidth   float64
	tileHeight  float64
	minTileX    float64
	maxTileX    float64
	minTileY    float64
	maxTileY    float64
}

const defaultGranularity = 20

// NewGeofence is the construct for Geofence, vertices: {{(1,2),(2,3)}, {(1,0)}}.
// 1st array contains polygon vertices. 2nd array contains holes.
func NewGeofence(points [][]*geo.Point, args ...interface{}) *Geofence {
	geofence := &Geofence{}
	if len(args) > 0 {
		geofence.granularity = args[0].(int64)
	} else {
		geofence.granularity = defaultGranularity
	}
	geofence.vertices = points[0]
	if len(points) > 1 {
		geofence.holes = points[1:]
	}
	geofence.tiles = make(map[float64]string)

	geofence.setInclusionTiles()
	return geofence
}

// Inside checks whether a given point is inside the geofence
func (geofence *Geofence) Inside(point *geo.Point) bool {
	// Bbox check first
	if point.Lat() < geofence.minX || point.Lat() > geofence.maxX || point.Lng() < geofence.minY || point.Lng() > geofence.maxY {
		return false
	}

	tileHash := (project(point.Lng(), geofence.tileHeight)-geofence.minTileY)*float64(geofence.granularity) + (project(point.Lat(), geofence.tileWidth) - geofence.minTileX)
	intersects := geofence.tiles[tileHash]

	if intersects == "i" {
		return true
	} else if intersects == "x" {
		polygon := geo.NewPolygon(geofence.vertices)
		inside := polygon.Contains(point)
		if !inside || len(geofence.holes) == 0 {
			return inside
		}

		// if we hanve holes cut out, and the point falls within the outer ring,
		// ensure no inner rings exclude this point
		for i := 0; i < len(geofence.holes); i++ {
			holePoly := geo.NewPolygon(geofence.holes[i])
			if holePoly.Contains(point) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func (geofence *Geofence) setInclusionTiles() {
	xVertices := geofence.getXVertices()
	yVertices := geofence.getYVertices()

	geofence.minX = getMin(xVertices)
	geofence.minY = getMin(yVertices)
	geofence.maxX = getMax(xVertices)
	geofence.maxY = getMax(yVertices)

	xRange := geofence.maxX - geofence.minX
	yRange := geofence.maxY - geofence.minY
	geofence.tileWidth = xRange / float64(geofence.granularity)
	geofence.tileHeight = yRange / float64(geofence.granularity)

	geofence.minTileX = project(geofence.minX, geofence.tileWidth)
	geofence.minTileY = project(geofence.minY, geofence.tileHeight)
	geofence.maxTileX = project(geofence.maxX, geofence.tileWidth)
	geofence.maxTileY = project(geofence.maxY, geofence.tileHeight)

	geofence.setExclusionTiles(geofence.vertices, true)
	if len(geofence.holes) > 0 {
		for _, hole := range geofence.holes {
			geofence.setExclusionTiles(hole, false)
		}
	}
}

func (geofence *Geofence) setExclusionTiles(vertices []*geo.Point, inclusive bool) {
	var tileHash float64
	var bBoxPoly []*geo.Point
	for tileX := geofence.minTileX; tileX <= geofence.maxTileX; tileX++ {
		for tileY := geofence.minTileY; tileY <= geofence.maxTileY; tileY++ {
			tileHash = (tileY-geofence.minTileY)*float64(geofence.granularity) + (tileX - geofence.minTileX)
			bBoxPoly = []*geo.Point{geo.NewPoint(tileX*geofence.tileWidth, tileY*geofence.tileHeight), geo.NewPoint((tileX+1)*geofence.tileWidth, tileY*geofence.tileHeight), geo.NewPoint((tileX+1)*geofence.tileWidth, (tileY+1)*geofence.tileHeight), geo.NewPoint(tileX*geofence.tileWidth, (tileY+1)*geofence.tileHeight), geo.NewPoint(tileX*geofence.tileWidth, tileY*geofence.tileHeight)}

			if haveIntersectingEdges(bBoxPoly, vertices) || hasPointInPolygon(vertices, bBoxPoly) {
				geofence.tiles[tileHash] = "x"
			} else if hasPointInPolygon(bBoxPoly, vertices) {
				if inclusive {
					geofence.tiles[tileHash] = "i"
				} else {
					geofence.tiles[tileHash] = "o"
				}
			} // else all points are outside the poly
		}
	}
}

func (geofence *Geofence) getXVertices() []float64 {
	xVertices := make([]float64, len(geofence.vertices))
	for i := 0; i < len(geofence.vertices); i++ {
		xVertices[i] = geofence.vertices[i].Lat()
	}
	return xVertices
}

func (geofence *Geofence) getYVertices() []float64 {
	yVertices := make([]float64, len(geofence.vertices))
	for i := 0; i < len(geofence.vertices); i++ {
		yVertices[i] = geofence.vertices[i].Lng()
	}
	return yVertices
}

func getMin(slice []float64) float64 {
	var min float64
	if len(slice) > 0 {
		min = slice[0]
	}
	for i := 1; i < len(slice); i++ {
		if slice[i] < min {
			min = slice[i]
		}
	}
	return min
}

func getMax(slice []float64) float64 {
	var max float64
	if len(slice) > 0 {
		max = slice[0]
	}
	for i := 1; i < len(slice); i++ {
		if slice[i] > max {
			max = slice[i]
		}
	}
	return max
}
