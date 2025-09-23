
package messages 

import "github.com/chetanjangir0/blueboy/internal/network"

type ConnectDeviceMsg struct {
	Output string
	Err    error
}

type StartScanMsg struct {
	Devices []network.Device
	Err     error
}

type FetchPairedMsg struct {
	Devices []network.Device
	Err     error
}
