package vrp

import (
	"gsprit/problem"
	"gsprit/problem/cost"
	"gsprit/problem/driver"
	"gsprit/problem/job"
	gmock "gsprit/problem/mock"
	"gsprit/problem/solution/route"
	"gsprit/problem/vehicle"
	"gsprit/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createMockPickup(t *testing.T) *gmock.MockPickup {
	p := gmock.NewMockPickup(t)
	p.On("JobType").Return(problem.JobTypePickupService)
	p.On("Activities").Return([]problem.Activity{})
	p.On("SetIndex", mock.Anything).Return()
	return p
}

func createMockDelivery(t *testing.T) *gmock.MockDelivery {
	d := gmock.NewMockDelivery(t)
	d.On("JobType").Return(problem.JobTypeDeliveryService)
	d.On("Activities").Return([]problem.Activity{})
	d.On("SetIndex", mock.Anything).Return()
	return d
}

func createMockService(t *testing.T) *gmock.MockService {
	d := gmock.NewMockService(t)
	d.On("JobType").Return(problem.JobTypeService)
	d.On("Activities").Return([]problem.Activity{})
	d.On("SetIndex", mock.Anything).Return()
	return d
}

func TestNewBuilder_DefaultFleetSizeIsInfinite(t *testing.T) {
	builder := NewBuilder()
	assert.Equal(t, Infinite, builder.Build().FleetSize())
}

func TestBuilder_SetFleetSize_Finite(t *testing.T) {
	builder := NewBuilder()
	builder.SetFleetSize(Finite)
	assert.Equal(t, Finite, builder.Build().FleetSize())
}

func TestBuilder_AddMultipleVehicles_ShouldContainCorrectNumberOfVehicles(t *testing.T) {
	builder := NewBuilder()

	v1 := vehicle.NewVehicleBuilder("v1").SetStartLocation(problem.NewLocationWithID("start")).Build()
	v2 := vehicle.NewVehicleBuilder("v2").SetStartLocation(problem.NewLocationWithID("start")).Build()
	v3 := vehicle.NewVehicleBuilder("v3").SetStartLocation(problem.NewLocationWithID("start")).Build()
	v4 := vehicle.NewVehicleBuilder("v4").SetStartLocation(problem.NewLocationWithID("start")).Build()

	builder.AddVehicle(v1).
		AddVehicle(v2).
		AddVehicle(v3).
		AddVehicle(v4)

	vrpInstance := builder.Build()
	assert.Len(t, vrpInstance.Vehicles(), 4)
}

func TestBuilder_AddAllVehiclesAtOnce_ShouldContainCorrectNumberOfVehicles(t *testing.T) {
	builder := NewBuilder()

	v1 := vehicle.NewVehicleBuilder("v1").SetStartLocation(problem.NewLocationWithID("start")).Build()
	v2 := vehicle.NewVehicleBuilder("v2").SetStartLocation(problem.NewLocationWithID("start")).Build()
	v3 := vehicle.NewVehicleBuilder("v3").SetStartLocation(problem.NewLocationWithID("start")).Build()
	v4 := vehicle.NewVehicleBuilder("v4").SetStartLocation(problem.NewLocationWithID("start")).Build()

	builder.AddAllVehicles([]problem.Vehicle{v1, v2, v3, v4})

	vrpInstance := builder.Build()
	assert.Len(t, vrpInstance.Vehicles(), 4)
}

