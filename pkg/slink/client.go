package slink

import "context"

type RequestType int

const (
	Query = RequestType(iota)
	Procedure
)

type Client interface {
	Do(ctx context.Context, requestType RequestType, contentType string, xrpcMethod string, parameters map[string]any, input any, output any) error
}
