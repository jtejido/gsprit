package vrp

import (
	"fmt"
	"gsprit/problem"
	"gsprit/problem/cost"
	"gsprit/problem/solution/route"
	"gsprit/problem/solution/route/activity"
	"gsprit/util"
	"log"
	"sort"
	"sync"
)

type FleetSize string

const (
	Finite   FleetSize = "FINITE"
	Infinite FleetSize = "INFINITE"
)

type JobActivityFactory interface {
	CreateActivities(job problem.Job) []problem.AbstractActivity
}

type defaultJobActivityFactory struct {
	serviceActivityFactory  *activity.DefaultTourActivityFactory
	shipmentActivityFactory *activity.DefaultShipmentActivityFactory
}

func NewDefaultJobActivityFactory() *defaultJobActivityFactory {
	return &defaultJobActivityFactory{
		//serviceActivityFactory:  serviceFactory,
		shipmentActivityFactory: new(activity.DefaultShipmentActivityFactory),
	}
}

func (f *defaultJobActivityFactory) CreateActivities(job problem.Job) []problem.AbstractActivity {
	var acts []problem.AbstractActivity

	if job.JobType().IsBreak() {
		acts = append(acts, activity.NewBreakActivity(job.(problem.Break)))
	} else if job.JobType().IsService() {
		acts = append(acts, f.serviceActivityFactory.CreateActivity(job.(problem.Service)))
	} else if job.JobType().IsShipment() {
		acts = append(acts, f.shipmentActivityFactory.CreatePickup(job.(problem.Shipment)))
		acts = append(acts, f.shipmentActivityFactory.CreateDelivery(job.(problem.Shipment)))
	}

	return acts
}

type Builder struct {
	sync.Mutex
	transportCosts                                                       cost.VehicleRoutingTransportCosts
	activityCosts                                                        cost.VehicleRoutingActivityCosts
	jobs                                                                 map[string]problem.Job
	jobsWithLocation                                                     []problem.Job
	tentativeJobs                                                        map[string]problem.Job
	jobsInInitialRoutes                                                  map[string]problem.Job
	tentativeCoordinates                                                 map[string]*util.Coordinate
	fleetSize                                                            FleetSize
	vehicleTypes                                                         map[string]problem.VehicleType
	initialRoutes                                                        []*route.VehicleRoute
	uniqueVehicles                                                       map[string]problem.Vehicle
	addedVehicleIds                                                      map[string]bool
	jobActivityFactory                                                   JobActivityFactory
	activityMap                                                          map[problem.Job][]problem.AbstractActivity
	allLocations                                                         map[string]*problem.Location
	vehicleIndexCounter, activityIndexCounter, vehicleTypeIdIndexCounter int
	typeKeyIndices                                                       map[string]int
	nonJobActivities                                                     []problem.AbstractActivity
}

func NewBuilder() *Builder {
	return &Builder{
		jobs:                      make(map[string]problem.Job),
		jobsWithLocation:          []problem.Job{},
		tentativeJobs:             make(map[string]problem.Job),
		jobsInInitialRoutes:       make(map[string]problem.Job),
		tentativeCoordinates:      make(map[string]*util.Coordinate),
		fleetSize:                 Infinite,
		vehicleTypes:              make(map[string]problem.VehicleType),
		initialRoutes:             []*route.VehicleRoute{},
		uniqueVehicles:            make(map[string]problem.Vehicle),
		addedVehicleIds:           make(map[string]bool),
		activityMap:               make(map[problem.Job][]problem.AbstractActivity),
		allLocations:              make(map[string]*problem.Location),
		vehicleIndexCounter:       1,
		activityIndexCounter:      1,
		vehicleTypeIdIndexCounter: 1,
		activityCosts:             new(cost.WaitingTimeCosts),
		jobActivityFactory:        NewDefaultJobActivityFactory(),
		typeKeyIndices:            make(map[string]int),
		nonJobActivities:          []problem.AbstractActivity{},
	}
}

