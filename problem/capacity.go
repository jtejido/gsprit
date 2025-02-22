package problem

import (
	"fmt"
)

type CapacityBuilder struct {
	dimensions []int
}

func NewCapacityBuilder() *CapacityBuilder {
	return &CapacityBuilder{
		dimensions: make([]int, 1),
	}
}

func (cb *CapacityBuilder) AddDimension(index, dimValue int) *CapacityBuilder {
	if index < len(cb.dimensions) {
		cb.dimensions[index] = dimValue
	} else {
		requiredSize := index + 1
		newDimensions := make([]int, requiredSize)
		copy(newDimensions, cb.dimensions) // Copy old values into the new slice
		newDimensions[index] = dimValue
		cb.dimensions = newDimensions
	}
	return cb
}

func (cb *CapacityBuilder) Build() *Capacity {
	capacityCopy := make([]int, len(cb.dimensions))
	copy(capacityCopy, cb.dimensions)
	return &Capacity{dimensions: capacityCopy}
}

type Capacity struct {
	dimensions []int
}

func NewCapacityFromBuilder(b *CapacityBuilder) *Capacity {
	return &Capacity{dimensions: b.dimensions}
}

// NewCapacity creates a new Capacity with the specified dimensions
func NewCapacity(dimensions []int) *Capacity {
	return &Capacity{dimensions: append([]int{}, dimensions...)}
}

func NewDefaultCapacity() *Capacity {
	return &Capacity{
		dimensions: make([]int, 1),
	}
}

// AddUp adds two capacities together
func AddUp(cap1, cap2 *Capacity) *Capacity {
	if cap1 == nil || cap2 == nil {
		panic("arguments must not be null")
	}
	maxLen := max(len(cap1.dimensions), len(cap2.dimensions))
	newDims := make([]int, maxLen)
	for i := 0; i < maxLen; i++ {
		newDims[i] = cap1.Get(i) + cap2.Get(i)
	}
	return NewCapacity(newDims)
}

// Subtract subtracts cap2 from cap1
func Subtract(cap1, cap2 *Capacity) *Capacity {
	if cap1 == nil || cap2 == nil {
		panic("arguments must not be null")
	}
	maxLen := max(len(cap1.dimensions), len(cap2.dimensions))
	newDims := make([]int, maxLen)
	for i := 0; i < maxLen; i++ {
		newDims[i] = cap1.Get(i) - cap2.Get(i)
	}
	return NewCapacity(newDims)
}

// Invert inverts the values of a Capacity
func Invert(cap *Capacity) *Capacity {
	if cap == nil {
		panic("arguments must not be null")
	}
	newDims := make([]int, len(cap.dimensions))
	for i, val := range cap.dimensions {
		newDims[i] = -val
	}
	return NewCapacity(newDims)
}

func (c *Capacity) AddDimension(index, dimValue int) {
	if index < len(c.dimensions) {
		c.dimensions[index] = dimValue
	} else {
		newDimensions := make([]int, index+1)
		copy(newDimensions, c.dimensions)
		newDimensions[index] = dimValue
		c.dimensions = newDimensions
	}
}

// Get returns the value of the given dimension
func (c *Capacity) Get(index int) int {
	if index < len(c.dimensions) {
		return c.dimensions[index]
	}
	return 0
}

func (c *Capacity) NuOfDimensions() int {
	return len(c.dimensions)
}

// String returns a string representation of the Capacity
func (c *Capacity) String() string {
	return fmt.Sprintf("%v", c.dimensions)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
