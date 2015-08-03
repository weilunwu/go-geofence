package geofence

import "github.com/kellydunn/golang-geo"

// region is a struct to store good and bad areas in region
type region struct {
	whiteouts []*Geofence
	blackouts []*Geofence
}

// GeofenceGroup is a struct to store k-v pair of id and regions
type GeofenceGroup struct {
	entries map[int]*region
}

// NewGeofenceGroup is the constructor for GeofenceGroup
func NewGeofenceGroup() *GeofenceGroup {
	entries := make(map[int]*region)
	return &GeofenceGroup{entries: entries}
}

// Add adds an fencelist to the group
func (group *GeofenceGroup) Add(id int, whiteoutGfs []*Geofence, blackoutGfs []*Geofence) {
	entry, ok := group.entries[id]
	if ok {
		entry.whiteouts = append(entry.whiteouts, whiteoutGfs...)
		entry.blackouts = append(entry.blackouts, blackoutGfs...)
	}

	list := &region{whiteouts: whiteoutGfs, blackouts: blackoutGfs}
	group.entries[id] = list
}

// IsGoodPoint checks whether the point is in the id-th fencelist's whiteouts
func (group *GeofenceGroup) IsGoodPoint(id int, point *geo.Point) bool {
	entry, ok := group.entries[id]
	if ok {
		return isPointValid(point, entry.whiteouts, entry.blackouts)
	}
	return false
}

// GetValidKeys returns all group id that contains point in whiteouts
// TODO: Using map as return type is not good enough
func (group *GeofenceGroup) GetValidKeys(point *geo.Point) map[int]bool {
	result := make(map[int]bool)
	for id, entry := range group.entries {
		if isPointValid(point, entry.whiteouts, entry.blackouts) {
			result[id] = true
		}
	}
	return result
}

// GetWhiteouts return the whiteout area of a region
func (group *GeofenceGroup) GetWhiteouts(id int) []*Geofence {
	entry, ok := group.entries[id]
	if ok {
		return entry.whiteouts
	}
	return nil
}

// GetBlackouts return the blackout area of a region
func (group *GeofenceGroup) GetBlackouts(id int) []*Geofence {
	entry, ok := group.entries[id]
	if ok {
		return entry.blackouts
	}
	return nil
}

func isPointValid(point *geo.Point, whiteoutGfs []*Geofence, blackoutGfs []*Geofence) bool {
	// if a point is inside any of the blackout geofences, then it is invalid
	if blackoutGfs != nil && len(blackoutGfs) > 0 {
		for i := 0; i < len(blackoutGfs); i++ {
			if blackoutGfs[i].Inside(point) {
				return false
			}
		}
	}

	// if a point is inside any of the whiteout geofences, then it is valid
	if whiteoutGfs != nil && len(whiteoutGfs) > 0 {
		for i := 0; i < len(whiteoutGfs); i++ {
			if whiteoutGfs[i].Inside(point) {
				return true
			}
		}
		// point is not in any of whiteout geofences
		return false
	}
	// not inside any blackouts and there are no whiteouts, so point is valid
	return true
}
