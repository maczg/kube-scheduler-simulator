package main

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apiserver/pkg/storage/names"
	"math/rand"
)

func NewRandomPodSpec(name string) *corev1.Pod {
	p := &corev1.Pod{}
	minRequest := 100
	maxRequest := 1000
	randomCpuRequest := rand.Intn(maxRequest-minRequest) + minRequest
	randomMemoryRequest := rand.Intn(maxRequest-minRequest) + minRequest
	p.Name = names.SimpleNameGenerator.GenerateName(fmt.Sprintf("random-%s-", name))
	p.Namespace = "default"
	p.Spec.Containers = []corev1.Container{
		{
			Name:  "server",
			Image: "nginx:latest",
			Ports: []corev1.ContainerPort{
				{
					ContainerPort: 80,
				},
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceCPU:    resource.MustParse(fmt.Sprintf("%dm", randomCpuRequest)),
					corev1.ResourceMemory: resource.MustParse(fmt.Sprintf("%dMi", randomMemoryRequest)),
				},
			},
		},
	}
	return p
}
