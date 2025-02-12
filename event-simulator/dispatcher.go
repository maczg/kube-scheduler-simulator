package main

import (
	"container/heap"
	"context"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"time"
)

type Dispatcher struct {
	EventQueue  *EventQueue
	Notifier    chan Event
	CurrentTime time.Time
	ClientSet   *kubernetes.Clientset
}

type DispatcherOption func(*Dispatcher)

func WithClientSet(client *kubernetes.Clientset) DispatcherOption {
	return func(d *Dispatcher) {
		d.ClientSet = client
	}
}

func NewDispatcher(opts ...DispatcherOption) *Dispatcher {
	eq := &EventQueue{}
	initCurrentTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00Z")
	heap.Init(eq)
	d := &Dispatcher{
		EventQueue:  eq,
		Notifier:    make(chan Event),
		CurrentTime: initCurrentTime,
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

func (d *Dispatcher) DispatchPodWithTime(p *corev1.Pod, ArrivalTime int, DepartureTime int) {
	scheduleAt := d.CurrentTime.Add(time.Duration(ArrivalTime) * time.Second)
	unscheduleAt := d.CurrentTime.Add(time.Duration(DepartureTime) * time.Second)

	scheduleEvent := &Event{
		PodSpec: p,
		Type:    EventTypeSchedulePod,
		Time:    scheduleAt,
	}

	unscheduleEvent := &Event{
		PodSpec: p,
		Type:    EventTypeUnschedulePod,
		Time:    unscheduleAt,
	}

	heap.Push(d.EventQueue, scheduleEvent)
	heap.Push(d.EventQueue, unscheduleEvent)
}

func (d *Dispatcher) Run() {
	for d.EventQueue.Len() > 0 {
		event := heap.Pop(d.EventQueue).(*Event)
		d.CurrentTime = event.Time
		switch event.Type {
		case EventTypeSchedulePod:
			d.handleSchedule(event.PodSpec)
		case EventTypeUnschedulePod:
			d.handleUnschedule(event.PodSpec)
		}
		time.Sleep(10 * time.Second)
	}
}

func (d *Dispatcher) handleSchedule(p *corev1.Pod) {
	ctx := context.Background()
	logrus.Infof("[%s] scheduling pod %s", d.CurrentTime.Format("15:04:05"), p.Name)
	pod, err := d.ClientSet.CoreV1().Pods(p.Namespace).Create(ctx, p, metav1.CreateOptions{})
	if err != nil {
		logrus.Errorf("error scheduling pod %s: %v", p.Name, err)
		return
	}
	if err != nil {
		logrus.Errorf("error watching pod %s: %v", pod.Name, err)
		return
	}

	logrus.Infof("pod %s scheduled at %s", pod.Name, d.CurrentTime.Format("15:04:05"))
	return
}

func (d *Dispatcher) handleUnschedule(p *corev1.Pod) {
	ctx := context.Background()
	err := d.ClientSet.CoreV1().Pods(p.Namespace).Delete(ctx, p.Name, metav1.DeleteOptions{})
	if err != nil {
		logrus.Errorf("error unscheduling pod %s: %v", p.Name, err)
		return
	}
	logrus.Infof("pod %s unscheduled at %s", p.Name, d.CurrentTime.Format("15:04:05"))
	return
}

func (d *Dispatcher) WatchPod(ctx context.Context, p *corev1.Pod) {
	watch, err := d.ClientSet.CoreV1().Pods(p.Namespace).Watch(ctx, metav1.ListOptions{
		FieldSelector: "metadata.name=" + p.Name,
	})
	defer watch.Stop()
	if err != nil {
		logrus.Errorf("error watching pod %s: %v", p.Name, err)
		return
	}
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-watch.ResultChan():
			_pod, ok := e.Object.(*corev1.Pod)
			if !ok {
				logrus.Errorf("error converting to pod %s: %v", _pod.Name, err)
				return
			}
			logrus.Infof("pod %s is %s", _pod.Name, _pod.Status.Phase)
		}
	}
}
