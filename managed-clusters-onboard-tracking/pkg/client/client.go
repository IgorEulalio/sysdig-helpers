package client

import (
	"IgorEulalio/sysdig-helpers/managed-clusters-onboard-tracking/pkg/logging"
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
	MaxRetries     int
}

const (
	maxElapsedTime = 10
	maxInterval    = 500
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
	attempt := 0

	operation := func() error {
		attempt++
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			logging.Log.Debugf("successfully called endpoint %s with status code %d", req.URL, resp.StatusCode)
			if v != nil {
				return json.NewDecoder(resp.Body).Decode(v)
			}
			return nil
		}

		if resp.StatusCode == http.StatusTooManyRequests && attempt < c.MaxRetries {
			logging.Log.Errorf("Attempt %d calling API %s: Received status code 429. Retrying...", attempt, req.URL)
			return fmt.Errorf("you were trothled in API %s. status code: %d", req.URL, resp.StatusCode)
		}

		return backoff.Permanent(fmt.Errorf("status code: %d", resp.StatusCode))
	}

	// Exponential backoff configuration
	expBackOff := backoff.NewExponentialBackOff()
	expBackOff.MaxElapsedTime = maxElapsedTime * time.Second
	expBackOff.MaxInterval = maxInterval * time.Millisecond

	// Retry with exponential backoff
	return backoff.Retry(operation, backoff.WithMaxRetries(expBackOff, uint64(c.MaxRetries)))
}