func TestBuilder_AddVehiclesWithTypes_ShouldContainCorrectNumberOfTypes(t *testing.T) {
	builder := NewBuilder()

	type1 := vehicle.NewVehicleTypeBuilder("type1").Build()
	type2 := vehicle.NewVehicleTypeBuilder("type2").Build()

	v1 := vehicle.NewVehicleBuilder("v1").SetStartLocation(problem.NewLocationWithID("yo")).SetType(type1).Build()
	v2 := vehicle.NewVehicleBuilder("v2").SetStartLocation(problem.NewLocationWithID("yo")).SetType(type1).Build()
	v3 := vehicle.NewVehicleBuilder("v3").SetStartLocation(problem.NewLocationWithID("yo")).SetType(type2).Build()
	v4 := vehicle.NewVehicleBuilder("v4").SetStartLocation(problem.NewLocationWithID("yo")).SetType(type2).Build()

	builder.AddVehicle(v1).
		AddVehicle(v2).
		AddVehicle(v3).
		AddVehicle(v4)

	vrpInstance := builder.Build()
	assert.Len(t, vrpInstance.Types(), 2)
}

func TestBuilder_ShipmentsAreAdded_ShouldContainThem(t *testing.T) {
	s := job.NewShipmentBuilder("s").AddSizeDimension(0, 10).SetPickupLocation(problem.NewLocationBuilder().SetId("foofoo").Build()).SetDeliveryLocation(problem.NewLocationWithID("foo")).Build()
	s2 := job.NewShipmentBuilder("s2").AddSizeDimension(0, 100).SetPickupLocation(problem.NewLocationBuilder().SetId("foofoo").Build()).SetDeliveryLocation(problem.NewLocationWithID("foo")).Build()
	vrpBuilder := NewBuilder()
	vrpBuilder.AddJob(s)
	vrpBuilder.AddJob(s2)
	vrp := vrpBuilder.Build()
	assert.Equal(t, 2, len(vrp.Jobs()))
	assert.Equal(t, s, vrp.Jobs()["s"])
	assert.Equal(t, s2, vrp.Jobs()["s2"])
	assert.Equal(t, 2, len(vrp.AllLocations()))
}

func TestBuilder_ServicesWithNoLocationAreAdded_ShouldKnowThat(t *testing.T) {
	s1 := job.NewServiceBuilder[*job.Service]("s1").SetLocation(problem.NewLocationBuilder().SetIndex(1).Build()).Build()
	s2 := job.NewServiceBuilder[*job.Service]("s2").Build()
	vrpBuilder := NewBuilder()
	vrpBuilder.AddJob(s1)
	vrpBuilder.AddJob(s2)
	vrp := vrpBuilder.Build()
	assert.Equal(t, 2, len(vrp.Jobs()))
	assert.Equal(t, s1, vrp.Jobs()["s1"])
	assert.Equal(t, s2, vrp.Jobs()["s2"])
	assert.Equal(t, 1, len(vrp.AllLocations()))
	assert.Equal(t, 1, len(vrp.JobsWithLocation()))
}

func TestBuilder_ServicesAreAdded_ShouldContainThem(t *testing.T) {
	s1 := job.NewServiceBuilder[*job.Service]("s1").SetLocation(problem.NewLocationBuilder().SetIndex(1).Build()).Build()
	s2 := job.NewServiceBuilder[*job.Service]("s2").SetLocation(problem.NewLocationBuilder().SetIndex(1).Build()).Build()
	vrpBuilder := NewBuilder()
	vrpBuilder.AddJob(s1)
	vrpBuilder.AddJob(s2)
	vrp := vrpBuilder.Build()
	assert.Equal(t, 2, len(vrp.Jobs()))
	assert.Equal(t, s1, vrp.Jobs()["s1"])
	assert.Equal(t, s2, vrp.Jobs()["s2"])
	assert.Equal(t, 1, len(vrp.AllLocations()))
}

func TestBuilder_PickupsAreAdded_ShouldContainThem(t *testing.T) {
	s1 := createMockPickup(t)
	s1.On("Id").Return("s1")
	s1.On("Location").Return(problem.NewLocationBuilder().SetIndex(1).Build()).Maybe()

	s2 := createMockPickup(t)
	s2.On("Id").Return("s2")
	s2.On("Location").Return(problem.NewLocationBuilder().SetIndex(1).Build()).Maybe()

	vrpBuilder := NewBuilder()
	vrpBuilder.AddJob(s1).AddJob(s2)
	vrp := vrpBuilder.Build()

	// Assertions
	assert.Equal(t, 2, len(vrp.Jobs()))
	assert.Equal(t, s1, vrp.Jobs()["s1"])
	assert.Equal(t, s2, vrp.Jobs()["s2"])

	// Verify mock expectations
	s1.AssertExpectations(t)
	s2.AssertExpectations(t)
}

