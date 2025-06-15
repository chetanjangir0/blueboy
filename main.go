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
		MainOptions:   []string{"Scan Devices", "Paired Connections"},
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
					m.CurrentMenu = PairedMenu
				}
				m.cursor = 0 // reset cursor pos

			}
		case "b", "esc":
			// out, err := exec.Command("ls").CombinedOutput()
			// if err != nil {
			// 	log.Print("Error", err)
			//
			// }
			// items := strings.Split(string(out), "\n")
			// for _, item := range items {
			// 	log.Println(item)
			// }
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