func (b *Builder) incActivityIndexCounter() {
	b.activityIndexCounter++
}

func (b *Builder) incVehicleTypeIdIndexCounter() {
	b.vehicleTypeIdIndexCounter++
}

func (b *Builder) LocationMap() map[string]*util.Coordinate {
	return b.tentativeCoordinates
}

func (b *Builder) Locations() func(string) *util.Coordinate {
	return func(id string) *util.Coordinate {
		return b.tentativeCoordinates[id]
	}
}

func (b *Builder) SetRoutingCost(costs cost.VehicleRoutingTransportCosts) *Builder {
	b.transportCosts = costs
	return b
}

func (b *Builder) SetJobActivityFactory(jobActivityFactory JobActivityFactory) *Builder {
	b.jobActivityFactory = jobActivityFactory
	return b
}

func (b *Builder) SetActivityCosts(costs cost.VehicleRoutingActivityCosts) *Builder {
	b.activityCosts = costs
	return b
}

func (b *Builder) SetFleetSize(size FleetSize) *Builder {
	b.fleetSize = size
	return b
}

func (b *Builder) AddJob(job problem.Job) *Builder {
	if _, exists := b.tentativeJobs[job.Id()]; exists {
		panic(fmt.Sprintf("Job with ID %s already exists", job.Id()))
	}
	b.tentativeJobs[job.Id()] = job
	b.addLocationToTentativeLocationsFromJob(job)
	return b
}

func (b *Builder) addLocationToTentativeLocationsFromJob(job problem.Job) {
	for _, act := range job.Activities() {
		b.addLocationToTentativeLocations(act.Location())
	}
}

func (b *Builder) addLocationToTentativeLocations(location *problem.Location) {
	if location == nil {
		return
	}

	b.tentativeCoordinates[location.Id()] = location.Coordinate()
	b.allLocations[location.Id()] = location
}

func (b *Builder) AddVehicle(vehicle problem.Vehicle) *Builder {
	if _, exists := b.addedVehicleIds[vehicle.Id()]; exists {
		panic(fmt.Sprintf("Vehicle with ID %s already exists", vehicle.Id()))
	}

	b.addedVehicleIds[vehicle.Id()] = true
	if _, exists := b.uniqueVehicles[vehicle.Id()]; !exists {
		vehicle.(problem.AbstractVehicle).SetIndex(b.vehicleIndexCounter)
		b.incVehicleIndexCounter()
	}

	if v, exists := b.typeKeyIndices[vehicle.VehicleTypeIdentifier().String()]; exists {
		vehicle.VehicleTypeIdentifier().SetIndex(v)
	} else {
		vehicle.VehicleTypeIdentifier().SetIndex(b.vehicleTypeIdIndexCounter)
		b.typeKeyIndices[vehicle.VehicleTypeIdentifier().String()] = b.vehicleTypeIdIndexCounter
		b.incVehicleTypeIdIndexCounter()
	}
	b.uniqueVehicles[vehicle.Id()] = vehicle

	if v, exists := b.vehicleTypes[vehicle.Type().TypeId()]; !exists {
		b.vehicleTypes[vehicle.Type().TypeId()] = vehicle.Type()
	} else {
		if !vehicle.Type().Equals(v) {
			panic(fmt.Sprintf("A type with type id %s already exists. However, types are different. Please use unique vehicle types only.", vehicle.Type().TypeId()))
		}
	}

	startLocationId := vehicle.StartLocation().Id()
	b.addLocationToTentativeLocations(vehicle.StartLocation())

	if vehicle.EndLocation().Id() != startLocationId {
		b.addLocationToTentativeLocations(vehicle.EndLocation())
	}

	return b
}

func (b *Builder) incVehicleIndexCounter() {

	b.vehicleIndexCounter++
}