func TestBuilder_DelivieriesAreAdded_ShouldContainThem(t *testing.T) {
	s1 := createMockDelivery(t)
	s1.On("Id").Return("s1")
	s1.On("Size").Return(problem.NewCapacityBuilder().Build()).Maybe()
	s1.On("Location").Return(problem.NewLocationBuilder().SetIndex(1).Build()).Maybe()

	s2 := createMockPickup(t)
	s2.On("Id").Return("s2")
	s2.On("Size").Return(problem.NewCapacityBuilder().Build()).Maybe()
	s2.On("Location").Return(problem.NewLocationBuilder().SetIndex(1).Build()).Maybe()

	vrpBuilder := NewBuilder()
	vrpBuilder.AddJob(s1).AddJob(s2)
	vrp := vrpBuilder.Build()

	// Assertions
	assert.Equal(t, 2, len(vrp.Jobs()))
	assert.Equal(t, s1, vrp.Jobs()["s1"])
	assert.Equal(t, s2, vrp.Jobs()["s2"])

	// Verify mock expectations
	s1.AssertExpectations(t)
	s2.AssertExpectations(t)
}

func TestBuilder_DelivieriesAreAddedAllAtOnce_ShouldContainThem(t *testing.T) {
	s1 := createMockDelivery(t)
	s1.On("Id").Return("s1")
	s1.On("Size").Return(problem.NewCapacityBuilder().Build()).Maybe()
	s1.On("Location").Return(problem.NewLocationBuilder().SetIndex(1).Build()).Maybe()

	s2 := createMockPickup(t)
	s2.On("Id").Return("s2")
	s2.On("Size").Return(problem.NewCapacityBuilder().Build()).Maybe()
	s2.On("Location").Return(problem.NewLocationBuilder().SetIndex(1).Build()).Maybe()

	vrpBuilder := NewBuilder()
	vrpBuilder.AddAllJobs([]problem.Job{s1, s2})
	vrp := vrpBuilder.Build()

	// Assertions
	assert.Equal(t, 2, len(vrp.Jobs()))
	assert.Equal(t, s1, vrp.Jobs()["s1"])
	assert.Equal(t, s2, vrp.Jobs()["s2"])

	// Verify mock expectations
	s1.AssertExpectations(t)
	s2.AssertExpectations(t)
}

func TestBuilder_ServicesAreAddedAllAtOnce_ShouldContainThem(t *testing.T) {
	s1 := createMockService(t)
	s1.On("Id").Return("s1")
	s1.On("Location").Return(problem.NewLocationBuilder().SetIndex(1).Build()).Maybe()

	s2 := createMockService(t)
	s2.On("Id").Return("s2")
	s2.On("Location").Return(problem.NewLocationBuilder().SetIndex(1).Build()).Maybe()

	vrpBuilder := NewBuilder()
	vrpBuilder.AddAllJobs([]problem.Job{s1, s2})
	vrp := vrpBuilder.Build()

	// Assertions
	assert.Equal(t, 2, len(vrp.Jobs()))
	assert.Equal(t, s1, vrp.Jobs()["s1"])
	assert.Equal(t, s2, vrp.Jobs()["s2"])

	// Verify mock expectations
	s1.AssertExpectations(t)
	s2.AssertExpectations(t)
}

type testVehicleRoutingActivityCosts struct{}

func (t *testVehicleRoutingActivityCosts) ActivityCost(tourAct problem.TourActivity, arrivalTime float64, driver problem.Driver, vehicle problem.Vehicle) float64 {
	return 4.0
}

