package pod

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
	Use:   "pod",
	Short: "Generate a pod list yaml file",
	Long:  "Generate a pod list yaml file with the specified number of pods and pod configurations",
	Run: func(cmd *cobra.Command, args []string) {
		cpuLimit, _ := cmd.Flags().GetString("cpu-limit")
		memoryLimit, _ := cmd.Flags().GetString("memory-limit")
		cpuRequest, _ := cmd.Flags().GetString("cpu-request")
		memoryRequest, _ := cmd.Flags().GetString("memory-request")
		schedulerName, _ := cmd.Flags().GetString("scheduler-name")
		containerName, _ := cmd.Flags().GetString("container-name")
		image, _ := cmd.Flags().GetString("image")
		namespace, _ := cmd.Flags().GetString("namespace")

		number, _ := cmd.Flags().GetInt("number")
		dir, _ := cmd.Flags().GetString("dir")
		fileName, _ := cmd.Flags().GetString("filename")

		pods := make([]pkg.Pod, number)
		for i := 0; i < number; i++ {
			pods[i] = pkg.Pod{
				Name:          "pod-" + strconv.Itoa(i),
				CPULimit:      cpuLimit,
				MemoryLimit:   memoryLimit,
				CPURequest:    cpuRequest,
				MemoryRequest: memoryRequest,
				SchedulerName: schedulerName,
				ContainerName: containerName,
				Image:         image,
				Namespace:     namespace,
			}
		}
		podListdata := pkg.PodList{
			Pods: pods,
		}

		podTemplate, err := template.ParseFS(internal.EmbedFs, "template/pod_list.tmpl")
		if err != nil {
			panic(fmt.Errorf("failed to parse pod list template: %w", err))
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

		err = podTemplate.Execute(file, podListdata)
		if err != nil {
			panic(fmt.Errorf("failed to execute node list template: %w", err))
		}
	},
}

func init() {
	Cmd.Flags().StringP("filename", "f", "pod_list", "Filename of the pod list")
	Cmd.Flags().IntP("number", "n", 1, "Number of pods to generate")
	Cmd.Flags().StringP("cpu-limit", "c", "2", "CPU limit of the pod")
	Cmd.Flags().StringP("memory-limit", "m", "4", "Memory limit of the pod")
	Cmd.Flags().StringP("cpu-request", "r", "2", "CPU request of the pod")
	Cmd.Flags().StringP("memory-request", "o", "4", "Memory request of the pod")
	Cmd.Flags().StringP("scheduler-name", "s", "default-scheduler", "Scheduler name of the pod")
	Cmd.Flags().StringP("container-name", "t", "container", "Container name of the pod")
	Cmd.Flags().StringP("image", "i", "registry.k8s.io/pause:3.5", "Image of the pod")
	Cmd.Flags().StringP("namespace", "z", "default", "Namespace of the pod")
	Cmd.Flags().StringP("dir", "d", "output", "Output directory")
}
