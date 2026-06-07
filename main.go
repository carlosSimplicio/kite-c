package main

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"net/http"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type styles struct {
	window  lipgloss.Style
	overlay lipgloss.Style
}

type Model struct {
	panels        []string
	selectedPanel string
	styles        *styles
	width         int
	height        int

	searchActive bool
	searchString string
	textInput    textinput.Model

	fetching bool
	result   string
}

func initialModel() Model {
	s := new(styles)
	s.window = lipgloss.NewStyle().Padding(2).Align(lipgloss.Center).Background(lipgloss.Color("1"))
	s.overlay = lipgloss.NewStyle().Padding(2).Width(20).Height(5)
	textInput := textinput.New()
	textInput.Placeholder = "Busca"
	textInput.Focus()
	textInput.SetWidth(20)

	return Model{
		panels:        []string{"Rest Client", "SQL Client"},
		selectedPanel: "Rest Client",
		searchActive:  false,
		styles:        s,
		textInput:     textInput,
	}
}

func (m *Model) setSelectedPanel(msg tea.KeyPressMsg) {
	splitted := strings.Split(msg.String(), "")
	panelNumberStr := splitted[len(splitted)-1]
	panelNumber, err := strconv.Atoi(panelNumberStr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if len(m.panels) < panelNumber {
		log.Println("panel number is greater than available panels")
		return
	}

	m.selectedPanel = m.panels[panelNumber-1]
	log.Printf("selectedPanel: %s\n", m.selectedPanel)
}

func (m Model) Init() tea.Cmd {
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("message: %v", msg)
	m.fetching = false
	m.searchString = ""

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		log.Printf("width: %d, height: %d", msg.Width, msg.Height)
		m.styles.window = m.styles.window.Width(msg.Width)
		m.styles.window = m.styles.window.Height(msg.Height)
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			m.searchActive = false
			m.searchString = m.textInput.Value()
			m.textInput.Reset()
			return m, nil
		case "ctrl+q":
			return m, tea.Quit
		case "g":
			if !m.searchActive {
				m.fetching = true
				return m, fetch
			}
		case "alt+1":
			m.setSelectedPanel(msg)
			log.Printf("model: %v", m)
			return m, nil
		case "alt+2":
			m.setSelectedPanel(msg)
			log.Printf("model: %v", m)
			return m, nil
		case "ctrl+k":
			m.searchActive = true
			return m, nil
		case "esc":
			m.searchActive = false
			return m, nil
		default:
			log.Printf("%s was pressed", msg)
		}
	case fetchResult:
		m.result = string(msg)
	}

	var cmd tea.Cmd
	if m.searchActive {
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() tea.View {
	var s strings.Builder
	s.WriteString("Kite Client")
	s.WriteString("\n")

	s.WriteString(m.selectedPanel)
	s.WriteString("\n")

	if m.fetching {
		s.WriteString("\n\nFetching your mom!")
	}

	if m.result != "" {
		s.WriteString("\n\n\n\n")
		s.WriteString(m.result)
	}

	if m.searchString != "" {
		s.WriteString("\n\n\n")
		s.WriteString("You just typed: ")
		s.WriteString(m.searchString)
	}

	window := m.styles.window.Render(s.String())
	windowLayer := lipgloss.NewLayer(window)

	log.Printf("busca ativa: %t\n", m.searchActive)
	if m.searchActive {
		overlay := m.styles.overlay.Render(m.textInput.View())
		overlayLayer := lipgloss.NewLayer(overlay).X(m.width / 2).Y(m.height / 2).Z(1)
		windowLayer.AddLayers(overlayLayer)
	}

	compositor := lipgloss.NewCompositor(windowLayer)
	v := tea.NewView(compositor.Render())
	v.AltScreen = true

	return v
}

func main() {
	logOutputFile := "logs.txt"
	file, err := os.OpenFile(logOutputFile, os.O_TRUNC|os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	log.SetOutput(file)

	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Failed to run program %s", err.Error())
	}
}
