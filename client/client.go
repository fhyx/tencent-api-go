package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

func (c *Client) Do(method, uri string, body io.Reader) ([]byte, error) {
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

	logger().Debugw("client.do", "method", method, "uri", uri)
	req, e := http.NewRequest(method, uri, body)
	if e != nil {
		log.Println(e, method, uri)
		return nil, e
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

	return doRequest(c.httpClient, req)
}

func doRequest(client *http.Client, req *http.Request) ([]byte, error) {
	resp, e := client.Do(req)
	if e != nil {
		logger().Infow("client.do fail", "method", req.Method, "uri", req.RequestURI, "err", e)
		return nil, e
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		logger().Infow("http fail", "code", resp.StatusCode, "status", resp.Status)
		return nil, fmt.Errorf("Expecting HTTP status code 20x, but got %v", resp.StatusCode)
	}

	rbody, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		logger().Infow("resp read fail", "err", e)
		return nil, e
	}
	logger().Debugw("resp read", "body", string(rbody))

	return rbody, nil

}

func DoHTTP(method, uri string, auths string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}
	if auths != "" {
		req.Header.Set("Authorization", auths)
	}
	hc := &http.Client{Transport: tr}
	return doRequest(hc, req)
}

func (c *Client) Get(uri string) ([]byte, error) {
	return c.Do("GET", uri, nil)
}

func (c *Client) Post(uri string, data []byte) ([]byte, error) {
	return c.Do("POST", uri, bytes.NewReader(data))
}

func (c *Client) GetJSON(uri string, obj interface{}) error {
	body, err := c.Get(uri)
	if err != nil {
		return err
	}
	err = parseResult(body, obj)
	if err != nil {
		log.Printf("GetJSON(uri %s) ERR %s", uri, err)
	}
	return err
}

func (c *Client) PostJSON(uri string, data []byte, obj interface{}) error {
	body, err := c.Post(uri, data)
	if err != nil {
		return err
	}
	err = parseResult(body, obj)
	if err != nil {
		logger().Infow("PostJSON fail", "uri", uri, "data", len(data), "err", err)
	}
	return err
}

func parseResult(resp []byte, obj interface{}) error {
	// log.Printf("parse result: %s", string(resp))
	exErr := &Error{}
	if e := json.Unmarshal(resp, exErr); e != nil {
		log.Printf("unmarshal api err %s", e)
		return e
	}

	if exErr.ErrCode != 0 {
		log.Printf("apiError %s", exErr)
		return exErr
	}

	if e := json.Unmarshal(resp, obj); e != nil {
		log.Printf("unmarshal obj err %s", e)
		return e
	}

	return nil
}
