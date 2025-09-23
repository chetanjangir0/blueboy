package network

import "context"

import "github.com/chetanjangir0/blueboy/internal/ui"

type NetworkManager interface {
	FetchPaired(ctx context.Context) ([]ui.Device, error)
	ScanDevices(ctx context.Context) ([]ui.Device, error)
	Connect(ctx context.Context, UUID string) error
	Pair(ctx context.Context, device ui.Device, password string) error
}
