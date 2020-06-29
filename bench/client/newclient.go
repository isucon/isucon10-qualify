package client

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"time"
	"net/http"
)

const (
	NumOfClient = 10
	DefaultAPITimeoutTime = 2000 * time.Millisecond
)

var clients [NumOfClient]*Client

func InitializeClients() {
	for i := 0; i < NumOfClient; i++ {
		userAgent := fmt.Sprintf("isucon-%v-user", i)
		clients[i] = newClient(userAgent)
	}
}

func PickClient() *Client {
	return clients[rand.Intn(len(clients))]
}

func newClient(userAgent string) *Client {
	return &Client{
		userAgent: userAgent,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					// HTTPの時は無視されるだけ
					ServerName: ShareTargetURLs.TargetHost,
				},
			},
			Timeout: DefaultAPITimeoutTime,
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
