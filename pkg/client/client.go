package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/agentio/sidecar"
	"github.com/agentio/slink/pkg/common"
	"github.com/charmbracelet/log"
)

type Client struct {
	Host string
}

func NewClient() *Client {
	host := os.Getenv("SLINK_HOST")
	if host == "" {
		host = "http://localhost:5050"
	}
	return &Client{
		Host: host,
	}
}

func (c *Client) Do(
	ctx context.Context,
	kind common.RequestType,
	contentType string,
	method string,
	params map[string]any,
	bodyobj any,
	out any,
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
	case common.Query:
		m = "GET"
	case common.Procedure:
		m = "POST"
	default:
		return fmt.Errorf("unsupported request kind: %d", kind)
	}

	var paramStr string
	if len(params) > 0 {
		paramStr = "?" + makeParams(params)
	}

	host := c.Host
	if strings.HasPrefix(host, "unix:") {
		host = "http://unix"
	}
	path := host + "/xrpc/" + method + paramStr
	log.Infof("%s %s", m, path)
	req, err := http.NewRequest(m, path, body)
	if err != nil {
		return err
	}

	if bodyobj != nil && contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("User-Agent", "slink")

	authorization := os.Getenv("SLINK_AUTH")
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
		log.Infof("%s", common.TruncateToLength(authorization, 16))
	}

	client := sidecar.NewClient(sidecar.ClientOptions{
		Address: c.Host,
	})

	resp, err := client.HttpClient.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Infof("%d (%d bytes)", resp.StatusCode, len(b))
	if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		log.Debugf("%s", string(b))
	}

	if resp.StatusCode != 200 {
		var xe XRPCError
		if err := json.Unmarshal(b, &xe); err != nil {
			return errorFromHTTPResponse(resp, fmt.Errorf("failed to decode xrpc error message: %w", err))
		}
		return errorFromHTTPResponse(resp, &xe)
	}

	if out != nil {
		bufferPointer, ok := out.(*[]byte)
		if ok {
			*bufferPointer = b
			return nil
		}
		responseContentType := resp.Header.Get("Content-Type")
		if strings.HasPrefix(responseContentType, "application/json") {
			if err := json.Unmarshal(b, out); err != nil {
				return fmt.Errorf("decoding xrpc response: %w", err)
			}
			return nil
		}
		return fmt.Errorf("unexpected content type: %s", responseContentType)
	}
	return nil
}
