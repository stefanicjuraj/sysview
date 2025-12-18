package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/stefanicjuraj/sysview/internal/network"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sysview",
	Short: "List active network processes and ports",
	Run: func(cmd *cobra.Command, args []string) {
		processes, err := network.GetActivePorts()
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
			return
		}

		if len(processes) == 0 {
			fmt.Println("No active network processes found.")
			return
		}

		fmt.Printf("Active Network Processes (%d)\n", len(processes))
		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("%-8s %-20s %-12s %-8s %-18s %-12s\n", "PID", "COMMAND", "USER", "PORT", "ADDRESS", "STATE")
		fmt.Println(strings.Repeat("-", 80))

		for _, proc := range processes {
			state := proc.State
			if state == "" {
				state = "-"
			}
			fmt.Printf("%-8s %-20s %-12s %-8d %-18s %-12s\n",
				proc.PID,
				truncateString(proc.Command, 20),
				truncateString(proc.User, 12),
				proc.Port,
				proc.Address,
				state,
			)
		}
	},
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

