package main

import (
	corev1 "k8s.io/api/core/v1"
	"time"
)

// EventType is the type of event
type EventType string

const (
	// EventTypeSchedulePod is the event type for scheduling a pod
	EventTypeSchedulePod EventType = "SchedulePod"
	// EventTypeUnschedulePod is the event type for unscheduling a pod
	EventTypeUnschedulePod EventType = "UnschedulePod"
)

// Event is the event that is sent to the dispatcher
type Event struct {
	// PodSpec is the pod that the event is related to
	PodSpec *corev1.Pod
	// Type is the type of event
	Type EventType
	// Time is the time at which the event occurs
	Time time.Time
}
