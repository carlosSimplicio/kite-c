package components

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"log"
	"strconv"
	"strings"
)

type windowStyles struct {
	window lipgloss.Style
}

type WindowModel struct {
	panels         []string
	selectedPanel  string
	styles         *windowStyles
	width          int
	height         int
	commandPallete CommandPalleteModel
	restModeWindow RestWindowModel
}

func NewWindowModel() WindowModel {
	s := new(windowStyles)
	s.window = lipgloss.NewStyle().
		Padding(2).
		Align(lipgloss.Center).
		Background(lipgloss.Color("#141A24")).
		Border(lipgloss.ThickBorder()).
		BorderBackground(lipgloss.Color("#141A24")).BorderForeground(lipgloss.Color("#8A2BE2"))

	return WindowModel{
		panels:         []string{"Rest Mode", "SQL Mode"},
		selectedPanel:  "Rest Mode",
		styles:         s,
		commandPallete: NewCommandPalleteModel(),
		restModeWindow: NewRestWindowModel(),
	}
}

func (m *WindowModel) setSelectedPanel(msg tea.KeyPressMsg) {
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

func (m WindowModel) Init() tea.Cmd {
	return nil
}

func (m WindowModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("message: %v", msg)

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
		case "ctrl+q", "ctrl+c":
			return m, tea.Quit
		case "alt+1":
			m.setSelectedPanel(msg)
			log.Printf("model: %v", m)
			return m, nil
		case "alt+2":
			m.setSelectedPanel(msg)
			log.Printf("model: %v", m)
			return m, nil
		default:
			log.Printf("%s was pressed", msg)
		}
	}

	var cmd tea.Cmd
	m.commandPallete, cmd = m.commandPallete.Update(msg)
	if cmd != nil {
		return m, cmd
	}

	m.restModeWindow, cmd = m.restModeWindow.Update(msg)
	if cmd != nil {
		return m, cmd
	}

	return m, nil
}

func (m WindowModel) View() tea.View {
	var s strings.Builder
	s.WriteString("Kite Client")
	s.WriteString("\n")

	var modeWindow string
	if m.selectedPanel == "Rest Mode" {
		modeWindow = m.restModeWindow.View()
	} else {
		modeWindow = "SQL Mode"
	}

	s.WriteString(modeWindow)
	s.WriteString("\n")

	if m.commandPallete.SearchString != "" {
		s.WriteString("\n\n\n")
		s.WriteString("You just typed: ")
		s.WriteString(m.commandPallete.SearchString)
	}

	window := m.styles.window.Render(s.String())
	windowLayer := lipgloss.NewLayer(window)

	commandPalleteView := m.commandPallete.View()
	if commandPalleteView != "" {
		overlayLayer := lipgloss.NewLayer(commandPalleteView).
			X((m.width / 2) - lipgloss.Width(commandPalleteView)/2).
			Y(m.height / 2).Z(1)
		windowLayer.AddLayers(overlayLayer)
	}

	compositor := lipgloss.NewCompositor(windowLayer)
	v := tea.NewView(compositor.Render())
	v.AltScreen = true

	return v
}
