package commands

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/blueboy/internal/network"
	"github.com/chetanjangir0/blueboy/internal/messages"
)

func pairNewDevice(newDevice network.Device, password string, nm network.NetworkManager) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := nm.Pair(ctx,newDevice, password)
		if err != nil {
			return messages.ConnectDeviceMsg{Err: err}
		}
		return messages.ConnectDeviceMsg{Output: "connection successfully activated"}
	}
}