func (t *testVehicleRoutingActivityCosts) ActivityDuration(tourAct problem.TourActivity, arrivalTime float64, driver problem.Driver, vehicle problem.Vehicle) float64 {
	return tourAct.OperationTime()
}

func TestBuilder_SettingActivityCosts_ShouldContainIt(t *testing.T) {
	vrpBuilder := NewBuilder()
	vrpBuilder.SetActivityCosts(new(testVehicleRoutingActivityCosts))
	problem := vrpBuilder.Build()
	assert.Equal(t, 4.0, problem.ActivityCosts().ActivityCost(nil, 0.0, nil, nil))
}

type testAbstractForwardVehicleRoutingTransportCosts struct {
	cost.AbstractForwardVehicleRoutingTransportCosts
}

func newTestRoutingCosts() *testAbstractForwardVehicleRoutingTransportCosts {
	res := &testAbstractForwardVehicleRoutingTransportCosts{}
	res.Spi = res
	return res
}

func (t *testAbstractForwardVehicleRoutingTransportCosts) Distance(from, to *problem.Location, departureTime float64, vehicle problem.Vehicle) float64 {
	return 0
}
func (t *testAbstractForwardVehicleRoutingTransportCosts) TransportTime(from, to *problem.Location, departureTime float64, driver problem.Driver, vehicle problem.Vehicle) float64 {
	return 0
}
func (t *testAbstractForwardVehicleRoutingTransportCosts) TransportCost(from, to *problem.Location, departureTime float64, driver problem.Driver, vehicle problem.Vehicle) float64 {
	return 4.0
}

func (t *testAbstractForwardVehicleRoutingTransportCosts) String() string {
	return "test"
}

func TestBuilder_SettingRoutingCosts_ShouldContainIt(t *testing.T) {
	vrpBuilder := NewBuilder()
	vrpBuilder.SetRoutingCost(newTestRoutingCosts())
	problem := vrpBuilder.Build()
	assert.Equal(t, 4.0, problem.TransportCosts().TransportCost(loc("a"), loc("b"), 0.0, nil, nil))
}

func loc(i string) *problem.Location {
	return problem.NewLocationBuilder().SetId(i).Build()
}

func TestBuilder_AddingVehiclesWithSameId_ItShouldThrowException(t *testing.T) {
	builder := NewBuilder()
	tt := vehicle.NewVehicleTypeBuilder("type").Build()
	v1 := vehicle.NewVehicleBuilder("v").SetStartLocation(problem.NewLocationWithID("loc")).SetType(tt).Build()
	v2 := vehicle.NewVehicleBuilder("v").SetStartLocation(problem.NewLocationWithID("loc")).SetType(tt).Build()

	assert.Panics(t, func() {
		builder.AddVehicle(v1)
		builder.AddVehicle(v2)
	})
}

func TestBuilder_AddingVehicleTypesWithSameIdButDifferentCosts_ItShouldThrowException(t *testing.T) {
	builder := NewBuilder()
	type1 := vehicle.NewVehicleTypeBuilder("type").Build()
	type2 := vehicle.NewVehicleTypeBuilder("type").Build()
	v1 := vehicle.NewVehicleBuilder("v").SetStartLocation(problem.NewLocationWithID("loc")).SetType(type1).Build()
	v2 := vehicle.NewVehicleBuilder("v").SetStartLocation(problem.NewLocationWithID("loc")).SetType(type2).Build()

	assert.Panics(t, func() {
		builder.AddVehicle(v1)
		builder.AddVehicle(v2)
	})
}

