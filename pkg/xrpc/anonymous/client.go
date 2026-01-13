package xrpc_anonymous

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/agentio/slink/pkg/xrpc"
	"github.com/agentio/slink/pkg/xrpc/common"
)

type Client struct {
	Host string
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Do(
	ctx context.Context,
	kind xrpc.RequestType,
	contentType string,
	method string,
	params map[string]interface{},
	bodyobj interface{},
	out interface{},
) error {
	var body io.Reader
	if bodyobj != nil {
		if rr, ok := bodyobj.(io.Reader); ok {
			body = rr
		} else {
			b, err := json.Marshal(bodyobj)
			if err != nil {
				return err
			}

			body = bytes.NewReader(b)
		}
	}

	var m string
	switch kind {
	case xrpc.Query:
		m = "GET"
	case xrpc.Procedure:
		m = "POST"
	default:
		return fmt.Errorf("unsupported request kind: %d", kind)
	}

	var paramStr string
	if len(params) > 0 {
		paramStr = "?" + common.MakeParams(params)
	}

	req, err := http.NewRequest(m, c.Host+"/xrpc/"+method+paramStr, body)
	if err != nil {
		return err
	}

	if bodyobj != nil && contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("User-Agent", "atiquette-anonymous")

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var xe common.XRPCError
		if err := json.NewDecoder(resp.Body).Decode(&xe); err != nil {
			return common.ErrorFromHTTPResponse(resp, fmt.Errorf("failed to decode xrpc error message: %w", err))
		}
		return common.ErrorFromHTTPResponse(resp, &xe)
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("decoding xrpc response: %w", err)
		}
	}

	return nil
}
