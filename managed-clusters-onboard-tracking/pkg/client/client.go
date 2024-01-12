package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	ServiceName    string
	BaseURL        string
	SecureApiToken string
}

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
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return err
	}
	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return err
		}
	}
	return nil
}
