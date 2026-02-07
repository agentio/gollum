package froda

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"
)

type httpClient struct {
	Host       string
	HttpClient *http.Client
}

type httpClientOptions struct {
	Address  string
	Insecure bool
	Headers  []string
}

// NewClient creates a client representation from an address.
// Addresses must be in the format "HOSTNAME:PORT" or "unix:@SOCKET".
// Connections to port 443 use TLS. All others are cleartext (http1 or h2c).
func newHTTPClient(options httpClientOptions) *httpClient {
	// Expect TLS on port 443 and use the default HTTP client.
	if strings.HasSuffix(options.Address, ":443") {
		return (&httpClient{
			Host: "https://" + options.Address,
			HttpClient: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: options.Insecure},
				},
			},
		})
	}
	protocols := new(http.Protocols)
	protocols.SetUnencryptedHTTP2(true) // Enable h2c (HTTP/2 cleartext)
	protocols.SetHTTP1(true)            // Explicitly allow HTTP/1.1
	protocols.SetHTTP2(false)           // Explicitly disallow encrypted HTTP/2 (HTTPS)
	// If required, create a client that can call unix sockets.
	if strings.HasPrefix(options.Address, "unix:") {
		address := strings.TrimPrefix(options.Address, "unix:")
		return (&httpClient{
			Host: "http://socket", // The name "socket" is arbitrary.
			HttpClient: &http.Client{
				Transport: &http.Transport{
					Protocols: protocols,
					DialContext: func(ctx context.Context, _ string, _ string) (net.Conn, error) {
						return net.DialTimeout("unix", address, 5*time.Second)
					},
				},
			},
		})
	}
	// Create a client for networked h2c connections.
	return (&httpClient{
		Host: "http://" + options.Address,
		HttpClient: &http.Client{
			Transport: &http.Transport{
				Protocols: protocols,
			},
		},
	})
}
