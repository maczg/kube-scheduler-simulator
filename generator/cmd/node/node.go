package node

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sigs.k8s.io/kube-scheduler-simulator/generator/internal"
	"sigs.k8s.io/kube-scheduler-simulator/generator/pkg"
	"strconv"
	"text/template"
)

var Cmd = &cobra.Command{
	Use:   "node",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

		cpu, _ := cmd.Flags().GetString("cpu")
		memory, _ := cmd.Flags().GetString("memory")
		pod, _ := cmd.Flags().GetString("pod")
		number, _ := cmd.Flags().GetInt("number")
		dir, _ := cmd.Flags().GetString("dir")
		fileName, _ := cmd.Flags().GetString("filename")

		nodes := make([]pkg.Node, number)
		for i := 0; i < number; i++ {
			nodes[i] = pkg.Node{
				Name:           "node-" + strconv.Itoa(i),
				CPUCapacity:    cpu,
				MemoryCapacity: memory,
				PodCapacity:    pod,
			}
		}
		nodeListData := pkg.NodeList{
			Nodes: nodes,
		}

		nodeListTemplate, err := template.ParseFS(internal.EmbedFs, "template/node_list.tmpl")
		if err != nil {
			panic(fmt.Errorf("failed to parse node list template: %w", err))
		}

		outputFile := filepath.Join(dir, fileName+".yaml")
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}

		file, err := os.Create(outputFile)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		err = nodeListTemplate.Execute(file, nodeListData)
		if err != nil {
			panic(fmt.Errorf("failed to execute node list template: %w", err))
		}
	},
}

func init() {
	Cmd.Flags().StringP("filename", "f", "node_list", "Name of the node")
	Cmd.Flags().IntP("number", "n", 1, "Number of nodes to generate")
	Cmd.Flags().StringP("cpu", "c", "4", "CPU capacity of the node")
	Cmd.Flags().StringP("memory", "m", "16", "Memory capacity of the node")
	Cmd.Flags().StringP("pod", "p", "110", "Pod capacity of the node")
	Cmd.Flags().StringP("dir", "d", "output", "Output directory")
}
