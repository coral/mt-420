package controller

import (
	"math"
)

type Selector struct {
	min   int
	max   int
	value int
	steps int
}

func NewSelector(initialValue int, steps int, min int, max int) Selector {
	return Selector{
		min:   min,
		max:   max,
		value: initialValue,
		steps: steps,
	}
}

func (s *Selector) Increment() bool {
	if s.value+s.steps <= s.max {
		s.value = s.value + s.steps
		return true
	}
	return false
}

func (s *Selector) Decrement() bool {
	if s.value-s.steps >= s.min {
		s.value = s.value - s.steps
		return true
	}
	return false
}

func (s *Selector) Value() int {
	return s.value
}

//Fuck Golang's lack of generics.

type FloatSelector struct {
	min   float64
	max   float64
	value float64
	steps float64
}

func NewFloatSelector(initialValue float64, steps float64, min float64, max float64) FloatSelector {
	return FloatSelector{
		min:   min,
		max:   max,
		value: initialValue,
		steps: steps,
	}
}

func (s *FloatSelector) Increment() bool {
	if s.value+s.steps <= s.max {
		s.value = math.Round((s.value+s.steps)*100) / 100
		return true
	}
	return false
}

func (s *FloatSelector) Decrement() bool {
	if s.value-s.steps >= s.min {
		s.value = math.Round((s.value-s.steps)*100) / 100
		return true
	}
	return false
}

func (s *FloatSelector) Value() float64 {
	return s.value
}
