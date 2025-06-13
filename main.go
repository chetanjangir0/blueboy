package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

func main() {
	program := tea.NewProgram(initialModel())
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}

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
	ScanResults   []string
	PairedDevices []string
}

func initialModel() model {
	return model{
		cursor:        0,
		CurrentMenu:   MainMenu,
		MainOptions:   []string{"Scan Devices", "Paired Connections", "Quit"},
		ScanResults:   []string{"testres1", "testres2", "testres3"},
		PairedDevices: []string{"paired1", "Paired2", "paired3"},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if val, ok := msg.(tea.KeyMsg); ok {
		key := val.String()

		switch key {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "j", "down":
			var maxOptions int

			switch m.CurrentMenu {

			case MainMenu:
				maxOptions = len(m.MainOptions)
			case ScanMenu:
				maxOptions = len(m.ScanResults)
			case PairedMenu:
				maxOptions = len(m.PairedDevices)
			}
			if m.cursor < maxOptions-1 {
				m.cursor++
			}

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "enter":
			if m.CurrentMenu == MainMenu {

				switch m.MainOptions[m.cursor]{
			
				case "Scan Devices":
					m.CurrentMenu = ScanMenu
				case "Paired Connections":
					m.CurrentMenu = PairedMenu
				}
			}
			// _, ok := m.done[m.cursor]
			// if ok {
			// 	delete(m.done, m.cursor)
			// } else {
			// 	m.done[m.cursor] = struct{}{}
			// }
		}

	}
	return m, nil
}

func (m model) View() string {
	// the header
	s := "Blueman\n\n"

	switch m.CurrentMenu {

	case MainMenu:
		for i, option := range m.MainOptions {
			cursor := ""
			if m.cursor == i {
				cursor = "<"
			}

			s += fmt.Sprintf("%s %s\n", option, cursor)
		}

	case ScanMenu:
		for i, option := range m.ScanResults {
			cursor := ""
			if m.cursor == i {
				cursor = "<"
			}

			s += fmt.Sprintf("%s %s\n", option, cursor)
		}
	case PairedMenu:
		for i, option := range m.PairedDevices {
			cursor := ""
			if m.cursor == i {
				cursor = "<"
			}

			s += fmt.Sprintf("%s %s\n", option, cursor)
		}

	}

	// footer
	s += "\npress q to quit.\n"

	return s
}
