package services

import (
	"errors"
	"strings"
)

type HttpMethod string

const (
	POST   HttpMethod = "POST"
	GET    HttpMethod = "GET"
	DELETE HttpMethod = "DELETE"
	PUT    HttpMethod = "PUT"
)

type ParsedRestCommand struct {
	Method  HttpMethod
	Url     string
	Headers []string
}

// POST https://example.com/1
// {
// 	Content-Type: application/json
// }
// {
// 	"name": "foo",
// 	"age": 30
// }

func ParseRestCommand(command string) (*ParsedRestCommand, error) {
	splittedCommand := strings.Split(command, " ")
	if len(splittedCommand) < 2 {
		return &ParsedRestCommand{}, errors.New("Command lenght should be gte 2")
	}

	method := splittedCommand[0]
	if isValidMethod(method) == false {
		return nil, errors.New("Invalid HTTP method")
	}

	url := splittedCommand[1]

	return &ParsedRestCommand{
		Method: HttpMethod(method),
		Url:    url,
	}, nil
}

func isValidMethod(s string) bool {
	switch HttpMethod(s) {
	case POST, GET, DELETE, PUT:
		return true
	default:
		return false
	}

}
