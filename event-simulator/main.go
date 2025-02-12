package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
	flag.Parse()
	if *kubeconfig == "" {
		logrus.Fatalf("Kubeconfig is required")
	}

	c := Client{}
	go func() {
		err := c.Connect("http://localhost:1212/api/v1/listwatchresources")
		if err != nil {
			logrus.Errorf("Error connecting to server: %s", err.Error())
		}
	}()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		logrus.Fatalf("Error building kubeconfig: %s", err.Error())
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		logrus.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}
	d := NewDispatcher(WithClientSet(client))
	p := NewRandomPodSpec("first")
	p.Spec.Containers[0].Resources.Requests["cpu"] = resource.MustParse(fmt.Sprintf("%dm", 10000))

	d.DispatchPodWithTime(p, 0, 10)
	p = NewRandomPodSpec("second")
	d.DispatchPodWithTime(p, 5, 15)
	p = NewRandomPodSpec("third")
	d.DispatchPodWithTime(p, 12, 22)
	p = NewRandomPodSpec("fourth")
	d.DispatchPodWithTime(p, 18, 28)
	d.Run()
}
