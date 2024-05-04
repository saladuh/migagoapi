package migagoapi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const BaseEndpoint = "https://api.migadu.com/v1/"

type Client struct {
	User       string
	Token      string
	Endpoint   string
	Domain     string
	HttpClient *http.Client
}

type ErrorRequestStatus struct {
	StatusBody []byte
	StatusCode int
}

func (e ErrorRequestStatus) Error() string {
	return fmt.Sprintf("status code: %d, response body: %s", e.StatusCode, e.StatusBody)
}

// Timeout will be used as the client's timeout duration and will overrule
// any context timeout duration (if longer). When the timeout of the client
// making the request is reached, it is treated as if the context of the
// request has ended.
func NewClient(user, token, endpoint, domain *string, timeout *time.Duration) (*Client, error) {
	if user == nil {
		return nil, errors.New("No username supplied.")
	}

	if token == nil {
		return nil, errors.New("No token supplied.")
	}

	if domain == nil {
		return nil, errors.New("No domain supplied")
	}

	c := Client{
		*user,
		*token,
		BaseEndpoint,
		*domain,
		&http.Client{Timeout: time.Duration(time.Second * 120)},
	}

	if endpoint != nil {
		c.Endpoint = *endpoint
	}

	if timeout != nil {
		c.HttpClient.Timeout = *timeout
	}

	return &c, nil
}

func (c *Client) Get(ctx context.Context, path string) ([]byte, error) {
	url := fmt.Sprintf("%s/domains/%s/%s", c.Endpoint, c.Domain, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	body, err := c.doReq(req)

	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) Post(ctx context.Context, path string, json_body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/domains/%s/%s", c.Endpoint, c.Domain, path)

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, url, bytes.NewBuffer(json_body),
	)
	if err != nil {
		return nil, err
	}
	body, err := c.doReq(req)

	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) Put(ctx context.Context, path string, json_body []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/domains/%s/%s", c.Endpoint, c.Domain, path)

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPut, url, bytes.NewBuffer(json_body),
	)
	if err != nil {
		return nil, err
	}
	body, err := c.doReq(req)

	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) Delete(ctx context.Context, path string) ([]byte, error) {
	url := fmt.Sprintf("%s/domains/%s/%s", c.Endpoint, c.Domain, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, http.NoBody)

	if err != nil {
		return nil, err
	}
	body, err := c.doReq(req)

	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) commonMethod(
	ctx context.Context, path string,
	reader_body io.Reader, method string) ([]byte, error) {
	url := fmt.Sprintf("%s/domains/%s/%s", c.Endpoint, c.Domain, path)

	req, err := http.NewRequestWithContext(ctx, method, url, reader_body)

	if err != nil {
		return nil, err
	}
	body, err := c.doReq(req)

	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) doReq(request *http.Request) ([]byte, error) {
	request.SetBasicAuth(c.User, c.Token)
	request.Header.Add("Content-Type", "application/json")
	resp, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	jsonBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &ErrorRequestStatus{jsonBody, resp.StatusCode}
	}

	return jsonBody, nil
}
