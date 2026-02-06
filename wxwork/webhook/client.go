package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	DefaultUrlTemplate = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s"
)

var hclient = &http.Client{
	Timeout: time.Second * 15,
}

type Notifier interface {
	Notify(msg *Message) error
	NotifyMarkdown(text string, mentions ...string) error
	NotifyText(text string, mentions ...string) error
}

type client struct {
	uri string
}

func NewClient(uri string) Notifier {
	return &client{uri}
}

func NewWithKey(key string) Notifier {
	return NewClient(fmt.Sprintf(DefaultUrlTemplate, key))
}

func (c *client) NotifyMarkdown(text string, mentions ...string) error {
	msg := NewMarkdownMessage(text)
	if len(mentions) > 0 {
		msg.MentionedList = mentions
	}
	return c.Notify(msg)
}

func (c *client) NotifyText(text string, mentions ...string) error {
	msg := NewTextMessage(text)
	if len(mentions) > 0 {
		msg.MentionedList = mentions
	}
	return c.Notify(msg)
}

func (c *client) Notify(msg *Message) error {
	if len(c.uri) == 0 {
		logger().Infow("empty uri, notify to log", "msg", msg)
		return nil
	}
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
	if resp.StatusCode > 300 {
		logger().Infow("notified", "msg", msg, "status", resp.Status)
	} else {
		logger().Debugw("notified", "msg", msg, "status", resp.StatusCode)
	}
	_ = resp.Body.Close()

	return nil
}
