package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Password struct {
	isAsking      bool
	passwordInput textinput.Model
}


func NewPasswordModel() Password{
	ti := textinput.New()
	ti.Placeholder = "Enter password"
	ti.CharLimit = 20
	ti.Width = 30
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	return Password{isAsking: false, passwordInput: ti}
}

func (p Password) Update(msg tea.Msg) (Password, tea.Cmd) {
	var cmd tea.Cmd
	p.passwordInput, cmd = p.passwordInput.Update(msg)
	return p, cmd
}

func (p Password) View() string {
	return p.passwordInput.View()
}
