package sirius

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

const ErrUnauthorized ClientError = "unauthorized"

type ClientError string

func (e ClientError) Error() string {
	return string(e)
}

type StatusError struct {
	Code   int    `json:"code"`
	URL    string `json:"url"`
	Method string `json:"method"`
	Body   string `json:"body"`
}

func newStatusError(resp *http.Response) *StatusError {
	data, _ := io.ReadAll(resp.Body)

	return &StatusError{
		Code:   resp.StatusCode,
		URL:    resp.Request.URL.String(),
		Method: resp.Request.Method,
		Body:   string(data),
	}
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("%s %s returned %d", e.Method, e.URL, e.Code)
}

func (e *StatusError) Title() string {
	return "unexpected response from Sirius"
}

func (e *StatusError) Data() interface{} {
	return e
}

type Context struct {
	Context   context.Context
	Cookies   []*http.Cookie
	XSRFToken string
}

func NewClient(httpClient *http.Client, baseURL string) (*Client, error) {
	return &Client{
		http:    httpClient,
		baseURL: baseURL,
	}, nil
}

type Client struct {
	http    *http.Client
	baseURL string
}

func (c *Client) newRequest(ctx Context, method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx.Context, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	for _, c := range ctx.Cookies {
		req.AddCookie(c)
	}

	req.Header.Add("OPG-Bypass-Membrane", "1")
	req.Header.Add("X-XSRF-TOKEN", ctx.XSRFToken)
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, err
}
