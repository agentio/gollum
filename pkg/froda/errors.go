package froda

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type XRPCError struct {
	Code    int    `json:"code"`
	Title   string `json:"error"`
	Message string `json:"message"`
}

func (e *XRPCError) Error() string {
	return fmt.Sprintf("XRPC ERROR %d: %s (%s)", e.Code, e.Title, e.Message)
}

func xrpcErrorFromResponse(resp *http.Response, b []byte) error {
	var xrpcError XRPCError
	if err := json.Unmarshal(b, &xrpcError); err != nil {
		xrpcError.Title = "failed to decode xrpc error message"
		xrpcError.Message = strings.TrimSpace(string(b))
	}
	xrpcError.Code = resp.StatusCode
	return &xrpcError
}
