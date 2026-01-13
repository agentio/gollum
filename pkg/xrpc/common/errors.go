package common

import (
	"fmt"
	"net/http"
)

type ErrorWithMessage struct {
	ErrStr  string `json:"error"`
	Message string `json:"message"`
}

func (e *ErrorWithMessage) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrStr, e.Message)
}

type Error struct {
	StatusCode int
	Wrapped    error
}

func (e *Error) Error() string {
	// Preserving "XRPC ERROR %d" prefix for compatibility - previously matching this string was the only way
	// to obtain the status code.
	if e.Wrapped == nil {
		return fmt.Sprintf("XRPC ERROR %d", e.StatusCode)
	}
	return fmt.Sprintf("XRPC ERROR %d: %s", e.StatusCode, e.Wrapped)
}

func (e *Error) Unwrap() error {
	if e.Wrapped == nil {
		return nil
	}
	return e.Wrapped
}

type XRPCError struct {
	ErrStr  string `json:"error"`
	Message string `json:"message"`
}

func (xe *XRPCError) Error() string {
	return fmt.Sprintf("%s: %s", xe.ErrStr, xe.Message)
}

func ErrorFromHTTPResponse(resp *http.Response, err error) error {
	r := &Error{
		StatusCode: resp.StatusCode,
		Wrapped:    err,
	}
	return r
}
