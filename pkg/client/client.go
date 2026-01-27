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
	"github.com/agentio/slink/pkg/slink"
	"github.com/charmbracelet/log"
)

type Client struct {
	Host          string
	Authorization string
	ATProtoProxy  string
	ProxySession  string
	UserDid       string
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

func (c *Client) SetSessionHeaders(r *http.Request) *Client {
	c.ProxySession = r.Header.Get("proxy-session")
	c.UserDid = r.Header.Get("user-did")
	return c
}

type ClientOptions struct {
	Host          string
	Authorization string
	ATProtoProxy  string
	ProxySession  string
	UserDid       string
}

func NewClientWithOptions(options ClientOptions) *Client {
	return &Client{
		Host:          options.Host,
		Authorization: options.Authorization,
		ATProtoProxy:  options.ATProtoProxy,
		ProxySession:  options.ProxySession,
		UserDid:       options.UserDid,
	}
}

func (c *Client) Do(
	ctx context.Context,
	kind slink.RequestType,
	contentType string,
	method string,
	params map[string]any,
	bodyvalue any,
	out any,
) error {
	var body io.Reader
	if bodyvalue != nil {
		if bodyreader, ok := bodyvalue.(io.Reader); ok {
			body = bodyreader
		} else {
			b, err := json.Marshal(bodyvalue)
			if err != nil {
				return err
			}
			body = bytes.NewReader(b)
		}
	}

	var m string
	switch kind {
	case slink.Query:
		m = "GET"
	case slink.Procedure:
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

	if bodyvalue != nil && contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("User-Agent", "slink")

	authorization := c.Authorization
	if authorization == "" {
		authorization = os.Getenv("SLINK_AUTH")
	}
	if authorization != "" {
		req.Header.Set("authorization", authorization)
		log.Infof("authorization: %s", slink.TruncateToLength(authorization, 16))
	}

	atprotoproxy := c.ATProtoProxy
	if atprotoproxy == "" {
		atprotoproxy = os.Getenv("SLINK_ATPROTOPROXY")
	}
	if atprotoproxy != "" {
		req.Header.Set("atproto-proxy", atprotoproxy)
		log.Infof("atproto-proxy: %s", atprotoproxy)
	}

	proxysession := c.ProxySession
	if proxysession == "" {
		proxysession = os.Getenv("SLINK_PROXYSESSION")
	}
	if proxysession != "" {
		req.Header.Set("proxy-session", proxysession)
		log.Infof("proxy-session: %s", proxysession)
	}

	userdid := c.UserDid
	if userdid == "" {
		userdid = os.Getenv("SLINK_USERDID")
	}
	if userdid != "" {
		req.Header.Set("user-did", userdid)
		log.Infof("user-did: %s", userdid)
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
		return xrpcErrorFromResponse(resp, b)
	}

	if out == nil {
		return nil
	}

	if outBytesPointer, ok := out.(*[]byte); ok {
		*outBytesPointer = b
		return nil
	}

	responseContentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(responseContentType, "application/json") {
		return fmt.Errorf("unexpected content type: %s", responseContentType)
	}

	if err := json.Unmarshal(b, out); err != nil {
		return fmt.Errorf("decoding xrpc response: %w", err)
	}
	return nil

}
