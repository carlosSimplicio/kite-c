package components

import (
	"context"
	"database/sql"
	"io"
	"log"
	"strconv"

	"net/http"
	"strings"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	db "kite-c/database/sqlc"
	commandRepository "kite-c/services/repository"
	"kite-c/services/rest"
)

const REST_COMMAND_TYPE_ID int64 = 1

type RestWindowModel struct {
	inputTextArea    textarea.Model
	responseTextArea textarea.Model
	restCommand      *services.ParsedRestCommand
	restUrl          string
	fetchResult      string

	width  int
	height int

	httpClient *http.Client
	repository *commandRepository.CommandRepository
}

func customPromptFunc(info textarea.PromptInfo) string {
	return (strconv.Itoa(info.LineNumber) + " ")
}

func NewRestWindowModel(repository *commandRepository.CommandRepository) RestWindowModel {
	ita := textarea.New()
	ita.SetVirtualCursor(true)
	ita.ShowLineNumbers = false
	ita.SetPromptFunc(4, customPromptFunc)
	ita.SetWidth(60)
	ita.SetValue("GET https://jsonplaceholder.typicode.com/todos")

	itaStyle := textarea.StyleState{
		Base: lipgloss.NewStyle().Background(lipgloss.Color("1")).BorderBackground(lipgloss.Color("2")),
	}

	ita.SetStyles(textarea.Styles{Focused: itaStyle, Blurred: itaStyle})

	rta := textarea.New()
	rta.SetVirtualCursor(true)
	rta.ShowLineNumbers = false
	rta.SetPromptFunc(4, customPromptFunc)
	rta.SetWidth(60)
	rtaStyle := textarea.StyleState{
		Base: lipgloss.NewStyle().Background(lipgloss.Color("2")),
	}

	rta.SetStyles(textarea.Styles{Focused: rtaStyle, Blurred: rtaStyle})

	return RestWindowModel{
		inputTextArea:    ita,
		responseTextArea: rta,
		httpClient:       &http.Client{},
		repository:       repository,
	}

}

func (m RestWindowModel) Init() tea.Cmd {
	return nil
}

type FetchResult struct {
	success bool
	data    string
}

func (m RestWindowModel) doRequest() tea.Msg {
	req, err := http.NewRequest(string(m.restCommand.Method), m.restCommand.Url, nil)
	if err != nil {
		return FetchResult{
			success: false,
			data:    err.Error(),
		}
	}
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return FetchResult{
			success: false,
			data:    err.Error(),
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return FetchResult{
			success: false,
			data:    err.Error(),
		}
	}

	return FetchResult{
		success: true,
		data:    string(body),
	}

}

func (m *RestWindowModel) parseCommand(command string) {
	parsedCommand, err := services.ParseRestCommand(command)
	if err != nil {
		return
	}

	m.restCommand = parsedCommand
}

func (m *RestWindowModel) saveCommand() tea.Msg {
	log.Println("Trying to save command")
	cmd, err := m.repository.CreateCommand(context.TODO(), db.CreateCommandParams{
		Name: "This is a name for a test",
		Description: sql.NullString{
			String: "This is a description for a test",
			Valid:  true,
		},
		CommandQuery: m.inputTextArea.Value(),
		TypeID:       REST_COMMAND_TYPE_ID,
	})

	if err != nil {
		log.Printf("Failed to save command: %s\n", err.Error())
		return nil
	}
	log.Printf("Created command: %v\n", cmd)

	return nil
}

func (m RestWindowModel) Update(msg tea.Msg) (RestWindowModel, tea.Cmd) {
	var cmd tea.Cmd
	log.Printf("update rest window: %v\n", msg)
	log.Printf("msg type: %v\n", msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyPressMsg:
		log.Println("detected key press msg")
		switch msg.Key().String() {
		case "ctrl+f":
			m.responseTextArea.Blur()
			m.inputTextArea.Focus()
			return m, cmd
		case "ctrl+r":
			m.inputTextArea.Blur()
			m.responseTextArea.Focus()
			return m, cmd
		case "ctrl+a":
			m.inputTextArea.Blur()
			m.responseTextArea.Blur()
			return m, cmd
		case "ctrl+s":
			return m, m.saveCommand
		case "alt+enter":
			if !m.inputTextArea.Focused() {
				break
			}
			m.inputTextArea.Blur()
			m.responseTextArea.Blur()
			m.parseCommand(m.inputTextArea.Value())
			return m, m.doRequest
		}

	case FetchResult:
		log.Println("detected fetch result msg")
		log.Printf("fetch result: %v", msg)
		m.fetchResult = msg.data
		m.responseTextArea.SetValue(msg.data)
		m.responseTextArea.MoveToBegin()
		return m, cmd
	}

	if m.inputTextArea.Focused() {
		m.inputTextArea, cmd = m.inputTextArea.Update(msg)
	}

	if m.responseTextArea.Focused() {
		m.responseTextArea, cmd = m.responseTextArea.Update(msg)
	}

	return m, cmd
}

func (m RestWindowModel) View() string {
	header := "Rest Mode"
	requestWindowTitle := "Request"
	responseWindowTitle := "Response"

	headerHeight := lipgloss.Height(header)
	requestTitleHeight := lipgloss.Height(requestWindowTitle)
	responseTitleHeight := lipgloss.Height(responseWindowTitle)
	inputHeights := (m.height - headerHeight - requestTitleHeight - responseTitleHeight - 6) / 2

	m.inputTextArea.SetHeight(inputHeights)
	m.inputTextArea.SetWidth(m.width - 5)
	m.responseTextArea.SetHeight(inputHeights)
	m.responseTextArea.SetWidth(m.width - 5)

	inputTextArea := m.inputTextArea.View()

	responseTextArea := m.responseTextArea.View()
	restWindowStyle := lipgloss.NewStyle().Width(m.width - 2)

	content := restWindowStyle.Render(
		strings.Join([]string{
			header,
			"Request",
			inputTextArea,
			"Response",
			responseTextArea,
		}, "\n\n"),
	)

	return content
}
