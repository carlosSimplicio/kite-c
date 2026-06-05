package main

import (
	"log"
	"strconv"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	choices []string
	cursor  int
	count   int
	fs      bool
}

func initialModel() model {
	return model{
		choices: []string{"Search", "Add request", "Add query"},
		fs:      false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.count++
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "space":
			m.cursor = 1
			return m, nil
		case "q":
			return m, tea.Quit
		case "f":
			m.fs = !m.fs
			return m, nil
		default:
			m.cursor = 2
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	s := "Request-Query Lookup"
	s += ": " + strconv.Itoa(m.count)

	for _, item := range m.choices {
		s += "\n" + item
	}

	switch m.cursor {
	case 1:
		s += "\n\nspace was pressed"
	case 2:
		s += "\n\nDont know what happened dude"
	}

	v := tea.NewView(s)
	if m.fs {
		v.AltScreen = true
	}
	return v
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Failed to run program %s", err.Error())
	}
}
