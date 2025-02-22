package problem

type TourActivity interface {
	HasIndex
	SetTheoreticalEarliestOperationStartTime(earliest float64)
	SetTheoreticalLatestOperationStartTime(latest float64)
	Name() string
	Location() *Location
	TheoreticalEarliestOperationStartTime() float64
	TheoreticalLatestOperationStartTime() float64
	OperationTime() float64
	ArrTime() float64
	EndTime() float64
	SetArrTime(arrTime float64)
	SetEndTime(endTime float64)
	Size() *Capacity
	Duplicate() TourActivity
	String() string
}

type TourActivities interface {
	Activities() []TourActivity
	IsEmpty() bool
	Jobs() []Job
	ServesJob(job Job) bool
	RemoveJob(job Job) bool
	RemoveActivity(activity TourActivity) bool
	AddActivity(index int, act TourActivity) error
	AddActivityToEnd(act TourActivity) error
	JobSize() int
	String() string
	Copy() TourActivities
}

type JobActivity interface {
	TourActivity
	Job() Job
}

type AbstractActivity interface {
	TourActivity
	SetIndex(index int)
}

type BaseActivity struct {
	index int
}

func (a *BaseActivity) Index() int {
	return a.index
}

func (a *BaseActivity) SetIndex(index int) {
	a.index = index
}
