package controller

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
