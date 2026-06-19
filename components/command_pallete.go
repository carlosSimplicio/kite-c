package components

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"context"
	commandRepository "kite-c/services/repository"
	"log"
	"strings"
)

type styles struct {
	overlay lipgloss.Style
}

type CommandPalleteModel struct {
	textInput    textinput.Model
	searchActive bool
	SearchString string
	styles       *styles

	repository   *commandRepository.CommandRepository
	searchResult []SearchResultRow
}

func NewCommandPalleteModel(repository *commandRepository.CommandRepository) CommandPalleteModel {
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
		repository:   repository,
	}
}

func (m CommandPalleteModel) Init() tea.Cmd {
	return nil
}

type SearchResultRow struct {
	title string
	desc  string
}

func (r SearchResultRow) FilterValue() string {
	return ""
}

type SearchResult struct {
	result  []SearchResultRow
	success bool
}

func (m *CommandPalleteModel) searchCommand() tea.Msg {
	if m.textInput.Value() == "" {
		return SearchResult{
			result:  nil,
			success: false,
		}
	}

	result, err := m.repository.SearchCommand(context.TODO(), m.textInput.Value())
	if err != nil {
		log.Printf("Search command failed: %s\n", err.Error())
		return SearchResult{
			result:  nil,
			success: false,
		}
	}

	searchResultRows := make([]SearchResultRow, len(result))
	for index, row := range result {
		searchResultRows[index] = SearchResultRow{
			title: row.Name,
			desc:  row.Description,
		}
	}

	return SearchResult{
		result:  searchResultRows,
		success: true,
	}
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
			m.searchResult = nil
			return m, nil
		}
	case SearchResult:
		switch msg.success {
		case true:
			m.searchResult = msg.result
			return m, nil
		default:
			m.searchResult = nil
			return m, nil
		}
	}

	var cmd tea.Cmd
	if m.searchActive {
		cmds := make([]tea.Cmd, 2)
		m.textInput, cmd = m.textInput.Update(msg)
		log.Printf("text input cmd: %v", cmd)

		cmds = append(cmds, m.searchCommand)
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m CommandPalleteModel) View() string {
	log.Printf("current search result: %v", m.searchResult)
	if m.searchActive == false {
		return ""
	}

	s := strings.Builder{}
	s.WriteString(m.textInput.View())
	s.WriteString("\n")

	if len(m.searchResult) > 0 {
		itemsList := make([]list.Item, len(m.searchResult))
		for index, value := range m.searchResult {
			itemsList[index] = value
		}

		resultList := list.New(itemsList, list.NewDefaultDelegate(), 0, 0)
		s.WriteString(resultList.View())
	}

	return m.styles.overlay.Render(s.String())
}
