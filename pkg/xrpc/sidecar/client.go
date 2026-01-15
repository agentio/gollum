package xrpc_sidecar

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/agentio/sidecar"
	"github.com/agentio/slink/pkg/xrpc"
	"github.com/agentio/slink/pkg/xrpc/common"
)

type Client struct {
	Host string
}

func NewClient() *Client {
	return &Client{
		Host: "http://localhost:5050",
	}
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

	host := c.Host
	if strings.HasPrefix(host, "unix:") {
		host = "http://unix"
	}
	path := host + "/xrpc/" + method + paramStr
	req, err := http.NewRequest(m, path, body)
	if err != nil {
		return err
	}

	if bodyobj != nil && contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("User-Agent", "slink-sidecar")

	client := sidecar.NewClient(sidecar.ClientOptions{
		Address: c.Host,
	})

	resp, err := client.HttpClient.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {

		b, _ := io.ReadAll(resp.Body)
		log.Printf("%s", string(b))

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
