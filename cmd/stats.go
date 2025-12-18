package cmd

import (
	"fmt"
	"os"
	"github.com/stefanicjuraj/sysview/internal/network"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show network statistics summary",
	Run: func(cmd *cobra.Command, args []string) {
		stats, err := network.GetNetworkStats()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		fmt.Printf("Total Connections: %d\n", stats["total"])
		fmt.Printf("TCP: %d\n", stats["tcp"])
		fmt.Printf("UDP: %d\n", stats["udp"])
		fmt.Printf("IPv4: %d\n", stats["ipv4"])
		fmt.Printf("IPv6: %d\n", stats["ipv6"])
		fmt.Printf("Listening: %d\n", stats["listen"])
		fmt.Printf("Established: %d\n", stats["established"])
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

