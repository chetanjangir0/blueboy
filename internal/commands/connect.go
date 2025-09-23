package commands

import (
	"context"
	"time"

	"github.com/chetanjangir0/blueboy/internal/network"
	"github.com/chetanjangir0/blueboy/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func ConnectDevice(UUID string, nm network.NetworkManager) func() tea.Msg {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := nm.Connect(ctx, UUID)
		if err != nil {
			return ui.ConnectDeviceMsg{Err: err}
		}
		return ui.ConnectDeviceMsg{Output: "connection successfully activated"}
	}
}
