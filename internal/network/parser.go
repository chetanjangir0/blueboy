package network

import (
	"fmt"
	"strings"
)

func parsePaired(output string) ([]Device, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var devices []Device
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) != 3 {
			return nil, fmt.Errorf("unexpected number of fields in line: %s", line)
		}
		devices = append(devices, Device{
			Name: fields[0], Type: fields[1], UUID: fields[2],
		})
	}
	return devices, nil
}

func parseScan(output string) []Device {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var devices []Device
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) == 0 || fields[0] == "" {
			continue
		}
		security := " "
		if len(fields) > 1 {
			security = fields[1]
		}
		devices = append(devices, Device{Name: fields[0], Type: "wifi", Security: security})
	}
	return devices
}
