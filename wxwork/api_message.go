package wxwork

import (
	"encoding/json"
	"errors"
	"strings"

	"daxv.cn/gopak/tencent-api-go/client"
)

// reqMessage 消息发送请求
type reqMessage struct {
	ToUser  []string
	ToParty []string
	ToTag   []string
	ChatID  string
	AgentID int
	MsgType string
	Content map[string]any
	IsSafe  bool
}

func (x reqMessage) MarshalJSON() ([]byte, error) {
	safeInt := 0
	if x.IsSafe {
		safeInt = 1
	}

	out := map[string]any{
		"msgtype": x.MsgType,
		"agentid": x.AgentID,
		"safe":    safeInt,
	}

	// msgtype polymorphism
	out[x.MsgType] = x.Content

	// 复用这个结构体，因为是 package-private 的所以这么做没风险
	if x.ChatID != "" {
		out["chatid"] = x.ChatID
	} else {
		out["touser"] = strings.Join(x.ToUser, "|")
		out["toparty"] = strings.Join(x.ToParty, "|")
		out["totag"] = strings.Join(x.ToTag, "|")
	}

	result, err := json.Marshal(out)
	if err != nil {
		// should never happen unless OOM or similar bad things
		// TODO: error_chain
		return nil, err
	}

	return result, nil
}

// respMessageSend 消息发送响应
type respMessageSend struct {
	client.Error

	InvalidUsers   string `json:"invaliduser"`
	InvalidParties string `json:"invalidparty"`
	InvalidTags    string `json:"invalidtag"`
}

// execMessageSend 发送应用消息
func (a *API) execMessageSend(req *reqMessage) (*respMessageSend, error) {
	var resp respMessageSend
	err := a.c.PostJSON(UriPrefix+"/message/send", client.MustMarshal(req), &resp)
	if err != nil {
		logger().Infow("execMessageSend fail", "req", req, "err", err)
		return nil, err
	}

	return &resp, nil
}

// execAppchatSend 应用推送消息
func (a *API) execAppchatSend(req *reqMessage) (*respMessageSend, error) {
	var resp respMessageSend
	err := a.c.PostJSON(UriPrefix+"/appchat/send", client.MustMarshal(req), &resp)
	if err != nil {
		logger().Infow("execAppchatSend fail", "req", req, "err", err)
		return nil, err
	}

	return &resp, nil
}

// sendMessage 发送消息底层接口
//
// 收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
// 否则为单纯的【发送应用消息】接口调用。
func (a *API) sendMessage(recipient *Recipient, msgtype string, content map[string]any, isSafe bool) error {
	isApichatSendRequest := false
	if !recipient.isValidForMessageSend() {
		if !recipient.isValidForAppchatSend() {
			// TODO: better error
			return errors.New("recipient invalid for message sending")
		}

		// 发送给群聊
		isApichatSendRequest = true
	}

	req := &reqMessage{
		ToUser:  recipient.UserIDs,
		ToParty: recipient.PartyIDs,
		ToTag:   recipient.TagIDs,
		ChatID:  recipient.ChatID,
		AgentID: a.AgentID,
		MsgType: msgtype,
		Content: content,
		IsSafe:  isSafe,
	}

	var resp *respMessageSend
	var err error
	if isApichatSendRequest {
		resp, err = a.execAppchatSend(req)
	} else {
		resp, err = a.execMessageSend(req)
	}

	if err != nil {
		logger().Infow("sendMessage fail", "err", err)
		return err
	}

	// TODO: what to do with resp?
	_ = resp
	return nil
}

// SendTextMessage 发送文本消息
//
// 收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
// 否则为单纯的【发送应用消息】接口调用。
func (a *API) SendTextMessage(recipient *Recipient, content string, isSafe bool) error {
	return a.sendMessage(recipient, "text", map[string]any{"content": content}, isSafe)
}

// SendImageMessage 发送图片消息
//
// 收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
// 否则为单纯的【发送应用消息】接口调用。
func (a *API) SendImageMessage(recipient *Recipient, mediaID string, isSafe bool) error {
	return a.sendMessage(recipient, "image", map[string]any{
		"media_id": mediaID,
	}, isSafe)
}

// SendTextCardMessage 发送文本卡片消息
//
// 收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
// 否则为单纯的【发送应用消息】接口调用。
func (a *API) SendTextCardMessage(
	recipient *Recipient,
	title string,
	description string,
	url string,
	buttonText string,
	isSafe bool,
) error {
	return a.sendMessage(
		recipient,
		"textcard",
		map[string]any{
			"title":       title,
			"description": description,
			"url":         url,
			"btntxt":      buttonText, // TODO: 零值
		}, isSafe,
	)
}

// SendMarkdownMessage 发送 Markdown 消息
//
// 仅支持 Markdown 的子集，详见[官方文档](https://work.weixin.qq.com/api/doc#90002/90151/90854/%E6%94%AF%E6%8C%81%E7%9A%84markdown%E8%AF%AD%E6%B3%95)。
//
// 收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
// 否则为单纯的【发送应用消息】接口调用。
func (a *API) SendMarkdownMessage(
	recipient *Recipient,
	content string,
	isSafe bool,
) error {
	return a.sendMessage(recipient, "markdown", map[string]interface{}{"content": content}, isSafe)
}