func TestBuilder_BuildingProblemWithSameBreakId_ItShouldThrowException(t *testing.T) {
	builder := NewBuilder()
	tt := vehicle.NewVehicleTypeBuilder("type").Build()

	v1 := vehicle.NewVehicleBuilder("v1").SetStartLocation(problem.NewLocationWithID("loc")).SetType(tt).SetBreak(job.NewBreakBuilder("break").Build()).Build()
	v2 := vehicle.NewVehicleBuilder("v2").SetStartLocation(problem.NewLocationWithID("loc")).SetType(tt).SetBreak(job.NewBreakBuilder("break").Build()).Build()

	builder.AddVehicle(v1)
	builder.AddVehicle(v2)
	builder.SetFleetSize(Finite)
	assert.Panics(t, func() {
		builder.Build()
	})
}

func TestBuilder_AddingAVehicle_AddedVehicleTypesShouldReturnItsType(t *testing.T) {
	builder := NewBuilder()
	tt := vehicle.NewVehicleTypeBuilder("type").Build()
	v := vehicle.NewVehicleBuilder("v1").SetStartLocation(problem.NewLocationWithID("loc")).SetType(tt).Build()
	builder.AddVehicle(v)

	assert.Equal(t, 1, len(builder.AddedVehicleTypes()))
	assert.Equal(t, tt, builder.AddedVehicleTypes()[0])
}

func TestBuilder_AddingTwoVehicleWithSameType_AddedVehicleTypesShouldReturnOnlyOneType(t *testing.T) {
	builder := NewBuilder()
	tt := vehicle.NewVehicleTypeBuilder("type").Build()
	v1 := vehicle.NewVehicleBuilder("v1").SetStartLocation(problem.NewLocationWithID("loc")).SetType(tt).Build()
	v2 := vehicle.NewVehicleBuilder("v2").SetStartLocation(problem.NewLocationWithID("loc")).SetType(tt).Build()
	builder.AddVehicle(v1).AddVehicle(v2)

	assert.Equal(t, 1, len(builder.AddedVehicleTypes()))
	assert.Equal(t, tt, builder.AddedVehicleTypes()[0])
}

func TestBuilder_AddingTwoVehicleWithDiffType_AddedVehicleTypesShouldReturnTheseType(t *testing.T) {
	builder := NewBuilder()
	type1 := vehicle.NewVehicleTypeBuilder("type").Build()
	type2 := vehicle.NewVehicleTypeBuilder("type2").Build()
	v1 := vehicle.NewVehicleBuilder("v1").SetStartLocation(problem.NewLocationWithID("loc")).SetType(type1).Build()
	v2 := vehicle.NewVehicleBuilder("v2").SetStartLocation(problem.NewLocationWithID("loc")).SetType(type2).Build()
	builder.AddVehicle(v1).AddVehicle(v2)

	assert.Equal(t, 2, len(builder.AddedVehicleTypes()))
}

func TestBuilder_AddingVehicleWithDiffStartAndEnd_StartLocationMustBeRegisteredInLocationMap(t *testing.T) {
	builder := NewBuilder()
	v1 := vehicle.NewVehicleBuilder("v1").SetStartLocation(problem.NewLocationWithID("start")).SetEndLocation(problem.NewLocationWithID("end")).Build()
	builder.AddVehicle(v1)
	_, exists := builder.LocationMap()["start"]
	assert.True(t, exists)
}

func TestBuilder_AddingVehicleWithDiffStartAndEnd_EndLocationMustBeRegisteredInLocationMap(t *testing.T) {
	builder := NewBuilder()
	v1 := vehicle.NewVehicleBuilder("v1").SetStartLocation(problem.NewLocationWithID("start")).SetEndLocation(problem.NewLocationWithID("end")).Build()
	builder.AddVehicle(v1)
	_, exists := builder.LocationMap()["end"]
	assert.True(t, exists)
}

func TestBuilder_AddingInitialRoute_ItShouldBeAddedCorrectly(t *testing.T) {
	builder := NewBuilder()
	v := vehicle.NewVehicleBuilder("v").SetStartLocation(problem.NewLocationWithID("start")).SetEndLocation(problem.NewLocationWithID("end")).Build()
	r := route.NewVehicleRouteBuilder(v, driver.NewNoDriver()).Build()
	builder.AddInitialVehicleRoute(r)
	vrp := builder.Build()
	assert.True(t, len(vrp.InitialVehicleRoutes()) != 0)
}

