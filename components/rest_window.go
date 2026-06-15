package components

import (
	"io"
	"log"
	"strconv"

	"net/http"
	"strings"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type RestWindowModel struct {
	inputTextArea    textarea.Model
	responseTextArea textarea.Model
	restMethod       string
	restUrl          string
	fetchResult      string

	width  int
	height int
}

func customPromptFunc(info textarea.PromptInfo) string {
	return (strconv.Itoa(info.LineNumber) + " ")
}

func NewRestWindowModel() RestWindowModel {
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
	client := &http.Client{}
	req, err := http.NewRequest(m.restMethod, m.restUrl, nil)
	if err != nil {
		return FetchResult{
			success: false,
			data:    err.Error(),
		}
	}
	resp, err := client.Do(req)
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
	splittedCommand := strings.Split(command, " ")
	m.restMethod = splittedCommand[0]
	m.restUrl = splittedCommand[1]
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
