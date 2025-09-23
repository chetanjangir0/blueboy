package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/blueboy/internal/commands"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			if !m.password.isAsking {
				return m, tea.Quit
			}

		case "j", "down":
			if m.cursor < m.itemCount()-1 && !m.password.isAsking {
				m.cursor++
			}

		case "k", "up":
			if m.cursor > 0 && !m.password.isAsking {
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
					return m, commands.StartScan(m.nm)
				case "Paired Connections":
					m.CurrentMenu = PairedMenu
					m.cursor = 0
					return m, commands.FetchPaired(m.nm)
				}
				m.cursor = 0 // reset cursor pos

			case PairedMenu:
				m.status = "connecting..."
				return m, ConnectDevice(m.PairedDevices[m.cursor].UUID, m.nm)
			case ScanMenu:
				if m.cursor >= len(m.ScanResults) {
					return m, nil
				}
				selectedDevice := m.ScanResults[m.cursor]
				if selectedDevice.Security == "" {
					m.status = "connecting..."
					return m, pairNewDevice(selectedDevice, "")
				}
				if !m.password.isAsking {
					m.password.isAsking = true
					m.password.passwordInput.Reset()
					m.password.passwordInput.Focus()
					m.status = "This network requires a password:"
					return m, textinput.Blink
				}

				enteredPassword := m.password.passwordInput.Value()
				m.password.passwordInput.SetValue("") // Clear immediately
				m.password.passwordInput.Blur()       // Remove focus
				m.password.isAsking = false

				m.status = "connecting..."
				return m, pairNewDevice(selectedDevice, enteredPassword)

			}
		case "b", "esc":
			if !m.password.isAsking {
				m.CurrentMenu = MainMenu
				m.cursor = 0
				m.status = ""
				return m, nil
			}
		}
	case connectDeviceMsg:
		if msg.err != nil {
			m.status = msg.err.Error()
			return m, nil
		}
		m.status = msg.output
		return m, nil
	case startScanMsg:
		if msg.err != nil {
			m.status = msg.err.Error()
			return m, nil
		}
		m.ScanResults = msg.devices
		m.status = ""
		return m, nil
	case fetchPairedMsg:
		if msg.err != nil {
			return m, nil
		}
		m.PairedDevices = msg.devices
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	if m.password.isAsking {
		m.password.passwordInput, cmd = m.password.passwordInput.Update(msg)
	}
	return m, cmd
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