func TestBuilder_AddingInitialRoute_TheyShouldBeAddedCorrectly(t *testing.T) {
	builder := NewBuilder()
	v1 := vehicle.NewVehicleBuilder("v").SetStartLocation(problem.NewLocationWithID("start")).SetEndLocation(problem.NewLocationWithID("end")).Build()
	r1 := route.NewVehicleRouteBuilder(v1, driver.NewNoDriver()).Build()

	v2 := vehicle.NewVehicleBuilder("v").SetStartLocation(problem.NewLocationWithID("start")).SetEndLocation(problem.NewLocationWithID("end")).Build()
	r2 := route.NewVehicleRouteBuilder(v2, driver.NewNoDriver()).Build()
	builder.AddInitialVehicleRoutes([]*route.VehicleRoute{r1, r2})
	vrp := builder.Build()
	assert.Equal(t, 2, len(vrp.InitialVehicleRoutes()))
	assert.Equal(t, 2, len(vrp.AllLocations()))
}

func TestBuilder_AddingInitialRoute_LocationOfVehicleMustBeMemorized(t *testing.T) {
	start := problem.NewLocationBuilder().SetId("start").SetCoordinate(util.NewCoordinate(0, 1)).Build()
	end := problem.NewLocationWithID("end")
	builder := NewBuilder()
	v := vehicle.NewVehicleBuilder("v").
		SetStartLocation(start).
		SetEndLocation(end).
		Build()
	r := route.NewVehicleRouteBuilder(v, driver.NewNoDriver()).Build()
	builder.AddInitialVehicleRoute(r)
	vrp := builder.Build()
	assert.Contains(t, vrp.AllLocations(), start)
	assert.Contains(t, vrp.AllLocations(), end)
}

func TestBuilder_AddingJobAndInitialRouteWithThatJobAfterwards_ThisJobShouldNotBeInFinalJobMap(t *testing.T) {
	service := job.NewServiceBuilder[*job.Service]("myService").SetLocation(problem.NewLocationWithID("loc")).Build()
	builder := NewBuilder()
	builder.AddJob(service)
	vehicle := vehicle.NewVehicleBuilder("v").
		SetStartLocation(problem.NewLocationBuilder().SetId("start").SetCoordinate(util.NewCoordinate(0, 1)).Build()).
		SetEndLocation(problem.NewLocationWithID("end")).Build()
	initialRoute := route.NewVehicleRouteBuilder(vehicle, driver.NewNoDriver()).AddService(service).Build()
	builder.AddInitialVehicleRoute(initialRoute)
	vrp := builder.Build()
	assert.NotContains(t, vrp.Jobs(), "myService")
	assert.Equal(t, 3, len(vrp.AllLocations()))
}

func TestBuilder_AddingTwoJobs_TheyShouldHaveProperIndeces(t *testing.T) {
	service := job.NewServiceBuilder[*job.Service]("myService").SetLocation(problem.NewLocationWithID("loc")).Build()
	shipment := job.NewShipmentBuilder("shipment").SetPickupLocation(problem.NewLocationBuilder().SetId("pick").Build()).
		SetDeliveryLocation(problem.NewLocationWithID("del")).Build()
	builder := NewBuilder()
	builder.AddJob(service)
	builder.AddJob(shipment)
	vrp := builder.Build()

	assert.Equal(t, 1, service.Index())
	assert.Equal(t, 2, shipment.Index())
	assert.Equal(t, 3, len(vrp.AllLocations()))
}

func TestBuilder_AddingTwoServicesWithTheSameId_ItShouldThrowException(t *testing.T) {
	service1 := job.NewServiceBuilder[*job.Service]("myService").SetLocation(problem.NewLocationWithID("loc")).Build()
	service2 := job.NewServiceBuilder[*job.Service]("myService").SetLocation(problem.NewLocationWithID("loc")).Build()
	builder := NewBuilder()

	assert.Panics(t, func() {
		builder.AddJob(service1).AddJob(service2)
	})
}