func (b *Builder) AddInitialVehicleRoute(route *route.VehicleRoute) *Builder {
	if _, ok := b.addedVehicleIds[route.Vehicle().Id()]; !ok {
		b.AddVehicle(route.Vehicle())
		b.addedVehicleIds[route.Vehicle().Id()] = true
	}
	for _, act := range route.Activities() {
		abstractAct, ok := act.(problem.AbstractActivity)
		if !ok {
			continue
		}
		abstractAct.SetIndex(b.activityIndexCounter)
		b.incActivityIndexCounter()
		if jobActivity, isJobActivity := act.(problem.JobActivity); isJobActivity {
			job := jobActivity.Job()
			b.jobsInInitialRoutes[job.Id()] = job
			b.addLocationToTentativeLocationsFromJob(job)
			b.registerJobAndActivity(abstractAct, job)
		}
	}
	b.initialRoutes = append(b.initialRoutes, route)
	return b
}

func (b *Builder) registerJobAndActivity(abstractAct problem.AbstractActivity, job problem.Job) {
	b.activityMap[job] = append(b.activityMap[job], abstractAct)
}

func (b *Builder) AddInitialVehicleRoutes(routes []*route.VehicleRoute) *Builder {
	for _, r := range routes {
		b.AddInitialVehicleRoute(r)
	}
	return b
}

func (b *Builder) AddNonJobActivities(nonJobActivities []problem.AbstractActivity) *Builder {
	for _, act := range nonJobActivities {
		act.SetIndex(b.activityIndexCounter)
		b.incActivityIndexCounter()
		b.nonJobActivities = append(b.nonJobActivities, act)
	}
	return b
}

func (b *Builder) AddLocation(locationId string, coordinate *util.Coordinate) *Builder {
	b.tentativeCoordinates[locationId] = coordinate
	return b
}

func (b *Builder) AddAllJobs(jobs []problem.Job) *Builder {
	for _, j := range jobs {
		b.AddJob(j)
	}
	return b
}

func (b *Builder) AddAllVehicles(vehicles []problem.Vehicle) *Builder {
	for _, v := range vehicles {
		b.AddVehicle(v)
	}
	return b
}

func (b *Builder) AddedVehicles() []problem.Vehicle {
	vehicles := []problem.Vehicle{}
	for _, v := range b.uniqueVehicles {
		vehicles = append(vehicles, v)
	}
	return vehicles
}

func (b *Builder) AddedVehicleTypes() []problem.VehicleType {
	types := []problem.VehicleType{}
	for _, v := range b.vehicleTypes {
		types = append(types, v)
	}
	return types
}

func (b *Builder) AddedJobs() []problem.Job {
	jobs := []problem.Job{}
	for _, v := range b.tentativeJobs {
		jobs = append(jobs, v)
	}
	return jobs
}

func (b *Builder) convertMapToSlice(vehicleMap map[string]problem.Vehicle) []problem.Vehicle {
	vehicles := []problem.Vehicle{}
	for _, v := range vehicleMap {
		vehicles = append(vehicles, v)
	}
	return vehicles
}

func (b *Builder) convertVehicleTypeMapToSlice() []problem.VehicleType {
	types := []problem.VehicleType{}
	for _, v := range b.vehicleTypes {
		types = append(types, v)
	}
	return types
}

func (b *Builder) addJobToFinalJobMapAndCreateActivities(job problem.Job) {
	b.addJobToFinalMap(job)
	jobActs := b.jobActivityFactory.CreateActivities(job)
	for _, act := range jobActs {
		act.SetIndex(b.activityIndexCounter)
		b.incActivityIndexCounter()
	}
	b.activityMap[job] = jobActs
}
func (b *Builder) addJobToFinalMap(job problem.Job) {
	if _, exists := b.jobs[job.Id()]; exists {
		log.Printf("The job %s has already been added to the job list. This overrides the existing job.", job.String())
	}
	b.addLocationToTentativeLocationsFromJob(job)
	b.jobs[job.Id()] = job
	hasLocation := true
	for _, activity := range job.Activities() {
		if activity.Location() == nil {
			hasLocation = false
		}
	}
	if hasLocation {
		b.jobsWithLocation = append(b.jobsWithLocation, job)
	}
}

