package main

import (
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	listView uint = iota
	titleView
	bodyView
)

type model struct {
	state       uint
	store       *Store
	notes       []Note
	currentNote Note
	listIndex   int
	tesxtarea   textarea.Model
	textinput   textinput.Model
}

func NewModel(store *Store) model {
	notes, err := store.GetNote()
	if err != nil {
		log.Fatalf("Unable to get notes : %v", err)
	}
	return model{
		state:     listView,
		store:     store,
		notes:     notes,
		tesxtarea: textarea.New(),
		textinput: textinput.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)
	m.tesxtarea, cmd = m.tesxtarea.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch m.state {
		case listView:
			switch key {
			case "q":
				return m, tea.Quit
			case "n":
				m.textinput.SetValue("")
				m.textinput.Focus()
				m.currentNote = Note{}
				m.state = titleView

			case "up", "k":
				if m.listIndex > 0 {
					m.listIndex--
				}
			case "down", "j":
				if m.listIndex < len(m.notes)-1 {
					m.listIndex++
				}
			case "enter":
				m.currentNote = m.notes[m.listIndex]
				m.tesxtarea.SetValue(m.currentNote.Body)
				m.tesxtarea.Focus()
				m.tesxtarea.CursorEnd()
				m.state = bodyView
			}
		case titleView:
			switch key {
			case "enter":
				title := m.textinput.Value()
				if title != "" {
					m.currentNote.Title = title

					m.tesxtarea.SetValue("")
					m.tesxtarea.Focus()
					m.tesxtarea.CursorEnd()

					m.state = bodyView
				}
			case "esc":
				m.state = listView
			}

		case bodyView:
			switch key {
			case "ctrl+s":
				body := m.tesxtarea.Value()
				m.currentNote.Body = body
				var err error
				if err = m.store.SaveNote(m.currentNote); err != nil {
					return m, tea.Quit
				}
				m.notes, err = m.store.GetNote()
				if err != nil {
					return m, tea.Quit
				}
				m.currentNote = Note{}
				m.state = listView
			case "esc":
				m.state = listView
			}
		}
	}
	return m, tea.Batch(cmds...)
}