func TestBuilder_AddingTwoShipmentsWithTheSameId_ItShouldThrowException(t *testing.T) {
	shipment1 := job.NewShipmentBuilder("shipment").SetPickupLocation(problem.NewLocationBuilder().SetId("pick").Build()).
		SetDeliveryLocation(problem.NewLocationWithID("del")).Build()
	shipment2 := job.NewShipmentBuilder("shipment").SetPickupLocation(problem.NewLocationBuilder().SetId("pick").Build()).
		SetDeliveryLocation(problem.NewLocationWithID("del")).Build()
	builder := NewBuilder()

	assert.Panics(t, func() {
		builder.AddJob(shipment1).AddJob(shipment2)
	})
}

func TestBuilder_AddingTwoVehicles_TheyShouldHaveProperIndices(t *testing.T) {
	vehicle1 := vehicle.NewVehicleBuilder("v1").
		SetStartLocation(problem.NewLocationBuilder().SetId("start").SetCoordinate(util.NewCoordinate(0, 1)).Build()).
		SetEndLocation(problem.NewLocationWithID("end")).Build()
	vehicle2 := vehicle.NewVehicleBuilder("v2").
		SetStartLocation(problem.NewLocationBuilder().SetId("start").SetCoordinate(util.NewCoordinate(0, 1)).Build()).
		SetEndLocation(problem.NewLocationWithID("end")).Build()
	builder := NewBuilder()
	builder.AddVehicle(vehicle1).AddVehicle(vehicle2)
	builder.Build()

	assert.Equal(t, 1, vehicle1.Index())
	assert.Equal(t, 2, vehicle2.Index())
}

func TestBuilder_AddingTwoVehiclesWithSameTypeIdentifier_TypeIdentifiersShouldHaveSameIndices(t *testing.T) {
	vehicle1 := vehicle.NewVehicleBuilder("v1").
		SetStartLocation(problem.NewLocationBuilder().SetId("start").SetCoordinate(util.NewCoordinate(0, 1)).Build()).
		SetEndLocation(problem.NewLocationWithID("end")).Build()
	vehicle2 := vehicle.NewVehicleBuilder("v2").
		SetStartLocation(problem.NewLocationBuilder().SetId("start").SetCoordinate(util.NewCoordinate(0, 1)).Build()).
		SetEndLocation(problem.NewLocationWithID("end")).Build()
	builder := NewBuilder()
	builder.AddVehicle(vehicle1)
	builder.AddVehicle(vehicle2)
	builder.Build()

	assert.Equal(t, 1, vehicle1.VehicleTypeIdentifier().Index())
	assert.Equal(t, 1, vehicle2.VehicleTypeIdentifier().Index())
}

func TestBuilder_AddingTwoVehiclesDifferentTypeIdentifier_TypeIdentifiersShouldHaveDifferentIndices(t *testing.T) {
	vehicle1 := vehicle.NewVehicleBuilder("v1").
		SetStartLocation(problem.NewLocationBuilder().SetId("start").SetCoordinate(util.NewCoordinate(0, 1)).Build()).
		SetEndLocation(problem.NewLocationWithID("end")).Build()
	vehicle2 := vehicle.NewVehicleBuilder("v2").
		SetStartLocation(problem.NewLocationBuilder().SetId("startLoc").SetCoordinate(util.NewCoordinate(0, 1)).Build()).
		SetEndLocation(problem.NewLocationWithID("end")).Build()
	builder := NewBuilder()
	builder.AddVehicle(vehicle1)
	builder.AddVehicle(vehicle2)
	builder.Build()

	assert.Equal(t, 1, vehicle1.VehicleTypeIdentifier().Index())
	assert.Equal(t, 2, vehicle2.VehicleTypeIdentifier().Index())
}
