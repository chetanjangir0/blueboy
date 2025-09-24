package main

//todo
// wifi on off function
// use ticks to update frequently
// s for scan
// make it show security and network strength of devices
// styles like loading spinner

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/blueboy/internal/ui"
	// "os"
)

func main() {
	// debug file
	// f, err := os.Create("debug.log")
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	// log.SetOutput(f)

	program := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
