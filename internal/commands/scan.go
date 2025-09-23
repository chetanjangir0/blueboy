package commands

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/blueboy/internal/network"
	"github.com/chetanjangir0/blueboy/internal/messages"
)

func StartScan(nm network.NetworkManager) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		devices, err := nm.ScanDevices(ctx)
		if err != nil {
			return messages.StartScanMsg{Err: err}
		}
		return messages.StartScanMsg{Devices: devices}
	}

}
