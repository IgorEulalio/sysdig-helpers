package client

import (
	"encoding/json"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	ServiceName    string
	BaseURL        string
	SecureApiToken string
}

const (
	maxElapsedTime = 3
	maxInterval    = 1
)

func (c *Client) NewRequest(method string, url *url.URL, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.SecureApiToken)
	return req, nil
}

func (c *Client) CreateUrl(path string, pathParams map[string]string) (*url.URL, error) {
	path = fmt.Sprintf(path, c.BaseURL)
	parsedUrl, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	q := parsedUrl.Query()
	for pathParam, pathParamValue := range pathParams {
		q.Set(pathParam, pathParamValue)
	}
	return parsedUrl, nil
}

func (c *Client) Do(req *http.Request, v interface{}) error {
	operation := func() error {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			if v != nil {
				return json.NewDecoder(resp.Body).Decode(v)
			}
			return nil
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			return backoff.Permanent(fmt.Errorf("status code: %d", resp.StatusCode))
		}

		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	// Exponential backoff configuration
	expBackOff := backoff.NewExponentialBackOff()
	expBackOff.MaxElapsedTime = maxElapsedTime * time.Second
	expBackOff.MaxInterval = maxInterval * time.Second

	return backoff.Retry(operation, expBackOff)
}
