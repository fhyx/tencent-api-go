package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	*TokenHolder
	httpClient *http.Client
	auths      string
	ctype      string
}

var (
	tr = &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
)

func NewClient(urlToken string) *Client {
	hc := &http.Client{Transport: tr}
	return &Client{
		httpClient:  hc,
		TokenHolder: NewTokenHolder(urlToken),
	}
}

func (c *Client) SetAuth(auths string) {
	if len(auths) > 5 {
		c.auths = auths
	}
}

func (c *Client) SetContentType(ctype string) {
	if ctype != "" {
		c.ctype = ctype
	}
}

func (c *Client) makeRequest(method, uri string, body io.Reader) (req *http.Request, err error) {
	token, err := c.GetAuthToken()
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("access_token", token)
	u.RawQuery = q.Encode()
	uri = u.String()

	// logger().Debug("client.do", "method", method, "uri", uri)
	req, err = http.NewRequest(method, uri, body)
	if err != nil {
		logger().LogAttrs(context.Background(), slog.LevelInfo, "NewRequest fail",
			slog.Any("err", err),
			slog.String("method", method),
			slog.String("uri", uri),
		)
		return
	}
	if method == "POST" {
		if c.ctype != "" {
			req.Header.Set("Content-Type", c.ctype)
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}
	if c.auths != "" {
		req.Header.Set("Authorization", c.auths)
	}
	return
}

func (c *Client) Do(method, uri string, body io.Reader, rf respFunc) error {
	req, err := c.makeRequest(method, uri, body)
	if err != nil {
		return err
	}
	return doRequest(c.httpClient, req, rf)
}

type respFunc func(hdr http.Header, r io.Reader, cl int64) error

func doRequest(client *http.Client, req *http.Request, rf respFunc) error {
	resp, e := client.Do(req)
	if e != nil {
		logger().LogAttrs(context.Background(), slog.LevelInfo, "client.do fail",
			slog.String("method", req.Method),
			slog.String("uri", req.RequestURI),
			slog.Any("err", e),
		)
		return e
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		logger().LogAttrs(context.Background(), slog.LevelInfo, "http fail",
			slog.Int("code", resp.StatusCode),
			slog.String("status", resp.Status),
		)
		return fmt.Errorf("expecting HTTP status code 20x, but got %d", resp.StatusCode)
	}

	if rf != nil {
		// logger().Debug("doRequest ok", "status", resp.StatusCode, "length", resp.ContentLength, "uri", req.URL.Path)
		if err := rf(resp.Header, resp.Body, resp.ContentLength); err != nil {
			logger().LogAttrs(context.Background(), slog.LevelInfo, "call respFunc fail",
				slog.String("uri", req.URL.Path),
				slog.Any("err", err),
			)
			return err
		}
	}

	return nil
}

func doRequestData(client *http.Client, req *http.Request) (out []byte, err error) {
	err = doRequest(client, req, func(_ http.Header, r io.Reader, _ int64) error {
		rbody, e := io.ReadAll(r)
		if e != nil {
			logger().LogAttrs(context.Background(), slog.LevelInfo, "resp read fail",
				slog.Any("err", e),
			)
			return e
		}
		out = rbody
		return nil
	})
	return
}

func DoHTTP(method, uri string, auths string, body io.Reader, rf respFunc) error {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return err
	}
	if auths != "" {
		req.Header.Set("Authorization", auths)
	}
	hc := &http.Client{Transport: tr}
	return doRequest(hc, req, rf)
}

func DoHTTPData(method, uri string, auths string, body io.Reader) (out []byte, err error) {
	err = DoHTTP(method, uri, auths, body, func(_ http.Header, r io.Reader, _ int64) error {
		rbody, e := io.ReadAll(r)
		if e != nil {
			logger().LogAttrs(context.Background(), slog.LevelInfo, "resp read fail",
				slog.Any("err", e),
			)
			return e
		}
		out = rbody
		return nil
	})
	return
}

func (c *Client) Get(uri string) (out []byte, err error) {
	req, err := c.makeRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	return doRequestData(c.httpClient, req)
}

func (c *Client) Post(uri string, data []byte) (out []byte, err error) {
	req, err := c.makeRequest("POST", uri, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return doRequestData(c.httpClient, req)
}

func (c *Client) GetJSON(uri string, obj any) error {
	err := c.Do("GET", uri, nil, func(_ http.Header, r io.Reader, _ int64) error {
		return jsonInto(r, obj)
	})
	if err != nil {
		logger().LogAttrs(context.Background(), slog.LevelInfo, "GetJSON fail",
			slog.String("uri", uri),
			slog.Any("err", err),
		)
	}
	return err
}

func (c *Client) PostJSON(uri string, data []byte, obj any) error {
	err := c.Do("POST", uri, bytes.NewReader(data), func(_ http.Header, r io.Reader, _ int64) error {
		return jsonInto(r, obj)
	})
	if err != nil {
		logger().LogAttrs(context.Background(), slog.LevelInfo, "PostJSON fail",
			slog.String("uri", uri),
			slog.String("data", string(data)),
			slog.Any("err", err),
		)
	}
	return err
}

func (c *Client) PostObj(uri string, in any, res any) (err error) {
	if in == nil {
		return fmt.Errorf("empty input")
	}
	var data []byte
	data, err = json.Marshal(in)
	if err != nil {
		return
	}
	return c.PostJSON(uri, data, res)
}

func MustMarshal(obj any) []byte {
	data, err := json.Marshal(obj)
	if err != nil {
		logger().LogAttrs(context.Background(), slog.LevelError, "marshal fail",
			slog.Any("obj", obj),
			slog.Any("err", err),
		)
		os.Exit(1)
	}
	return data
}

func jsonInto(r io.Reader, obj any) error {
	err := json.NewDecoder(r).Decode(obj)
	if err != nil {
		logger().LogAttrs(context.Background(), slog.LevelInfo, "resp decode fail",
			slog.Any("err", err),
		)
		return err
	}
	if ce, ok := obj.(ErrorCoder); ok {
		if code := ce.GetErrorCode(); code != 0 {
			logger().LogAttrs(context.Background(), slog.LevelInfo, "resp has error",
				slog.Int("code", code),
				slog.String("msg", ce.GetErrorMsg()),
			)
			return ce.GetErr()
		}

		logger().LogAttrs(context.Background(), slog.LevelDebug, "resp decode done",
			slog.String("msg", ce.GetErrorMsg()),
		)
	} else {
		logger().LogAttrs(context.Background(), slog.LevelDebug, "resp decode done",
			slog.Any("obj", obj),
		)
	}

	return err
}
