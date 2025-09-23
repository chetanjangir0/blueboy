package ui

type ConnectDeviceMsg struct {
	Output string
	Err    error
}

type StartScanMsg struct {
	Devices []Device
	Err     error
}

type FetchPairedMsg struct {
	Devices []Device
	Err     error
}
