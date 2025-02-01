package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"sigs.k8s.io/kube-scheduler-simulator/generator/cmd/node"
	"sigs.k8s.io/kube-scheduler-simulator/generator/cmd/pod"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "generator",
	Short: "Generator is a tool to generate yaml files for pod and node scheduler simulator",
	Long:  `Generator is a tool to generate yaml files for pod and node scheduler simulator.`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(pod.Cmd)
	rootCmd.AddCommand(node.Cmd)
}
