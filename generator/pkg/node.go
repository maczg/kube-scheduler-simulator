package pkg

type Node struct {
	Name           string
	CPUCapacity    string
	MemoryCapacity string
	PodCapacity    string
}

type NodeList struct {
	Nodes []Node
}
