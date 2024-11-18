package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/zulubit/mole/pkg/consts"
	"github.com/shirou/gopsutil/net"
)

type ports []int

type reservedPorts struct {
	Ports ports `json:"ports"`
}

var startingPort = 9000

// FindAndReserveMolePorts identifies and reserves the next three available ports starting from `startingPort`.
// The function also updates the reserved ports file with the new reserved ports.
func FindAndReserveMolePorts() (ports, error) {
	reservedAndUsedPorts, err := getReservedAndUsedPorts()
	if err != nil {
		return ports{}, fmt.Errorf("failed to retrieve reserved and used ports: %w", err)
	}

	newPorts := ports{}
	for i := startingPort; len(newPorts) < 3; i++ {
		isReserved := false
		for _, port := range reservedAndUsedPorts {
			if port == i {
				isReserved = true
				break
			}
		}

		if !isReserved {
			reservedAndUsedPorts = append(reservedAndUsedPorts, i)
			newPorts = append(newPorts, i)
		}
	}

	if err := saveReservedPorts(reservedAndUsedPorts); err != nil {
		return ports{}, fmt.Errorf("failed to save reserved ports: %w", err)
	}

	return newPorts, nil
}

// saveReservedPorts writes a sorted list of reserved ports to a JSON file.
func saveReservedPorts(portsToSave ports) error {
	sort.Ints(portsToSave)

	data := reservedPorts{
		Ports: portsToSave,
	}

	fileData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal reserved ports data: %w", err)
	}
	if err := os.MkdirAll(consts.GetBasePath(), 0755); err != nil {
		return fmt.Errorf("failed to create base directory: %w", err)
	}

	if err := os.WriteFile(path.Join(consts.GetBasePath(), "reservedPorts.json"), fileData, 0644); err != nil {
		return fmt.Errorf("failed to write reserved ports to file: %w", err)
	}

	return nil
}

// getReservedAndUsedPorts retrieves a list of unique, currently used, and reserved ports.
func getReservedAndUsedPorts() (ports, error) {
	connections, err := net.Connections("tcp")
	if err != nil {
		return ports{}, fmt.Errorf("failed to retrieve TCP connections: %w", err)
	}

	usedPorts := ports{}
	for _, conn := range connections {
		usedPorts = append(usedPorts, int(conn.Laddr.Port))
	}

	reservedData, err := os.ReadFile(path.Join(consts.GetBasePath(), "reservedPorts.json"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return ports{}, fmt.Errorf("failed to read reserved ports file: %w", err)
	}

	reservedPorts := reservedPorts{}
	if len(reservedData) > 0 {
		if err := json.Unmarshal(reservedData, &reservedPorts); err != nil {
			return ports{}, fmt.Errorf("failed to unmarshal reserved ports data: %w", err)
		}
	}

	portSet := map[int]bool{}
	for _, port := range usedPorts {
		portSet[port] = true
	}
	for _, port := range reservedPorts.Ports {
		portSet[port] = true
	}

	uniquePorts := ports{}
	for port := range portSet {
		uniquePorts = append(uniquePorts, port)
	}

	return uniquePorts, nil
}

// PortReport generates a comma-separated report of all active TCP ports, sorted in ascending order.
func PortReport() (string, error) {
	connections, err := net.Connections("tcp")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve TCP connections: %w", err)
	}

	uniquePorts := map[string]bool{}
	for _, conn := range connections {
		uniquePorts[fmt.Sprintf("%d", conn.Laddr.Port)] = true
	}

	sortedPorts := []string{}
	for port := range uniquePorts {
		sortedPorts = append(sortedPorts, port)
	}

	sort.Strings(sortedPorts)
	return strings.Join(sortedPorts, ", "), nil
}
