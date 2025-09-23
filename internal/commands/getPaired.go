package commands

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/blueboy/internal/network"
	"github.com/chetanjangir0/blueboy/internal/messages"
)

func FetchPaired(nm network.NetworkManager) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		devices, err := nm.FetchPaired(ctx)
		if err != nil {
			return messages.FetchPairedMsg{Err: err}
		}

		return messages.FetchPairedMsg{Devices: devices}
	}
}
