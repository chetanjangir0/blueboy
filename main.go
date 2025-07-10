package main

//todo
// wifi on off function
// use uuid instead of name to connect

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
	// optional: adds timestamps + file:line
	// log.SetFlags(log.LstdFlags | log.Lshortfile)

	program := tea.NewProgram(initialModel())
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}

type MenuState int

type Device struct {
	Name string
	Type string
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

type nmcliMsg struct {
	status string
	output string
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
	_, err := exec.LookPath("nmcli")
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
				handleMainMenu(&m)

			} else if m.CurrentMenu == PairedMenu {
				connectDevice(m.PairedDevices[m.cursor].Name)

			}
		case "b", "esc":
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

func renderList[T string | Device](list []T, cursor int) string {
	s := ""
	for i, item := range list {

		cursorView := ""
		if i == cursor {
			cursorView = "<"
		}

		switch v := any(item).(type) {
		case string:
			s += fmt.Sprintf("%s %s\n", v, cursorView)
		case Device:
			s += fmt.Sprintf("%s %s %s\n", v.Name, v.Type, cursorView)

		}

	}
	return s
}

func getPairedDevices() []Device {
	log.Println("fetching remembered devices")
	output, err := exec.Command("nmcli", "-t", "-f", "NAME,TYPE", "connection", "show").CombinedOutput()
	if err != nil {
		log.Println("Error: ", err)
		return nil
	}
	outputString := strings.Trim(string(output), "\n")
	log.Println(outputString)
	outputStringSlice := strings.Split(outputString, "\n")

	var devices = make([]Device, len(outputStringSlice))
	for i, d := range outputStringSlice {
		deviceInfo := strings.Split(d, ":") // "Device <name> <type>"
		if len(deviceInfo) != 2 {
			log.Println("unexpected number of fields")
			return nil
		}
		devices[i] = Device{
			Name: deviceInfo[0],
			Type: deviceInfo[1],
		}
	}
	return devices

}

func getScanResults() []Device {
	exec.Command("bluetoothctl", "power", "on")
	exec.Command("bluetoothctl", "power", "on")
	exec.Command("bluetoothctl", "agent", "on")
	exec.Command("bluetoothctl", "default-agent")
	_, err := exec.Command("bluetoothctl", "scan", "on").CombinedOutput()
	if err != nil {
		log.Println("Error: scan err", err)
		return nil
	}
	output2, err := exec.Command("bluetoothctl", "devices").CombinedOutput()
	if err != nil {
		log.Println("Error: devices err", err)
		return nil
	}
	outputString := strings.Trim(string(output2), "\n")
	log.Println(outputString)
	outputStringSlice := strings.Split(outputString, "\n")

	var devices = make([]Device, len(outputStringSlice))
	for i, d := range outputStringSlice {
		deviceInfo := strings.SplitN(d, " ", 3) // "Device <Mac address> <device name>"
		devices[i] = Device{
			Name: deviceInfo[1],
			Type: deviceInfo[2],
		}
	}
	return devices

}

func handleMainMenu(m *model) {
	switch m.MainOptions[m.cursor] {
	case "Scan Devices":
		m.ScanResults = getScanResults()
		m.CurrentMenu = ScanMenu
	case "Paired Connections":
		m.PairedDevices = getPairedDevices()
		m.CurrentMenu = PairedMenu
	}
	m.cursor = 0 // reset cursor pos
}

func connectDevice(deviceName string) {
	log.Printf("attempting to connect to %v\n", deviceName)
	_, err := exec.Command("nmcli", "connection", "up", deviceName).CombinedOutput()
	if err != nil {
		log.Println("Error while connecting:", err)
		return
	}
	log.Println("connected succefully")
}
