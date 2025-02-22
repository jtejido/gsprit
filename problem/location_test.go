package problem

import (
	"testing"

	"gsprit/util"

	"github.com/stretchr/testify/assert"
)

func TestWhenIndexSet_BuildLocation(t *testing.T) {
	l := NewLocationBuilder().SetIndex(1).Build()
	assert.Equal(t, 1, l.Index())
}

func TestWhenNameSet_BuildLocation(t *testing.T) {
	l := NewLocationBuilder().SetName("mystreet 6a").SetIndex(1).Build()
	assert.Equal(t, "mystreet 6a", l.Name())
}

func TestWhenIndexSetWithFactory_ReturnCorrectLocation(t *testing.T) {
	l := NewLocationWithIndex(1)
	assert.Equal(t, 1, l.Index())
}

func TestWhenIndexSmallerZero_ThrowException(t *testing.T) {
	assert.Panics(t, func() {
		NewLocationBuilder().SetIndex(-1).Build()
	})
}

func TestWhenCoordinateAndIdAndIndexNotSet_ThrowException(t *testing.T) {
	assert.Panics(t, func() {
		NewLocationBuilder().Build()
	})
}

func TestWhenIdSet_Build(t *testing.T) {
	l := NewLocationBuilder().SetId("id").Build()
	assert.Equal(t, "id", l.Id())
}

func TestWhenIdSetWithFactory_ReturnCorrectLocation(t *testing.T) {
	l := NewLocationWithID("id")
	assert.Equal(t, "id", l.Id())
}

func TestWhenCoordinateSet_Build(t *testing.T) {
	l := NewLocationBuilder().SetCoordinate(util.NewCoordinate(10, 20)).Build()
	assert.NotNil(t, l.Coordinate())
	assert.Equal(t, 10.0, l.Coordinate().X)
	assert.Equal(t, 20.0, l.Coordinate().Y)
}

func TestWhenCoordinateSetWithFactory_ReturnCorrectLocation(t *testing.T) {
	l := NewLocationWithCoordinate(10, 20)
	assert.NotNil(t, l.Coordinate())
	assert.Equal(t, 10.0, l.Coordinate().X)
	assert.Equal(t, 20.0, l.Coordinate().Y)
}

func TestWhenSettingUserData_ItIsAssociatedWithTheLocation(t *testing.T) {
	one := NewLocationBuilder().
		SetCoordinate(util.NewCoordinate(10, 20)).
		SetUserData(map[string]any{}).
		Build()

	two := NewLocationBuilder().
		SetIndex(1).
		SetUserData(42).
		Build()

	three := NewLocationBuilder().SetIndex(2).Build()

	assert.IsType(t, map[string]any{}, one.UserData())
	assert.Equal(t, 42, two.UserData())
	assert.Nil(t, three.UserData())
}
