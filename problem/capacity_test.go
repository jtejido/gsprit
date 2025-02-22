package problem

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSettingOneCapacityDimension_NuOfDimensionsMustBeCorrect(t *testing.T) {
	cb := NewCapacityBuilder()
	cb.AddDimension(0, 4)
	cap := cb.Build()
	assert.Equal(t, 1, cap.NuOfDimensions())
}

func TestSettingTwoCapacityDimensions_NuOfDimensionsMustBeCorrect(t *testing.T) {
	cb := NewCapacityBuilder()
	cb.AddDimension(0, 4)
	cb.AddDimension(1, 10)
	cap := cb.Build()
	assert.Equal(t, 2, cap.NuOfDimensions())
}

func TestSettingRandomCapacityDimensions_NuOfDimensionsMustBeCorrect(t *testing.T) {
	nuOfCapDimensions := 1 + rand.Intn(100)
	cb := NewCapacityBuilder()
	cb.AddDimension(nuOfCapDimensions-1, 4)
	cap := cb.Build()
	assert.Equal(t, nuOfCapDimensions, cap.NuOfDimensions())
}

func TestSettingOneDimensionValue_ValueMustBeCorrect(t *testing.T) {
	cb := NewCapacityBuilder()
	cb.AddDimension(0, 4)
	cap := cb.Build()
	assert.Equal(t, 4, cap.Get(0))
}

func TestGettingIndexHigherThanCapacityDimensions_ShouldReturnZero(t *testing.T) {
	cb := NewCapacityBuilder()
	cb.AddDimension(0, 4)
	cap := cb.Build()
	assert.Equal(t, 0, cap.Get(2))
}

func TestSettingNoDimension_DefaultIsOneDimensionWithZeroValue(t *testing.T) {
	cb := NewCapacityBuilder()
	cap := cb.Build()
	assert.Equal(t, 1, cap.NuOfDimensions())
	assert.Equal(t, 0, cap.Get(0))
}

// func TestCopyingCapacityWithTwoDimensions_CopiedObjectShouldHaveSameDimensions(t *testing.T) {
// 	cb := NewCapacityBuilder()
// 	cb.AddDimension(0, 4)
// 	cb.AddDimension(1, 10)
// 	cap := cb.Build()

// 	copiedCapacity := NewCapacityFromBuilder(cb)
// 	assert.Equal(t, 2, copiedCapacity.NuOfDimensions())
// }

// func TestCopyingCapacityWithTwoDimensions_CopiedObjectShouldHaveSameValues(t *testing.T) {
// 	cb := NewCapacityBuilder()
// 	cb.AddDimension(0, 4)
// 	cb.AddDimension(1, 10)
// 	cap := cb.Build()

// 	copiedCapacity := NewCapacityFromBuilder(cb)
// 	assert.Equal(t, 4, copiedCapacity.Get(0))
// 	assert.Equal(t, 10, copiedCapacity.Get(1))
// }

func TestCopyingNull_ShouldReturnNil(t *testing.T) {
	var nullCap *Capacity = nil
	assert.Nil(t, nullCap)
}

func TestAddingTwoOneDimensionalCapacities_ShouldReturnCorrectCapacityValues(t *testing.T) {
	cap1 := NewCapacityBuilder().AddDimension(0, 1).Build()
	cap2 := NewCapacityBuilder().AddDimension(0, 2).Build()
	result := AddUp(cap1, cap2)
	assert.Equal(t, 3, result.Get(0))
}

func TestAddingTwoOneDimensionalCapacities_ShouldReturnCorrectNuOfDimensions(t *testing.T) {
	cap1 := NewCapacityBuilder().AddDimension(0, 1).Build()
	cap2 := NewCapacityBuilder().AddDimension(0, 2).Build()
	result := AddUp(cap1, cap2)
	assert.Equal(t, 1, result.NuOfDimensions())
}

