package main

//todo
// wifi on off function
// use ticks to update frequently
// s for scan

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
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
	Name     string
	Type     string
	UUID     string
	Security string
}

type Password struct {
	state         string // normal asking processing
	passwordInput textinput.Model
	status        string
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
}

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

func initialModel() model {
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
		password:      Password{state: "normal", passwordInput: ti},
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
			if m.password.state != "asking" {
				return m, tea.Quit
			}

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
					m.status = "fetching devices..."
					return m, startScan()
				case "Paired Connections":
					m.CurrentMenu = PairedMenu
					m.cursor = 0
					return m, fetchPaired()
				}
				m.cursor = 0 // reset cursor pos

			case PairedMenu:
				log.Println("connecting")
				m.status = "connecting..."
				return m, connectDevice(m.PairedDevices[m.cursor].UUID)
			case ScanMenu:
				log.Println("connection pairing started")
				selectedDevice := m.ScanResults[m.cursor]
				if selectedDevice.Security == "" {
					m.status = "connecting..."
					return m, pairNewDevice(selectedDevice, "")
				}
				if m.password.state != "asking" {
					m.password.state = "asking"
					m.password.passwordInput.Reset()
					m.password.passwordInput.Focus()
					m.status = "This network requires a password:"
					return m, textinput.Blink
				}

				enteredPassword := m.password.passwordInput.Value()
				m.password.passwordInput.SetValue("") // Clear immediately
				m.password.passwordInput.Blur()       // Remove focus
				m.password.state = "normal"

				log.Printf("Processing password: \"%s\" (length %d)...", enteredPassword, len(enteredPassword))
				m.status = "connecting..."
				return m, pairNewDevice(selectedDevice, enteredPassword)

			}
		case "b", "esc":
			m.CurrentMenu = MainMenu
			m.cursor = 0
			m.status = ""
			return m, nil
		}
	case connectDeviceMsg:
		if msg.err != nil {
			log.Println("Error connecting:", msg.err)
			m.status = msg.err.Error()
			return m, nil
		}
		m.status = msg.output
		return m, nil
	case startScanMsg:
		if msg.err != nil {
			log.Println("Error scanning devices:", msg.err)
			m.status = msg.err.Error()
			return m, nil
		}
		m.ScanResults = msg.devices
		m.status = ""
		return m, nil
	case fetchPairedMsg:
		if msg.err != nil {
			log.Println("Error fetching paired devices:", msg.err)
			return m, nil
		}
		m.PairedDevices = msg.devices
		return m, nil

	}

	var cmd tea.Cmd
	if m.password.state == "asking" {
		m.password.passwordInput, cmd = m.password.passwordInput.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	// the header
	s := "Blueman\n\n"

	if m.password.state == "asking" {
		s += m.password.passwordInput.View()
	} else {
		switch m.CurrentMenu {

		case MainMenu:
			s += renderList(m.MainOptions, m.cursor)
		case ScanMenu:
			s += "Scan Results\n" + renderList(m.ScanResults, m.cursor)
		case PairedMenu:
			s += "Paired Devices\n" + renderList(m.PairedDevices, m.cursor)

		}
	}

	//status
	s += "\n" + m.status + "\n"

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
		output, err := exec.CommandContext(ctx, "nmcli", "-t", "-f", "SSID,SECURITY", "device", "wifi", "list").CombinedOutput()
		if ctx.Err() == context.DeadlineExceeded {
			return startScanMsg{err: fmt.Errorf("Error: connection timed out")}

		}
		if err != nil {
			return startScanMsg{err: err}
		}
		outputString := strings.Trim(string(output), "\n")
		log.Println(outputString)
		outputStringSlice := strings.Split(outputString, "\n")

		var devices = make([]Device, len(outputStringSlice))
		for i, d := range outputStringSlice {
			deviceInfo := strings.Split(d, ":")

			if len(deviceInfo) == 0 { // if there is no ssid/name
				continue
			} else if len(deviceInfo) == 1 { // if there is no security
				deviceInfo = append(deviceInfo, " ")
			}

			devices[i] = Device{
				Name:     deviceInfo[0],
				Type:     "wifi",
				Security: deviceInfo[1],
			}
		}
		return startScanMsg{devices: devices}
	}

}

func connectDevice(UUID string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		_, err := exec.CommandContext(ctx, "nmcli", "connection", "up", UUID).CombinedOutput()
		if ctx.Err() == context.DeadlineExceeded {
			return connectDeviceMsg{err: fmt.Errorf("Error: connection timed out")}
		}
		if err != nil {
			return connectDeviceMsg{err: err}
		}
		return connectDeviceMsg{output: "connection successfully activated"}
	}
}

func pairNewDevice(newDevice Device, password string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var err error
		if password != "" {
			_, err = exec.CommandContext(ctx, "nmcli", "device", "wifi", "connect", newDevice.Name, "password", password).CombinedOutput()
		} else {
			_, err = exec.CommandContext(ctx, "nmcli", "device", "wifi", "connect", newDevice.Name).CombinedOutput()
		}
		if ctx.Err() == context.DeadlineExceeded {
			return connectDeviceMsg{err: fmt.Errorf("Error: connection timed out")}
		}
		if err != nil {
			return connectDeviceMsg{err: err}
		}
		return connectDeviceMsg{output: "connection successfully activated"}
	}
}
