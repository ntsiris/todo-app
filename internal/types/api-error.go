package types

import "fmt"

type APIError struct {
	Code          int    `json:"code"`
	Message       string `json:"message"`
	Operation     string `json:"operation"`
	EmbeddedError error  `json:"embeddedError"`
}

func (e *APIError) Error() string {
	retStr := fmt.Sprintf("API error code: %d, message: %s", e.Code, e.Message)
	if e.EmbeddedError != nil {
		retStr = fmt.Sprintf("%s\n%s", retStr, e.EmbeddedError)
	}

	return retStr
}

func FormatOperation(method, path string) string {
	return fmt.Sprintf("%s %s", method, path)
}
