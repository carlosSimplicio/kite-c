package components

import (
	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"io"
	"log"
	"net/http"
	"strings"
)

type RestWindowModel struct {
	inputTextArea    textarea.Model
	responseTextArea textarea.Model
	restMethod       string
	restUrl          string
	fetchResult      string
}

func NewRestWindowModel() RestWindowModel {
	ita := textarea.New()
	ita.ShowLineNumbers = false
	ita.SetVirtualCursor(true)
	// https://jsonplaceholder.typicode.com/todos
	ita.SetValue("GET https://jsonplaceholder.typicode.com/todos")

	rta := textarea.New()
	rta.ShowLineNumbers = false
	rta.SetVirtualCursor(true)

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
	case tea.KeyPressMsg:
		log.Println("detected key press msg")
		switch msg.Key().String() {
		case "ctrl+f":
			m.inputTextArea.Focus()
			return m, cmd
		case "ctrl+r":
			m.responseTextArea.Focus()
			return m, cmd
		case "ctrl+a":
			m.inputTextArea.Blur()
			m.responseTextArea.Blur()
			return m, cmd
		case "alt+enter":
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

	inputTextArea := m.inputTextArea.View()

	m.responseTextArea.CursorStart()
	responseTextArea := m.responseTextArea.View()

	return strings.Join([]string{
		header,
		inputTextArea,
		responseTextArea,
	}, "\n")
}
