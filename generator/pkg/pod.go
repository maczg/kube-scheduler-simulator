package pkg

type Pod struct {
	Name          string
	Namespace     string
	Image         string
	ContainerName string
	CPULimit      string
	MemoryLimit   string
	CPURequest    string
	MemoryRequest string
	SchedulerName string
}
type PodList struct {
	Pods []Pod
}
