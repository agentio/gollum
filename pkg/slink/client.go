package slink

import (
	"context"
	"io"
)

type RequestType int

const (
	Query = RequestType(iota)
	Procedure
)

type Client interface {
	Do(ctx context.Context,
		requestType RequestType,
		contentType string,
		xrpcMethod string,
		parameters map[string]any,
		input any,
		output any,
	) error
	Subscribe(ctx context.Context,
		xrpcMethod string,
		params map[string]any,
		callback func(b io.Reader) error,
	) error
}
