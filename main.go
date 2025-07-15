package main

//todo
// wifi on off function
// add current status(loading, success, error) in the model and show in view
// use ticks to update frequently
// s for scan

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

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
	UUID string
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

type connectDeviceMsg struct {
	status string
	output string
}

type startScanMsg struct {
	devices []Device
	err     error
}

type fetchPairedMsg struct {
	devices []Device
	err     error
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

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
			switch m.CurrentMenu {
			case MainMenu:

				switch m.MainOptions[m.cursor] {
				case "Scan Devices":
					m.CurrentMenu = ScanMenu
					m.cursor = 0
					return m, startScan()
				case "Paired Connections":
					m.CurrentMenu = PairedMenu
					m.cursor = 0
					return m, fetchPaired()
				}
				m.cursor = 0 // reset cursor pos

			case PairedMenu:
				log.Println("connecting")
				return m, connectDevice(m.PairedDevices[m.cursor].UUID)

			}
		case "b", "esc":
			m.CurrentMenu = MainMenu
		}
	case connectDeviceMsg:
		log.Printf("%s:%s", msg.status, msg.output)
		return m, nil
	case startScanMsg:
		if msg.err != nil {
			log.Println("Error scanning devices:", msg.err)
			return m, nil
		}
		m.ScanResults = msg.devices
		return m, nil
	case fetchPairedMsg:
		if msg.err != nil {
			log.Println("Error fetching paired devices:", msg.err)
			return m, nil
		}
		m.PairedDevices= msg.devices
		return m, nil
		
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
			s += fmt.Sprintf("%s %s\n", v.Name, cursorView)

		}

	}
	return s
}


func fetchPaired() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		output, err := exec.CommandContext(ctx, "nmcli", "-t", "-f", "NAME,TYPE,UUID", "connection", "show").CombinedOutput()
		if ctx.Err() == context.DeadlineExceeded {
			return fetchPairedMsg{err: fmt.Errorf("timeout")}
		}
		if err != nil {
			return fetchPairedMsg{err: err}
		}
		outputString := strings.Trim(string(output), "\n")
		log.Println(outputString)
		outputStringSlice := strings.Split(outputString, "\n")

		var devices = make([]Device, len(outputStringSlice))
		for i, d := range outputStringSlice {
			deviceInfo := strings.Split(d, ":") // "Device <name> <type>"
			if len(deviceInfo) != 3 {
				return fetchPairedMsg{err: fmt.Errorf("unexpected number of fields")} 
			}
			devices[i] = Device{
				Name: deviceInfo[0],
				Type: deviceInfo[1],
				UUID: deviceInfo[2],
			}
		}
		return fetchPairedMsg{devices: devices} 
	}
}

func startScan() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		output, err := exec.CommandContext(ctx, "nmcli", "-t", "-f", "SSID", "device", "wifi", "list").CombinedOutput()
		if ctx.Err() == context.DeadlineExceeded {
			return startScanMsg{err: fmt.Errorf("timeout")}

		}
		if err != nil {
			return startScanMsg{err: err}
		}
		outputString := strings.Trim(string(output), "\n")
		log.Println(outputString)
		outputStringSlice := strings.Split(outputString, "\n")

		var devices = make([]Device, len(outputStringSlice))
		for i, d := range outputStringSlice {
			devices[i] = Device{
				Name: d,
				Type: "wifi",
				UUID: "",
			}
		}
		return startScanMsg{devices: devices}
	}

}

func connectDevice(UUID string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		output, err := exec.CommandContext(ctx, "nmcli", "connection", "up", UUID).CombinedOutput()
		if ctx.Err() == context.DeadlineExceeded {
			return connectDeviceMsg{status: "error", output: "connection timed out"}
		}
		if err != nil {
			return connectDeviceMsg{status: "error", output: string(output)}
		}
		return connectDeviceMsg{status: "success", output: string(output)}
	}
}
