package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/text/cases"
)

func main() {
	fmt.Println("hello world")

}

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{} // key = index of the choices slice
}

func initialModel() model {
	return model{
		choices:  []string{"buy vegetables", "go to gym", "brush your teeth"},
		selected: make(map[int]struct{}),
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
		}

	}
	return m, nil
}

func (m model) View() string {
	return "foo bar"
}
