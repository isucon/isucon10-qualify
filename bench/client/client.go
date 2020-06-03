package client

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/morikuni/failure"
	// "github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
)

const (
	DefaultAPITimeout = 10
)

type Client struct {
	userAgent  string
	httpClient *http.Client
}

type TargetURLs struct {
	AppURL     url.URL
	TargetHost string
}

var (
	ShareTargetURLs *TargetURLs
)

func SetShareTargetURLs(appURL, targetHost string) error {
	var err error
	ShareTargetURLs, err = newTargetURLs(appURL, targetHost)
	if err != nil {
		return err
	}

	return nil
}

func newTargetURLs(appURL, targetHost string) (*TargetURLs, error) {
	if len(appURL) == 0 {
		return nil, fmt.Errorf("client: missing url")
	}

	appParsedURL, err := urlParse(appURL)
	if err != nil {
		return nil, failure.Wrap(err, failure.Messagef("failed to parse url: %s", appURL))
	}

	return &TargetURLs{
		AppURL:     *appParsedURL,
		TargetHost: targetHost,
	}, nil
}

func urlParse(ref string) (*url.URL, error) {
	u, err := url.Parse(ref)
	if err != nil {
		return nil, err
	}

	if u.Host == "" {
		return nil, fmt.Errorf("host is empty")
	}

	return &url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
	}, nil
}

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
			Timeout: time.Duration(DefaultAPITimeout) * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return fmt.Errorf("redirect attempted")
			},
		},
	}
}

func (c *Client) newGetRequest(u url.URL, spath string) (*http.Request, error) {
	if len(spath) > 0 {
		u.Path = spath
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Host = ShareTargetURLs.TargetHost
	req.Header.Set("User-Agent", c.userAgent)

	return req, nil
}

func (c *Client) newGetRequestWithQuery(u url.URL, spath string, q url.Values) (*http.Request, error) {
	if len(spath) > 0 {
		u.Path = spath
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Host = ShareTargetURLs.TargetHost
	req.Header.Set("User-Agent", c.userAgent)

	return req, nil
}

func (c *Client) newPostRequest(u url.URL, spath string, body io.Reader) (*http.Request, error) {
	u.Path = spath

	req, err := http.NewRequest(http.MethodPost, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Host = ShareTargetURLs.TargetHost
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	return req, nil
}

func checkStatusCode(res *http.Response, expectedStatusCode int) error {
	prefixMsg := fmt.Sprintf("%s %s", res.Request.Method, res.Request.URL.Path)

	if res.StatusCode != expectedStatusCode {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return failure.Wrap(err, failure.Message(prefixMsg+": bodyの読み込みに失敗しました"))
		}
		return failure.Translate(
			fmt.Errorf("status code: %d; body: %s", res.StatusCode, b),
			fails.ErrApplication,
			failure.Messagef("%s: got response status code %d; expected %d", prefixMsg, res.StatusCode, expectedStatusCode),
		)
	}

	return nil
}

func checkStatusCodeWithMsg(res *http.Response, expectedStatusCode int, msg string) error {
	prefixMsg := fmt.Sprintf("%s %s", res.Request.Method, res.Request.URL.Path)

	if res.StatusCode != expectedStatusCode {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return failure.Wrap(err, failure.Message(prefixMsg+": bodyの読み込みに失敗しました "+msg))
		}
		return failure.Translate(fmt.Errorf("status code: %d; body: %s", res.StatusCode, b), fails.ErrApplication,
			failure.Messagef("%s: got response status code %d; expected %d %s", prefixMsg, res.StatusCode, expectedStatusCode, msg),
		)
	}

	return nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	res, err := c.httpClient.Do(req)
	if err != nil {
		if nerr, ok := err.(net.Error); ok {
			if nerr.Timeout() {
				return nil, failure.Translate(err, fails.ErrTimeout)
			} else if nerr.Temporary() {
				return nil, failure.Translate(err, fails.ErrTemporary)
			}
		}

		return nil, err
	}

	return res, nil
}
