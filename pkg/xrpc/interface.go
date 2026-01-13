package xrpc

import "context"

type RequestType int

const (
	Query = RequestType(iota)
	Procedure
)

type Client interface {
	Do(context.Context, RequestType, string, string, map[string]interface{}, interface{}, interface{}) error
}
