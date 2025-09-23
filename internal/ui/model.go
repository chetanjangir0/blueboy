package ui

import (
	"log"
	"os/exec"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)
type MenuState int

type Device struct {
	Name     string
	Type     string
	UUID     string
	Security string
}

type Password struct {
	isAsking      bool
	passwordInput textinput.Model
}

const (
	MainMenu MenuState = iota
	ScanMenu
	PairedMenu
)

type model struct {
	cursor        int
	CurrentMenu   MenuState
	MainOptions   []string
	ScanResults   []Device
	PairedDevices []Device
	status        string
	password      Password // new connection password
	width         int
	height        int
}

func InitialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter password"
	ti.CharLimit = 20
	ti.Width = 30
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'

	return model{
		cursor:        0,
		CurrentMenu:   MainMenu,
		MainOptions:   []string{"Scan Devices", "Paired Connections"},
		ScanResults:   []Device{},
		PairedDevices: []Device{},
		password:      Password{isAsking: false, passwordInput: ti},
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
