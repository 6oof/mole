package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/6oof/mole/pkg/consts"
	"github.com/shirou/gopsutil/net"
)

type ports []int

type reservedPorts struct {
	Ports ports `json:"ports"`
}

var reservedPortsFile = path.Join(consts.BasePath, "reservedPorts.json")
var startAt = 8000

func FindAndReserveMolePorts() (ports, error) {
	ru, err := reservedAndUsed()
	if err != nil {
		return ports{}, err
	}

	newPorts := ports{}

	for i := startAt; len(newPorts) < 3; i++ {
		found := false
		for _, p := range ru {
			if p == i {
				found = true
			}
		}

		if !found {
			ru = append(ru, i)
			newPorts = append(newPorts, i)
		}
	}

	err = writeReservedPorts(ru)
	if err != nil {
		return ports{}, err
	}

	return newPorts, nil
}

func writeReservedPorts(portsToWrite ports) error {
	sort.Ints(portsToWrite)

	wp := reservedPorts{
		Ports: portsToWrite,
	}

	pj, err := json.Marshal(wp)
	if err != nil {
		return err
	}

	os.MkdirAll(consts.BasePath, 0755)

	err = os.WriteFile(reservedPortsFile, pj, 0644)
	if err != nil {
		return err
	}

	return nil
}

func reservedAndUsed() (ports, error) {
	con, err := net.Connections("tcp")
	if err != nil {
		return ports{}, err
	}

	usedPorts := ports{}

	for _, conn := range con {
		usedPorts = append(usedPorts, int(conn.Laddr.Port))
	}

	ps, err := os.ReadFile(reservedPortsFile)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return ports{}, err
		}
	}

	reservedPorts := reservedPorts{}

	if len(ps) > 0 {
		err = json.Unmarshal(ps, &reservedPorts)
		if err != nil {
			return ports{}, err
		}
	}

	portMap := map[int]bool{}
	for _, p := range usedPorts {
		portMap[p] = true
	}
	for _, p := range reservedPorts.Ports {
		portMap[p] = true
	}

	uniquePorts := ports{}
	for p := range portMap {
		uniquePorts = append(uniquePorts, p)
	}

	return uniquePorts, nil
}

func PortReport() (string, error) {
	con, err := net.Connections("tcp")
	if err != nil {
		return "", err
	}

	portUni := map[string]bool{}

	for _, conn := range con {
		portUni[fmt.Sprintf("%d", conn.Laddr.Port)] = true
	}

	usedPorts := []string{}
	for ps := range portUni {
		usedPorts = append(usedPorts, ps)
	}

	sort.Strings(usedPorts)

	return strings.Join(usedPorts, ", "), nil
}
