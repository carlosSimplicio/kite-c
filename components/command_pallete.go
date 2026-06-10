package components

import (
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type styles struct {
	overlay lipgloss.Style
}

type CommandPalleteModel struct {
	textInput    textinput.Model
	searchActive bool
	SearchString string
	styles       *styles
}

func NewCommandPalleteModel() CommandPalleteModel {
	textInput := textinput.New()
	textInput.Placeholder = "Busca"
	textInput.Focus()
	textInput.SetWidth(100)

	s := new(styles)
	s.overlay = lipgloss.NewStyle().Padding(2).Width(100).Height(5)

	return CommandPalleteModel{
		textInput:    textInput,
		SearchString: "",
		styles:       s,
	}
}

func (m CommandPalleteModel) Init() tea.Cmd {
	return nil
}

func (m CommandPalleteModel) Update(msg tea.Msg) (CommandPalleteModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			m.searchActive = false
			m.SearchString = m.textInput.Value()
			m.textInput.Reset()
			return m, nil
		case "ctrl+k":
			m.searchActive = true
			m.SearchString = ""
			return m, nil
		case "esc":
			m.searchActive = false
			m.textInput.Reset()
			return m, nil
		}
	}

	var cmd tea.Cmd
	if m.searchActive {
		m.textInput, cmd = m.textInput.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m CommandPalleteModel) View() string {
	if m.searchActive {
		overlay := m.styles.overlay.Render(m.textInput.View())
		return tea.NewView(overlay).Content
	}

	return ""
}
