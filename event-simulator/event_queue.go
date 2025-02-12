package main

// EventQueue implements heap.Interface for priority queue
type EventQueue []*Event

// Len returns the length of the event queue
func (eq *EventQueue) Len() int {
	return len(*eq)
}

// Less returns true if the event at index i should be ordered before the event at index j
func (eq *EventQueue) Less(i, j int) bool {
	return (*eq)[i].Time.Before((*eq)[j].Time)
}

// Swap swaps the events at index i and j
func (eq *EventQueue) Swap(i, j int) {
	(*eq)[i], (*eq)[j] = (*eq)[j], (*eq)[i]
}

// Push pushes an event onto the event queue
func (eq *EventQueue) Push(x interface{}) {
	*eq = append(*eq, x.(*Event))
}

// Pop pops an event from the event queue
func (eq *EventQueue) Pop() interface{} {
	old := *eq
	n := len(old)
	x := old[n-1]
	*eq = old[0 : n-1]
	return x
}
