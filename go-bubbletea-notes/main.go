package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	store := &Store{}
	if err := store.Init(); err != nil {
		log.Fatalf("Unable to init Store : ", err)
	}
	m := NewModel(store)

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		log.Fatalf("Unable to run TUI : %v", err)
	}
}
