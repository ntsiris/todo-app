package types

import (
	"encoding/json"
	"fmt"
)

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

// MarshalJSON implements a custom way for marshaling the API Error type
func (e *APIError) MarshalJSON() ([]byte, error) {
	type Alias APIError // Avoid recursion by using an alias type
	return json.Marshal(&struct {
		*Alias
		EmbeddedError string `json:"embeddedError"`
	}{
		Alias:         (*Alias)(e),
		EmbeddedError: e.getEmbeddedError(),
	})
}

// Helper to unwrap the error chain into a string
func (e *APIError) getEmbeddedError() string {
	if e.EmbeddedError == nil {
		return ""
	}
	return e.EmbeddedError.Error()
}

func FormatOperation(method, path string) string {
	return fmt.Sprintf("%s %s", method, path)
}
