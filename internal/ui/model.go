package ui

import (
	"log"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/blueboy/internal/network"
)
type MenuState int


const (
	MainMenu MenuState = iota
	ScanMenu
	PairedMenu
)

type model struct {
	cursor        int
	CurrentMenu   MenuState
	MainOptions   []string
	ScanResults   []network.Device
	PairedDevices []network.Device
	status        string
	password      Password // new connection password
	width         int
	height        int
}

func InitialModel() model {

	return model{
		cursor:        0,
		CurrentMenu:   MainMenu,
		MainOptions:   []string{"Scan Devices", "Paired Connections"},
		ScanResults:   []network.Device{},
		PairedDevices: []network.Device{},
		password:      NewPasswordModel(),
	}
}

func (m model) Init() tea.Cmd {
	// the backend for this program is bluetoothctl
	_, err := exec.LookPath("nmcli")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	return nil
}
