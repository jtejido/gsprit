package problem

import (
	"fmt"
	"gsprit/util"
)

type LocationBuilder struct {
	index      int
	id         string
	coordinate *util.Coordinate
	name       string
	userData   any
}

func NewLocationBuilder() *LocationBuilder {
	return &LocationBuilder{
		index: NoIndex,
	}
}

func (b *LocationBuilder) SetUserData(userData any) *LocationBuilder {
	b.userData = userData
	return b
}

func (b *LocationBuilder) SetIndex(index int) *LocationBuilder {
	if index < 0 {
		panic("index must be >= 0")
	}
	b.index = index
	return b
}

func (b *LocationBuilder) SetCoordinate(coordinate *util.Coordinate) *LocationBuilder {
	b.coordinate = coordinate
	return b
}

func (b *LocationBuilder) SetId(id string) *LocationBuilder {
	b.id = id
	return b
}

func (b *LocationBuilder) SetName(name string) *LocationBuilder {
	b.name = name
	return b
}

func (b *LocationBuilder) Build() *Location {
	if b.id == "" && b.coordinate == nil {
		if b.index == -1 {
			panic("coordinate or index must be set")
		}
	}
	if b.coordinate != nil && b.id == "" {
		b.id = b.coordinate.String()
	}
	if b.index != -1 && b.id == "" {
		b.id = fmt.Sprintf("%d", b.index)
	}
	return &Location{
		index:      b.index,
		id:         b.id,
		name:       b.name,
		coordinate: b.coordinate,
		userData:   b.userData,
	}
}

type Location struct {
	index      int
	coordinate *util.Coordinate
	id         string
	name       string
	userData   any
}

const NoIndex = -1

// NewLocation creates a new location instance with coordinate
func NewLocationWithCoordinate(x, y float64) *Location {
	return NewLocationBuilder().SetCoordinate(util.NewCoordinate(x, y)).Build()
}

// NewLocationWithID creates a new location instance with ID
func NewLocationWithID(id string) *Location {
	return NewLocationBuilder().SetId(id).Build()
}

// NewLocationWithIndex creates a new location instance with Index
func NewLocationWithIndex(index int) *Location {
	return NewLocationBuilder().SetIndex(index).Build()
}

func (l *Location) Id() string {
	return l.id
}

func (l *Location) Index() int {
	return l.index
}

func (l *Location) Coordinate() *util.Coordinate {
	return l.coordinate
}

func (l *Location) Name() string {
	return l.name
}

func (l *Location) SetUserData(data interface{}) {
	l.userData = data
}

func (l *Location) UserData() interface{} {
	return l.userData
}

func (l *Location) String() string {
	return fmt.Sprintf("[id=%s][index=%d][coordinate=%v]", l.id, l.index, l.coordinate)
}
