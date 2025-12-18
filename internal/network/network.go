package network

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type NetworkProcess struct {
	PID     string
	Command string
	User    string
	Type    string
	Port    int
	Address string
	State   string
}

func GetActivePorts() ([]NetworkProcess, error) {
	cmd := exec.Command("lsof", "-i", "-P", "-n")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute lsof: %w", err)
	}

	return parseLsofOutput(string(output))
}

func parseLsofOutput(output string) ([]NetworkProcess, error) {
	var processes []NetworkProcess
	lines := strings.Split(output, "\n")

	if len(lines) < 2 {
		return processes, nil
	}

	seen := make(map[string]bool)

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 8 {
			continue
		}

		command := fields[0]
		pid := fields[1]
		user := fields[2]
		typeField := fields[4]

		name := strings.Join(fields[7:], " ")

		if !strings.HasPrefix(typeField, "IPv") {
			continue
		}

		port, address, state := parseNetworkInfo(name)

		if port == 0 {
			continue
		}

		if state != "LISTEN" {
			continue
		}

		key := fmt.Sprintf("%s:%s:%d", pid, address, port)
		if seen[key] {
			continue
		}
		seen[key] = true

		process := NetworkProcess{
			PID:     pid,
			Command: command,
			User:    user,
			Type:    typeField,
			Port:    port,
			Address: address,
			State:   state,
		}

		processes = append(processes, process)
	}

	sort.Slice(processes, func(i, j int) bool {
		if processes[i].Port != processes[j].Port {
			return processes[i].Port < processes[j].Port
		}
		return processes[i].Address < processes[j].Address
	})

	return processes, nil
}

func parseNetworkInfo(name string) (int, string, string) {
	parts := strings.Fields(name)
	if len(parts) < 2 {
		return 0, "", ""
	}

	protocol := parts[0]
	if protocol != "TCP" && protocol != "UDP" {
		return 0, "", ""
	}

	addrPort := parts[1]
	state := ""
	if len(parts) > 2 && strings.HasPrefix(parts[2], "(") {
		state = strings.Trim(parts[2], "()")
	}

	colonIdx := strings.LastIndex(addrPort, ":")
	if colonIdx == -1 {
		return 0, "", ""
	}

	address := addrPort[:colonIdx]
	portStr := addrPort[colonIdx+1:]

	if strings.Contains(portStr, "->") {
		portStr = strings.Split(portStr, "->")[0]
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0, "", ""
	}

	if address == "*" {
		address = "0.0.0.0"
	}

	if strings.HasPrefix(address, "[") && strings.HasSuffix(address, "]") {
		address = strings.Trim(address, "[]")
	}

	return port, address, state
}

