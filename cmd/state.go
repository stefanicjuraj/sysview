package cmd

import (
	"fmt"
	"strings"
	"github.com/stefanicjuraj/sysview/internal/network"
	"github.com/spf13/cobra"
)

var stateCmd = &cobra.Command{
	Use:   "state [STATE]",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stateFilter := strings.ToUpper(args[0])
		processes, err := network.GetActivePorts(stateFilter)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
			return
		}

		if len(processes) == 0 {
			fmt.Printf("No active network processes with state: %s\n", stateFilter)
			return
		}

		displayProcesses(processes)
	},
}

func init() {
	rootCmd.AddCommand(stateCmd)
}

