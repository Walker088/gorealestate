package http

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultUserAgent        = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36 Edg/110.0.1587.69"
	DefaultMaxBody    int64 = 1024 * 1024 * 1024 // 1GB
	DefaultRetryTimes       = 2
)

var (
	DefaultRetryHTTPCodes = []int{500, 502, 503, 504, 522, 524, 408}
)

type Client struct {
	*http.Client
	opt *Options
}

// Options is custom http.client options
type Options struct {
	MaxBodySize    int64
	RetryTimes     int
	RetryHTTPCodes []int
	ProxyFunc      func(*http.Request) (*url.URL, error)
}

func New(opt *Options) *Client {
	var proxyFunction = http.ProxyFromEnvironment
	if opt.ProxyFunc != nil {
		proxyFunction = opt.ProxyFunc
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: proxyFunction,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          0,    // Default: 100
			MaxIdleConnsPerHost:   1000, // Default: 2
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: time.Second * 180, // Google's timeout
	}
	return &Client{
		Client: httpClient,
		opt:    opt,
	}
}