func TestAddingTwoThreeDimensionalCapacities_ShouldReturnCorrectValues(t *testing.T) {
	cap1 := NewCapacityBuilder().AddDimension(0, 1).AddDimension(1, 2).AddDimension(2, 3).Build()
	cap2 := NewCapacityBuilder().AddDimension(0, 2).AddDimension(1, 3).AddDimension(2, 4).Build()
	result := AddUp(cap1, cap2)

	assert.Equal(t, 3, result.Get(0))
	assert.Equal(t, 5, result.Get(1))
	assert.Equal(t, 7, result.Get(2))
}

func TestAddingCapacitiesWithDifferentDimensions_ShouldAddCorrectly(t *testing.T) {
	cap1 := NewCapacityBuilder().AddDimension(0, 1).AddDimension(1, 2).Build()
	cap2 := NewCapacityBuilder().AddDimension(0, 2).Build()
	result := AddUp(cap1, cap2)

	assert.Equal(t, 3, result.Get(0))
	assert.Equal(t, 2, result.Get(1))
}

func TestSubtractingTwoCapacities_ShouldReturnCorrectValues(t *testing.T) {
	cap1 := NewCapacityBuilder().AddDimension(0, 1).Build()
	cap2 := NewCapacityBuilder().AddDimension(0, 2).Build()
	result := Subtract(cap2, cap1)

	assert.Equal(t, 1, result.Get(0))
}

func TestSubtractingCapacitiesWithDifferentDimensions_ShouldSubtractCorrectly(t *testing.T) {
	cap1 := NewCapacityBuilder().AddDimension(0, 1).AddDimension(1, 2).Build()
	cap2 := NewCapacityBuilder().AddDimension(0, 2).Build()
	result := Subtract(cap2, cap1)

	assert.Equal(t, 1, result.Get(0))
	assert.Equal(t, -2, result.Get(1))
}

func TestInvertingCapacity_ShouldBeDoneCorrectly(t *testing.T) {
	cap := NewCapacityBuilder().AddDimension(0, 2).AddDimension(1, 3).AddDimension(2, 4).Build()
	inverted := Invert(cap)

	assert.Equal(t, -2, inverted.Get(0))
	assert.Equal(t, -3, inverted.Get(1))
	assert.Equal(t, -4, inverted.Get(2))
}

func TestMaximumOfTwoCapacities_ShouldReturnMaxPerDimension(t *testing.T) {
	cap1 := NewCapacityBuilder().AddDimension(0, 3).AddDimension(1, 3).Build()
	cap2 := NewCapacityBuilder().AddDimension(0, 2).AddDimension(1, 4).Build()

	maxCap := NewCapacity([]int{max(cap1.Get(0), cap2.Get(0)), max(cap1.Get(1), cap2.Get(1))})
	assert.Equal(t, 3, maxCap.Get(0))
	assert.Equal(t, 4, maxCap.Get(1))
}

func TestDividingTwoCapacities_ShouldReturnCorrectRatio(t *testing.T) {
	cap1 := NewCapacityBuilder().AddDimension(0, 1).AddDimension(1, 2).Build()
	cap2 := NewCapacityBuilder().AddDimension(0, 2).AddDimension(1, 4).Build()

	div := float64(cap1.Get(0)) / float64(cap2.Get(0))
	assert.InDelta(t, 0.5, div, 0.001)
}

func TestEqualCapacities_ShouldReturnTrue(t *testing.T) {
	cap1 := NewCapacityBuilder().Build()
	cap2 := NewCapacityBuilder().Build()
	assert.Equal(t, cap1.String(), cap2.String())
}

func TestDifferentCapacities_ShouldReturnFalse(t *testing.T) {
	cap1 := NewCapacityBuilder().AddDimension(0, 10).AddDimension(1, 100).AddDimension(2, 1000).Build()
	cap2 := NewCapacityBuilder().AddDimension(0, 10).AddDimension(2, 1000).AddDimension(1, 100).Build()
	assert.Equal(t, cap1.String(), cap2.String()) // Should be true as same values
}
