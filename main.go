package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	choices []string
	cursor  int
}

func initialModel() model {
	return model{
		choices: []string{"Search", "Add request", "Add query"},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "space":
			m.cursor = 1
			return m, nil
		case "q":
			return m, tea.Quit
		default:
			m.cursor = 2
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	s := "Request-Query Lookup"

	switch m.cursor {
	case 1:
		s += "\n\nspace was pressed"
	case 2:
		s += "\n\nDont know what happened dude"
	}

	return tea.NewView(s)
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Failed to run program %s", err.Error())
	}
}
