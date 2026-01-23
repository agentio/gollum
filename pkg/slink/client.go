package slink

import "context"

type RequestType int

const (
	Query = RequestType(iota)
	Procedure
)

type Client interface {
	Do(context.Context, RequestType, string, string, map[string]any, any, any) error
}
