package network

import (
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type NetworkProcess struct {
	PID          string
	Command      string
	User         string
	Type         string
	Protocol     string
	LocalPort    int
	LocalAddress string
	RemotePort   int
	RemoteAddress string
	State        string
}

func GetActivePorts(filterState string) ([]NetworkProcess, error) {
	cmd := exec.Command("lsof", "-i", "-P", "-n")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute lsof: %w", err)
	}

	return parseLsofOutput(string(output), filterState)
}

func GetNetworkStats() (map[string]int, error) {
	processes, err := GetActivePorts("")
	if err != nil {
		return nil, err
	}

	stats := make(map[string]int)
	stats["total"] = len(processes)
	stats["tcp"] = 0
	stats["udp"] = 0
	stats["listen"] = 0
	stats["established"] = 0
	stats["ipv4"] = 0
	stats["ipv6"] = 0

	for _, proc := range processes {
		if proc.Protocol == "TCP" {
			stats["tcp"]++
		} else if proc.Protocol == "UDP" {
			stats["udp"]++
		}

		if proc.Type == "IPv4" {
			stats["ipv4"]++
		} else if proc.Type == "IPv6" {
			stats["ipv6"]++
		}

		if proc.State == "LISTEN" {
			stats["listen"]++
		} else if proc.State == "ESTABLISHED" {
			stats["established"]++
		}
	}

	return stats, nil
}

func parseLsofOutput(output string, filterState string) ([]NetworkProcess, error) {
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

		proc := parseNetworkInfo(name, typeField)
		if proc.LocalPort == 0 {
			continue
		}

		if filterState != "" && proc.State != filterState {
			continue
		}

		key := fmt.Sprintf("%s:%s:%d:%s:%d", pid, proc.LocalAddress, proc.LocalPort, proc.RemoteAddress, proc.RemotePort)
		if seen[key] {
			continue
		}
		seen[key] = true

		process := NetworkProcess{
			PID:           pid,
			Command:       command,
			User:          user,
			Type:          typeField,
			Protocol:      proc.Protocol,
			LocalPort:     proc.LocalPort,
			LocalAddress:  proc.LocalAddress,
			RemotePort:    proc.RemotePort,
			RemoteAddress: proc.RemoteAddress,
			State:         proc.State,
		}

		processes = append(processes, process)
	}

	sort.Slice(processes, func(i, j int) bool {
		if processes[i].LocalPort != processes[j].LocalPort {
			return processes[i].LocalPort < processes[j].LocalPort
		}
		return processes[i].LocalAddress < processes[j].LocalAddress
	})

	return processes, nil
}

type connectionInfo struct {
	Protocol      string
	LocalPort     int
	LocalAddress  string
	RemotePort    int
	RemoteAddress string
	State         string
}

func parseNetworkInfo(name, typeField string) connectionInfo {
	parts := strings.Fields(name)
	if len(parts) < 2 {
		return connectionInfo{}
	}

	protocol := parts[0]
	if protocol != "TCP" && protocol != "UDP" {
		return connectionInfo{}
	}

	addrPort := parts[1]
	state := ""
	if len(parts) > 2 && strings.HasPrefix(parts[2], "(") {
		state = strings.Trim(parts[2], "()")
	}

	localAddr := ""
	localPort := 0
	remoteAddr := ""
	remotePort := 0

	if strings.Contains(addrPort, "->") {
		parts := strings.Split(addrPort, "->")
		if len(parts) == 2 {
			localAddr, localPort = parseAddressPort(parts[0])
			remoteAddr, remotePort = parseAddressPort(parts[1])
		}
	} else {
		localAddr, localPort = parseAddressPort(addrPort)
	}

	return connectionInfo{
		Protocol:      protocol,
		LocalPort:     localPort,
		LocalAddress:  localAddr,
		RemotePort:    remotePort,
		RemoteAddress: remoteAddr,
		State:         state,
	}
}

func parseAddressPort(addrPort string) (string, int) {
	colonIdx := strings.LastIndex(addrPort, ":")
	if colonIdx == -1 {
		return "", 0
	}

	address := addrPort[:colonIdx]
	portStr := addrPort[colonIdx+1:]

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0
	}

	if address == "*" {
		address = "0.0.0.0"
	}

	if strings.HasPrefix(address, "[") && strings.HasSuffix(address, "]") {
		address = strings.Trim(address, "[]")
	}

	return address, port
}

