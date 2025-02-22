package problem

type HasID interface {
	Id() string
}

type HasIndex interface {
	Index() int
}

type TimeWindow interface {
	Larger(TimeWindow) bool
	Start() float64
	End() float64
	String() string
}

type Coordinate interface {
	String() string
	Equals(other any) bool
}
