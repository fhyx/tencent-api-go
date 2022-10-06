package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

var hclient = &http.Client{
	Timeout: time.Second * 15,
}

type Notifier interface {
	Notify(msg *Message) error
}

type client struct {
	uri string
}

func NewClient(uri string) Notifier {
	return &client{uri}
}

func (c *client) Notify(msg *Message) error {
	b, err := json.Marshal(msg)
	if err != nil {
		logger().Infow("marshal fail", "err", err)
		return err
	}
	body := bytes.NewReader(b)
	resp, err := hclient.Post(c.uri, "application/json; charset=UTF-8", body)
	if err != nil {
		logger().Infow("notify fail", "err", err, "msg", msg)
		return err
	}
	resp.Body.Close()
	logger().Infow("notify ok", "msg", msg, "status", resp.StatusCode)
	return nil
}
