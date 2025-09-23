package ui

type connectDeviceMsg struct {
	output string
	err    error
}

type startScanMsg struct {
	devices []Device
	err     error
}

type fetchPairedMsg struct {
	devices []Device
	err     error
}
