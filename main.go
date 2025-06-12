package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	program := tea.NewProgram(initialModel())
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
	tasks  []string
	cursor int
	done   map[int]struct{} // key = index of the tasks slice
}

func initialModel() model {
	return model{
		tasks: []string{"buy vegetables", "go to gym", "brush your teeth"},
		done:  make(map[int]struct{}),
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
			if m.cursor < len(m.tasks)-1 {
				m.cursor++
			}

		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "enter":
			_, ok := m.done[m.cursor]
			if ok {
				delete(m.done, m.cursor)
			} else {
				m.done[m.cursor] = struct{}{}
			}
		}

	}
	return m, nil
}

func (m model) View() string {
	// the header
	s := "Welcome to the todo list app\n\n"

	// the tasks
	for i, task := range m.tasks {
		cursor := ""
		if m.cursor == i {
			cursor = "<"
		}

		checked := " "
		if _, ok := m.done[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("[%s] %s %s\n", checked, task, cursor)
	}

	// footer
	s += "\npress q to quit.\n"

	return s
}
