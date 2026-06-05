package main

import (
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"net/http"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	choices  []string
	cursor   int
	count    int
	fs       bool
	fetching bool
	result   string
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

type fetchResult string

func fetch() tea.Msg {
	c := &http.Client{Timeout: 10 * time.Second}

	resp, err := c.Get("https://yourmom.zip:")
	if err != nil {
		return nil
	}

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return fetchResult(rawBody)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.count++
	m.fetching = false
	m.result = ""

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
		case "g":
			m.fetching = true
			return m, fetch
		default:
			m.cursor = 2
			return m, nil
		}
	case fetchResult:
		m.result = string(msg)
	}

	return m, nil
}

func (m model) View() tea.View {
	var s strings.Builder
	s.WriteString("Request-Query Lookup")
	s.WriteString(": ")
	s.WriteString(strconv.Itoa(m.count))

	for _, item := range m.choices {
		s.WriteString("\n")
		s.WriteString(item)
	}

	switch m.cursor {
	case 1:
		s.WriteString("\n\nspace was pressed")
	case 2:
		s.WriteString("\n\nDont know what happened dude")
	}

	if m.fetching {
		s.WriteString("\n\nFetching your mom!")
	}

	if m.result != "" {
		s.WriteString("\n\n\n\n")
		s.WriteString(m.result)
	}

	v := tea.NewView(s.String())

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
