package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func DoHTTP(method, url string, auths string, body io.Reader) ([]byte, error) {
	req, e := http.NewRequest(method, url, body)
	if e != nil {
		log.Println(e, method, url)
		return nil, e
	}

	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if auths != "" {
		req.Header.Set("Authorization", auths)
	}

	c := &http.Client{}
	resp, e := c.Do(req)
	if e != nil {
		return nil, e
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		log.Printf("http code error %d, %s", resp.StatusCode, resp.Status)
		return nil, fmt.Errorf("Expecting HTTP status code 20x, but got %v", resp.StatusCode)
	}

	rbody, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}

	return rbody, nil
}
