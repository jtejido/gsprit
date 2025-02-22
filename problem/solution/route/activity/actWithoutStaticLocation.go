package activity

import (
	"gsprit/problem"
)

// ActWithoutStaticLocation represents a service activity without a fixed static location
type ActWithoutStaticLocation struct {
	ServiceActivity
	previousLocation *problem.Location
	nextLocation     *problem.Location
}

// NewActWithoutStaticLocation creates a new instance
func NewActWithoutStaticLocation(service problem.Service) *ActWithoutStaticLocation {
	return &ActWithoutStaticLocation{
		ServiceActivity: *NewServiceActivity(service),
	}
}

// Location returns the previous location as the default location
func (a *ActWithoutStaticLocation) Location() *problem.Location {
	return a.previousLocation
}

// PreviousLocation returns the previous location
func (a *ActWithoutStaticLocation) PreviousLocation() *problem.Location {
	return a.previousLocation
}

// NextLocation returns the next location
func (a *ActWithoutStaticLocation) NextLocation() *problem.Location {
	return a.nextLocation
}

// SetPreviousLocation sets the previous location
func (a *ActWithoutStaticLocation) SetPreviousLocation(previousLocation *problem.Location) {
	a.previousLocation = previousLocation
}

// SetNextLocation sets the next location
func (a *ActWithoutStaticLocation) SetNextLocation(nextLocation *problem.Location) {
	a.nextLocation = nextLocation
}

// Duplicate creates a duplicate of the activity
func (a *ActWithoutStaticLocation) Duplicate() problem.TourActivity {
	return &ActWithoutStaticLocation{
		ServiceActivity:  *a.ServiceActivity.Duplicate().(*ServiceActivity),
		previousLocation: a.previousLocation,
		nextLocation:     a.nextLocation,
	}
}
