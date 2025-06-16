package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// debug file
	f, err := os.Create("debug.log")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	log.SetOutput(f)
	// log.SetFlags(log.LstdFlags | log.Lshortfile) // optional: adds timestamps + file:line

	program := tea.NewProgram(initialModel())
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}

type MenuState int

type Device struct {
	MacAddress string
	DeviceName string
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
}

func initialModel() model {
	return model{
		cursor:        0,
		CurrentMenu:   MainMenu,
		MainOptions:   []string{"Scan Devices", "Paired Connections"},
		ScanResults:   []Device{},
		PairedDevices: []Device{},
	}
}

func (m model) Init() tea.Cmd {
	// the backend for this program is bluetoothctl
	_, err := exec.LookPath("bluetoothctl")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	log.Println("bluetoothctl is available")
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if val, ok := msg.(tea.KeyMsg); ok {
		key := val.String()

		switch key {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "j", "down":
			if m.cursor < m.itemCount()-1 {
				m.cursor++
			}

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "enter":
			if m.CurrentMenu == MainMenu {

				switch m.MainOptions[m.cursor] {
				case "Scan Devices":
					m.CurrentMenu = ScanMenu
				case "Paired Connections":
					m.PairedDevices = getPairedDevices()
					m.CurrentMenu = PairedMenu
				}
				m.cursor = 0 // reset cursor pos

			}
		case "b", "esc":
			// 	out, err := exec.Command("pd").CombinedOutput()
			// 	if err != nil {
			// 		log.Print("Error: ", err)
			//
			// 	}
			// 	items := strings.Split(string(out), "\n")
			// 	for _, item := range items {
			// 		log.Println(item)
			// 	}
			m.CurrentMenu = MainMenu
		}

	}
	return m, nil
}

func (m model) View() string {
	// the header
	s := "Blueman\n\n"

	switch m.CurrentMenu {

	case MainMenu:
		s += renderList(m.MainOptions, m.cursor)
	case ScanMenu:
		s += "Scan Results\n" + renderList(m.ScanResults, m.cursor)
	case PairedMenu:
		s += "Paired Devices\n" + renderList(m.PairedDevices, m.cursor)

	}

	// footer
	s += "\nq: Quit, b/esc: Main menu.\n"

	return s
}

func (m model) itemCount() int {
	switch m.CurrentMenu {
	case MainMenu:
		return len(m.MainOptions)
	case ScanMenu:
		return len(m.ScanResults)
	case PairedMenu:
		return len(m.PairedDevices)
	default:
		return 0
	}
}

func renderList(list []string, cursor int) string {
	s := ""
	for i, item := range list {

		cursorView := ""
		if i == cursor {
			cursorView = "<"
		}

		s += fmt.Sprintf("%s %s\n", item, cursorView)
	}
	return s
}

func getPairedDevices() []Device{
	log.Println("fetching paired devices")
	output, err := exec.Command("bluetoothctl", "paired-devices").CombinedOutput()
	if err != nil {
		log.Println("Error: ", err)
		return nil
	}
	outputString := strings.Trim(string(output), "\n")
	log.Println(outputString)
	outputStringSlice := strings.Split(outputString, "\n")

	var devices = make([]Device, len(outputStringSlice))
	for i, d := range outputStringSlice {
		deviceInfo := strings.SplitN(d, " ", 3) // "Device <Mac address> <device name>"
		devices[i] = Device{
			MacAddress: deviceInfo[1],
			DeviceName: deviceInfo[2],
		}
	}	
	return devices

}
