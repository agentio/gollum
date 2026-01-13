package xrpc_session

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/agentio/gollum/pkg/xrpc"
	"github.com/agentio/gollum/pkg/xrpc/common"
)

type Client struct {
	UserAgent *string
	Headers   map[string]string
	Session   string
	Did       string
	Proxy     string
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) AuthorizedCopy(r *http.Request) *Client {
	nc := *c
	nc.Session = r.Header.Get("proxy-session")
	nc.Did = r.Header.Get("user-did")
	return &nc
}

func (c *Client) Do(ctx context.Context, kind xrpc.RequestType, contentType, method string, params map[string]any, bodyobj any, out any) error {
	var verb string
	switch kind {
	case xrpc.Query:
		verb = "GET"
	case xrpc.Procedure:
		verb = "POST"
	default:
		return fmt.Errorf("unsupported request kind: %d", kind)
	}

	var paramStr string
	if len(params) > 0 {
		paramStr = "?" + common.MakeParams(params)
	}

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
	ustr := c.Proxy + "/xrpc/" + method + paramStr
	req, err := http.NewRequest(verb, ustr, body)
	if err != nil {
		return err
	}

	if bodyobj != nil && contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("User-Agent", "atiquette-session")

	if c.Headers != nil {
		for k, v := range c.Headers {
			req.Header.Set(k, v)
		}
	}

	req.Header.Set("proxy-session", c.Session)
	req.Header.Set("user-did", c.Did)

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var xe common.ErrorWithMessage
		if err := json.NewDecoder(resp.Body).Decode(&xe); err != nil {
			return common.ErrorFromHTTPResponse(resp, fmt.Errorf("failed to decode xrpc error message: %w", err))
		}

		return common.ErrorFromHTTPResponse(resp, &xe)
	}

	if out != nil {
		if buf, ok := out.(*bytes.Buffer); ok {
			if resp.ContentLength < 0 {
				_, err := io.Copy(buf, resp.Body)
				if err != nil {
					return fmt.Errorf("reading response body: %w", err)
				}
			} else {
				n, err := io.CopyN(buf, resp.Body, resp.ContentLength)
				if err != nil {
					return fmt.Errorf("reading length delimited response body (%d < %d): %w", n, resp.ContentLength, err)
				}
			}
		} else {
			if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
				return fmt.Errorf("decoding xrpc response: %w", err)
			}
		}
	}

	return nil

}
