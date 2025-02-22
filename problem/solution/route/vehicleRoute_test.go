package route

import (
	"testing"

	"gsprit/problem"
	"gsprit/problem/driver"
	"gsprit/problem/job"
	"gsprit/problem/solution/route/activity"
	"gsprit/problem/vehicle"

	"github.com/stretchr/testify/assert"
)

// Setup test data
var testVehicle *vehicle.Vehicle
var testDriver *driver.NoDriver

func init() {
	testVehicle = vehicle.NewVehicleBuilder("v").SetStartLocation(problem.NewLocationWithID("loc")).SetType(vehicle.NewVehicleTypeBuilder("yo").Build()).Build()
	testDriver = driver.NewNoDriver()
}

// Test creating an empty route correctly
func TestBuildingEmptyRoute(t *testing.T) {
	route := NewVehicleRouteBuilder(vehicle.NewNoVehicle(), driver.NewNoDriver()).Build()
	assert.NotNil(t, route)
}

func TestBuildingEmptyRouteV2(t *testing.T) {
	route := EmptyRoute()
	assert.NotNil(t, route)
}

// Test activity iterator should iterate over zero activities in an empty route
func TestBuildingEmptyRoute_ActivityIteratorZeroActivities(t *testing.T) {
	route := EmptyRoute()
	count := len(route.TourActivities().Activities())
	assert.Equal(t, 0, count)
}

// Test building a route with null values should throw an error
func TestBuildingRouteWithNulls_ShouldThrowError(t *testing.T) {
	assert.Panics(t, func() {
		NewVehicleRouteBuilder(nil, nil)
	})
}

// Test adding service activities and iterating over them
func TestBuildingNonEmptyTour(t *testing.T) {
	routeBuilder := NewVehicleRouteBuilder(testVehicle, testDriver)
	routeBuilder.AddService(job.NewServiceBuilder[*job.Service]("2").AddSizeDimension(0, 30).SetLocation(problem.NewLocationWithID("1")).Build())
	r := routeBuilder.Build()

	assert.Equal(t, 1, len(r.TourActivities().Activities()))

	r.TourActivities().AddActivityToEnd(activity.NewServiceActivity(job.NewServiceBuilder[*job.Service]("3").AddSizeDimension(0, 30).SetLocation(problem.NewLocationWithID("1")).Build()))
	assert.Equal(t, 2, len(r.TourActivities().Activities()))
}

// Test reverse iterator over activities
func TestBuildingNonEmptyTour_ReverseIterator(t *testing.T) {
	r := NewVehicleRouteBuilder(testVehicle, testDriver).Build()
	iter := activity.NewReverseActivityIterator(r.TourActivities().Activities())
	count := 0

	for iter.HasNext() {
		iter.Next()
		count++
	}

	assert.Equal(t, 0, count)
}

func TestBuildingNonEmptyTourV2_ReverseIterator(t *testing.T) {
	routeBuilder := NewVehicleRouteBuilder(testVehicle, testDriver)
	routeBuilder.AddService(job.NewServiceBuilder[*job.Service]("2").AddSizeDimension(0, 30).SetLocation(problem.NewLocationWithID("1")).Build())
	r := routeBuilder.Build()
	iter := activity.NewReverseActivityIterator(r.TourActivities().Activities())

	count := 0
	for iter.HasNext() {
		iter.Next()
		count++
	}

	assert.Equal(t, 1, count)
}

func TestBuildingRouteWithDifferentStartAndEndLocations(t *testing.T) {
	v := vehicle.NewVehicleBuilder("v").SetStartLocation(problem.NewLocationWithID("start")).SetEndLocation(problem.NewLocationWithID("end")).Build()
	vRoute := NewVehicleRouteBuilder(v, driver.NewNoDriver()).Build()
	assert.Equal(t, "start", vRoute.Start().Location().Id())
	assert.Equal(t, "end", vRoute.End().Location().Id())
}

func TestBuildingRouteWithSameStartAndEndLocations(t *testing.T) {
	v := vehicle.NewVehicleBuilder("v").SetStartLocation(problem.NewLocationWithID("start")).SetEndLocation(problem.NewLocationWithID("start")).Build()
	vRoute := NewVehicleRouteBuilder(v, driver.NewNoDriver()).Build()
	assert.Equal(t, "start", vRoute.Start().Location().Id())
	assert.Equal(t, "start", vRoute.End().Location().Id())
}

func TestBuildingRouteWithVehicle_EarliestStartAndEndTime(t *testing.T) {
	v := vehicle.NewVehicleBuilder("v").SetStartLocation(problem.NewLocationWithID("start")).SetEndLocation(problem.NewLocationWithID("end")).SetEarliestStart(100).SetLatestArrival(200).Build()
	vRoute := NewVehicleRouteBuilder(v, driver.NewNoDriver()).Build()
	dep, err := vRoute.DepartureTime()
	assert.NoError(t, err)
	assert.Equal(t, 100.0, dep)
	assert.Equal(t, 100.0, vRoute.Start().EndTime())
	assert.Equal(t, 200.0, vRoute.End().TheoreticalLatestOperationStartTime())
}

func TestBuildingRouteWithNewVehicle_UpdatesStartAndEndLocation(t *testing.T) {
	oldV := vehicle.NewVehicleBuilder("v").SetStartLocation(problem.NewLocationWithID("start")).SetEndLocation(problem.NewLocationWithID("end")).SetEarliestStart(100).SetLatestArrival(200).Build()
	newV := vehicle.NewVehicleBuilder("new_v").SetStartLocation(problem.NewLocationWithID("new_start")).SetEndLocation(problem.NewLocationWithID("new_end")).SetEarliestStart(1000).SetLatestArrival(2000).Build()

	vRoute := NewVehicleRouteBuilder(oldV, driver.NewNoDriver()).Build()
	vRoute.SetVehicleAndDepartureTime(newV, 50.0)

	assert.Equal(t, "new_start", vRoute.Start().Location().Id())
	assert.Equal(t, "new_end", vRoute.End().Location().Id())
	dep, err := vRoute.DepartureTime()
	assert.NoError(t, err)
	assert.Equal(t, 1000.0, dep)
}

func TestAddingPickupToRoute(t *testing.T) {
	pickup := job.NewPickupBuilder("pick").SetLocation(problem.NewLocationWithID("pickLoc")).Build()
	v := vehicle.NewVehicleBuilder("vehicle").SetStartLocation(problem.NewLocationWithID("startLoc")).Build()
	r := NewVehicleRouteBuilder(v, driver.NewNoDriver()).AddService(pickup).Build()

	act := r.Activities()[0]
	assert.Equal(t, "pickup", act.Name())
	assert.IsType(t, &activity.PickupService{}, act)
}

func TestAddingDeliveryToRoute(t *testing.T) {
	delivery := job.NewDeliveryBuilder("delivery").SetLocation(problem.NewLocationWithID("deliveryLoc")).Build()
	v := vehicle.NewVehicleBuilder("vehicle").SetStartLocation(problem.NewLocationWithID("startLoc")).Build()
	r := NewVehicleRouteBuilder(v, driver.NewNoDriver()).AddService(delivery).Build()

	act := r.Activities()[0]
	assert.Equal(t, "delivery", act.Name())
	assert.IsType(t, &activity.DeliverService{}, act)
}
