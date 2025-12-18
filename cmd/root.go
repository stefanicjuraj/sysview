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
	Run: func(cmd *cobra.Command, args []string) {
		processes, err := network.GetActivePorts("")
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
			return
		}

		if len(processes) == 0 {
			fmt.Println("No active network processes.")
			return
		}

		displayProcesses(processes)
	},
}

func displayProcesses(processes []network.NetworkProcess) {
	fmt.Printf("Active Network Processes (%d)\n", len(processes))
	fmt.Println(strings.Repeat("-", 100))
	fmt.Printf("%-8s %-20s %-8s %-8s %-20s %-20s %-12s\n", "PID", "COMMAND", "PROTO", "PORT", "LOCAL ADDRESS", "REMOTE ADDRESS", "STATE")
	fmt.Println(strings.Repeat("-", 100))

	for _, proc := range processes {
		state := proc.State
		if state == "" {
			state = "-"
		}

		remoteAddr := proc.RemoteAddress
		if remoteAddr == "" {
			remoteAddr = "-"
		} else if proc.RemotePort > 0 {
			remoteAddr = fmt.Sprintf("%s:%d", remoteAddr, proc.RemotePort)
		}

		fmt.Printf("%-8s %-20s %-8s %-8d %-20s %-20s %-12s\n",
			proc.PID,
			truncateString(proc.Command, 20),
			proc.Protocol,
			proc.LocalPort,
			truncateString(proc.LocalAddress, 20),
			truncateString(remoteAddr, 20),
			state,
		)
	}
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

