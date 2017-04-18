package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	auths      string
	ctype      string
}

func NewClient() *Client {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	return &Client{
		httpClient: &http.Client{Transport: tr},
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

func (c *Client) Do(method, url string, body io.Reader) ([]byte, error) {
	req, e := http.NewRequest(method, url, body)
	if e != nil {
		log.Println(e, method, url)
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

	resp, e := c.httpClient.Do(req)
	if e != nil {
		log.Printf("client %s %s ERR %s", method, url, e)
		return nil, e
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		log.Printf("http code error %d, %s", resp.StatusCode, resp.Status)
		return nil, fmt.Errorf("Expecting HTTP status code 20x, but got %v", resp.StatusCode)
	}

	rbody, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		log.Printf("read body ERR %s", e)
		return nil, e
	}

	return rbody, nil
}

func DoHTTP(method, url string, auths string, body io.Reader) ([]byte, error) {
	c := NewClient()
	c.SetAuth(auths)
	return c.Do(method, url, body)
}

func (c *Client) Get(url string) ([]byte, error) {
	return c.Do("GET", url, nil)
}

func (c *Client) Post(url string, data []byte) ([]byte, error) {
	return c.Do("POST", url, bytes.NewReader(data))
}
