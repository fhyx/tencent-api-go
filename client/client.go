package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

	// logger().Debugw("client.do", "method", method, "uri", uri)
	req, err = http.NewRequest(method, uri, body)
	if err != nil {
		logger().Infow("NewRequest fail", "err", err, "method", method, "uri", uri)
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
		logger().Infow("client.do fail", "method", req.Method, "uri", req.RequestURI, "err", e)
		return e
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		logger().Infow("http fail", "code", resp.StatusCode, "status", resp.Status)
		return fmt.Errorf("Expecting HTTP status code 20x, but got %v", resp.StatusCode)
	}

	if rf != nil {
		// logger().Debugw("doRequest ok", "status", resp.StatusCode, "length", resp.ContentLength, "uri", req.URL.Path)
		if err := rf(resp.Header, resp.Body, resp.ContentLength); err != nil {
			logger().Infow("call respFunc fail", "uri", req.URL.Path, "err", err)
			return err
		}
	}

	return nil
}

func doRequestData(client *http.Client, req *http.Request) (out []byte, err error) {
	err = doRequest(client, req, func(_ http.Header, r io.Reader, _ int64) error {
		rbody, e := io.ReadAll(r)
		if e != nil {
			logger().Infow("resp read fail", "err", e)
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
			logger().Infow("resp read fail", "err", e)
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
		logger().Infow("GetJSON fail", "uri", uri, "err", err)
	}
	return err
}

func (c *Client) PostJSON(uri string, data []byte, obj any) error {
	err := c.Do("POST", uri, bytes.NewReader(data), func(_ http.Header, r io.Reader, _ int64) error {
		return jsonInto(r, obj)
	})
	if err != nil {
		logger().Infow("PostJSON fail", "uri", uri, "data", len(data), "err", err)
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
		logger().Fatalw("marshal fail", "obj", obj, "err", err)
	}
	return data
}

func jsonInto(r io.Reader, obj any) error {
	err := json.NewDecoder(r).Decode(obj)
	if err != nil {
		logger().Infow("resp decode fail", "err", err)
		return err
	}
	if ce, ok := obj.(ErrorCoder); ok {
		if code := ce.GetErrorCode(); code != 0 {
			err = ce.GetErr()
			logger().Infow("resp has error", "code", code, "err", err)
		} else {
			logger().Debugw("resp decode done", "ce", ce.GetErrorMsg())
		}
	} else {
		logger().Debugw("resp decode done", "obj", obj)
	}

	return err
}
