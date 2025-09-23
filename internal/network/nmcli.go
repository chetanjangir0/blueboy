package network

import (
	"context"
	"fmt"
	"os/exec"
)

type NmcliManager struct{}

func (NmcliManager) FetchPaired(ctx context.Context) ([]Device, error) {
	output, err := exec.CommandContext(ctx, "nmcli", "-t", "-f", "NAME,TYPE,UUID", "connection", "show").CombinedOutput()
	if err != nil {
		return nil, err
	}
	return parsePaired(string(output))
}

func (NmcliManager) ScanDevices(ctx context.Context) ([]Device, error) {
	output, err := exec.CommandContext(ctx, "nmcli", "-t", "-f", "SSID,SECURITY", "device", "wifi", "list").CombinedOutput()
	if err != nil {
		return nil, err
	}
	return parseScan(string(output)), nil
}

func (NmcliManager) Connect(ctx context.Context, UUID string) error {
	_, err := exec.CommandContext(ctx, "nmcli", "connection", "up", UUID).CombinedOutput()
	if err != nil {
		return fmt.Errorf("connect failed: %w", err)
	}
	return nil
}

func (NmcliManager) Pair(ctx context.Context, device Device, password string) error {
	var cmd *exec.Cmd
	if password != "" {
		cmd = exec.CommandContext(ctx, "nmcli", "device", "wifi", "connect", device.Name, "password", password)
	} else {
		cmd = exec.CommandContext(ctx, "nmcli", "device", "wifi", "connect", device.Name)
	}
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pair failed: %w", err)
	}
	return nil
}
