package client

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func NewClient(userAgent string) *Client {
	return &Client{
		userAgent: userAgent,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					// HTTPの時は無視されるだけ
					ServerName: ShareTargetURLs.TargetHost,
				},
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return fmt.Errorf("redirect attempted")
			},
		},
	}
}

func NewClientForInitialize() *Client {
	return &Client{
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					// HTTPのときには無視される
					ServerName: ShareTargetURLs.TargetHost,
				},
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return fmt.Errorf("redirect attempted")
			},
		},
	}
}
