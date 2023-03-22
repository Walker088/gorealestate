package error

import "fmt"

type Error struct {
	ErrorData
}

type InnerError struct {
	ErrorData
}

type ErrorData struct {
	Code    string       `json:"code,omitempty"`
	Message string       `json:"message,omitempty"`
	Target  string       `json:"target,omitempty"`
	Details []Error      `json:"details,omitempty"`
	Inner   *interface{} `json:"innererror,omitempty"`
}

func NewError(errorData ErrorData) *Error {
	return &Error{
		ErrorData: errorData,
	}
}

func NewInnerError(errorData ErrorData) *InnerError {
	return &InnerError{
		ErrorData: errorData,
	}
}

func NewErrorData(code string, message string, target string, details []Error, inner *interface{}) *ErrorData {
	return &ErrorData{
		Code:    code,
		Message: message,
		Target:  target,
		Details: details,
		Inner:   inner,
	}
}

func (e *Error) ToString() string {
	return fmt.Sprintf(
		"Error{code=%s, message=%s, target=%s, details=%v, innererror=%v}",
		e.Code,
		e.Message,
		e.Target,
		e.Details,
		e.Inner,
	)
}

func (e *InnerError) ToString() string {
	return fmt.Sprintf(
		"Error{code=%s, message=%s, target=%s, details=%v, innererror=%v}",
		e.Code,
		e.Message,
		e.Target,
		e.Details,
		e.Inner,
	)
}

func (e *ErrorData) ToString() string {
	return fmt.Sprintf(
		"Error{code=%s, message=%s, target=%s, details=%v, innererror=%v}",
		e.Code,
		e.Message,
		e.Target,
		e.Details,
		e.Inner,
	)
}
