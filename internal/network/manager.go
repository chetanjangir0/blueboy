package network

import "context"

type Device struct {
	Name     string
	Type     string
	UUID     string
	Security string
}

type NetworkManager interface {
	FetchPaired(ctx context.Context) ([]Device, error)
	ScanDevices(ctx context.Context) ([]Device, error)
	Connect(ctx context.Context, UUID string) error
	Pair(ctx context.Context, device Device, password string) error
}
