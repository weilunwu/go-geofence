package geofence

import "github.com/kellydunn/golang-geo"

// Geofence is a struct for efficient search whether a point is in polygon
type Geofence struct {
	vertices    []*geo.Point
	holes       [][]*geo.Point
	granularity int64
	minX        float64
	maxX        float64
	minY        float64
	maxY        float64
	titleWidth  float64
	titleHeight float64
	minTileX    float64
	maxTileX    float64
	minTileY    float64
	maxTileY    float64
}

// NewGeofence is the construct for Geofence, vertices: {{(1,2),(2,3)}, {(1,0)}}. 1st array are polygon vertices. 2nd array are holes. TODO: pass in granularity
func NewGeofence(points [][]*geo.Point) *Geofence {
	geofence := &Geofence{}
	geofence.granularity = 20
	geofence.vertices = points[0]
	if len(points) > 1 {
		geofence.holes = points[1:]
	}

	geofence.setInclusionTiles()
	return geofence
}

func (geofence *Geofence) setInclusionTiles() {
	xVertices := geofence.getXVertices()
	yVertices := geofence.getYVertices()

	geofence.minX = getMin(xVertices)
	geofence.minY = getMin(yVertices)
	geofence.maxX = getMax(xVertices)
	geofence.maxY = getMax(yVertices)
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
