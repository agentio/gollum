package common

import (
	"github.com/agentio/slink/pkg/xrpc"
	xrpc_sidecar "github.com/agentio/slink/pkg/xrpc/sidecar"
)

func NewClient() xrpc.Client {
	return xrpc_sidecar.NewClient()
}