func (b *Builder) addBreaksToActivityMap() bool {
	hasBreaks := false
	uniqueBreakIds := make(map[string]bool)

	for _, v := range b.uniqueVehicles {
		if v.Break() != nil {
			breakID := v.Break().Id()
			if uniqueBreakIds[breakID] {
				panic(fmt.Sprintf("The vehicle routing problem already contains a vehicle break with id %s. Please choose unique ids for each vehicle break.", breakID))
			}
			uniqueBreakIds[breakID] = true
			hasBreaks = true

			breakActivities := b.jobActivityFactory.CreateActivities(v.Break())
			if len(breakActivities) == 0 {
				panic("At least one activity for break needs to be created by activityFactory.")
			}

			for _, act := range breakActivities {
				act.SetIndex(b.activityIndexCounter)
				b.incActivityIndexCounter()
			}

			b.activityMap[v.Break()] = breakActivities
		}
	}
	return hasBreaks
}

func (b *Builder) Build() *VehicleRoutingProblem {
	if b.transportCosts == nil {
		b.transportCosts = cost.NewCrowFlyCosts(b.Locations())
	}

	for _, job := range b.tentativeJobs {
		if _, exists := b.jobsInInitialRoutes[job.Id()]; !exists {
			b.addJobToFinalJobMapAndCreateActivities(job)
		}
	}
	orderedJobs := make([]problem.Job, 0, len(b.jobs))
	for _, job := range b.jobs {
		orderedJobs = append(orderedJobs, job)
	}
	sort.Slice(orderedJobs, func(i, j int) bool {
		return orderedJobs[i].Id() < orderedJobs[j].Id() // Sort jobs by ID
	})
	jobIndexCounter := 1
	for _, job := range orderedJobs {
		job.(problem.AbstractJob).SetIndex(jobIndexCounter)
		jobIndexCounter++
	}
	for _, job := range b.jobsInInitialRoutes {
		job.(problem.AbstractJob).SetIndex(jobIndexCounter)
		jobIndexCounter++
	}

	hasBreaks := b.addBreaksToActivityMap()
	if hasBreaks && b.fleetSize == Infinite {
		panic("Breaks are not yet supported when dealing with infinite fleet. Either set it to finite or omit breaks.")
	}

	res := &VehicleRoutingProblem{
		transportCosts:       b.transportCosts,
		activityCosts:        b.activityCosts,
		jobs:                 b.jobs,
		jobsWithLocation:     b.jobsWithLocation,
		allJobs:              mergeMaps(b.jobs, b.jobsInInitialRoutes),
		vehicles:             b.convertMapToSlice(b.uniqueVehicles),
		vehicleTypes:         b.convertVehicleTypeMapToSlice(),
		initialVehicleRoutes: b.initialRoutes,
		allLocations:         b.convertLocationMapToSlice(),
		fleetSize:            b.fleetSize,
		activityMap:          b.activityMap,
		nuActivities:         b.activityIndexCounter,
	}
	res.jobActivityFactory = res.copyAndGetActivities
	return res
}

func mergeMaps[M ~map[K]V, K comparable, V any](src ...M) M {
	merged := make(M)
	for _, m := range src {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}

func (b *Builder) convertLocationMapToSlice() []*problem.Location {
	locations := []*problem.Location{}
	for _, v := range b.allLocations {
		locations = append(locations, v)
	}
	return locations
}

type VehicleRoutingProblem struct {
	sync.Mutex
	transportCosts       cost.VehicleRoutingTransportCosts
	activityCosts        cost.VehicleRoutingActivityCosts
	jobs                 map[string]problem.Job
	jobsWithLocation     []problem.Job
	allJobs              map[string]problem.Job
	vehicles             []problem.Vehicle
	vehicleTypes         []problem.VehicleType
	initialVehicleRoutes []*route.VehicleRoute
	allLocations         []*problem.Location
	fleetSize            FleetSize
	activityMap          map[problem.Job][]problem.AbstractActivity
	nuActivities         int
	jobActivityFactory   func(problem.Job) []problem.AbstractActivity
}

func (vrp *VehicleRoutingProblem) Jobs() map[string]problem.Job {
	copy := make(map[string]problem.Job)
	for k, v := range vrp.jobs {
		copy[k] = v
	}
	return copy
}

func (vrp *VehicleRoutingProblem) JobsWithLocation() []problem.Job {
	c := make([]problem.Job, len(vrp.jobsWithLocation))
	copy(c, vrp.jobsWithLocation)
	return c
}

func (vrp *VehicleRoutingProblem) JobsInclusiveInitialJobsInRoutes() map[string]problem.Job {
	copy := make(map[string]problem.Job)
	for k, v := range vrp.allJobs {
		copy[k] = v
	}
	return copy
}

func (vrp *VehicleRoutingProblem) InitialVehicleRoutes() []*route.VehicleRoute {
	copiedInitialRoutes := make([]*route.VehicleRoute, 0)
	for _, route := range vrp.initialVehicleRoutes {
		copiedInitialRoutes = append(copiedInitialRoutes, route.Copy())
	}
	return copiedInitialRoutes
}

func (vrp *VehicleRoutingProblem) Types() []problem.VehicleType {
	c := make([]problem.VehicleType, len(vrp.vehicleTypes))
	copy(c, vrp.vehicleTypes)
	return c
}

func (vrp *VehicleRoutingProblem) Vehicles() []problem.Vehicle {
	c := make([]problem.Vehicle, len(vrp.vehicles))
	copy(c, vrp.vehicles)
	return c
}

func (vrp *VehicleRoutingProblem) FleetSize() FleetSize {
	return vrp.fleetSize
}

func (vrp *VehicleRoutingProblem) TransportCosts() cost.VehicleRoutingTransportCosts {
	return vrp.transportCosts
}

func (vrp *VehicleRoutingProblem) ActivityCosts() cost.VehicleRoutingActivityCosts {
	return vrp.activityCosts
}

func (vrp *VehicleRoutingProblem) AllLocations() []*problem.Location {
	return vrp.allLocations
}

func (vrp *VehicleRoutingProblem) Activities(job problem.Job) []problem.AbstractActivity {
	c := make([]problem.AbstractActivity, len(vrp.activityMap[job]))
	copy(c, vrp.activityMap[job])
	return c
}

func (vrp *VehicleRoutingProblem) NuActivities() int {
	return vrp.nuActivities
}

func (vrp *VehicleRoutingProblem) JobActivityFactory() func(problem.Job) []problem.AbstractActivity {
	return vrp.jobActivityFactory
}

func (vrp *VehicleRoutingProblem) copyAndGetActivities(job problem.Job) []problem.AbstractActivity {
	vrp.Lock()
	defer vrp.Unlock()

	copiedActivities := []problem.AbstractActivity{}
	activities, exists := vrp.activityMap[job]
	if !exists {
		return copiedActivities
	}

	for _, act := range activities {
		copiedActivities = append(copiedActivities, act.Duplicate().(problem.AbstractActivity))
	}
	return copiedActivities
}

func (vrp *VehicleRoutingProblem) String() string {
	return fmt.Sprintf("[fleetSize=%s][#jobs=%d][#vehicles=%d][#vehicleTypes=%d][transportCost=%s][activityCosts=%v]",
		vrp.fleetSize, len(vrp.jobs), len(vrp.vehicles), len(vrp.vehicleTypes), vrp.transportCosts.String(), vrp.activityCosts)
}
