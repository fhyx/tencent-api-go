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

// ==================== 模板卡片消息 ====================

// TemplateCardSource 卡片来源样式
type TemplateCardSource struct {
	IconURL   string `json:"icon_url,omitempty"`    // 来源图片url
	Desc      string `json:"desc,omitempty"`        // 来源图片描述，建议不超过20个字
	DescColor int    `json:"desc_color,omitempty"` // 来源文字颜色: 0灰色(默认) 1黑色 2红色 3绿色
}

// TemplateCardActionMenu 右上角菜单
type TemplateCardActionMenu struct {
	Desc       string                        `json:"desc,omitempty"`        // 卡片副交互辅助文本说明
	ActionList []TemplateCardActionMenuItem `json:"action_list"`           // 操作列表，列表长度1~3
}

// TemplateCardActionMenuItem 右上角菜单项
type TemplateCardActionMenuItem struct {
	Text string `json:"text"` // 操作描述文案
	Key  string `json:"key"`  // 操作key，用户点击后回调时作为EventKey返回
}

// TemplateCardMainTitle 一级标题
type TemplateCardMainTitle struct {
	Title string `json:"title,omitempty"` // 一级标题，建议不超过36个字
	Desc  string `json:"desc,omitempty"`  // 标题辅助信息，建议不超过44个字
}

// TemplateCardQuoteArea 引用文献样式
type TemplateCardQuoteArea struct {
	Type      int    `json:"type,omitempty"`       // 点击事件类型: 0无 1跳转URL 2跳转小程序
	URL       string `json:"url,omitempty"`          // 点击跳转URL，type=1时必填
	AppID     string `json:"appid,omitempty"`        // 小程序appid，type=2时必填
	PagePath  string `json:"pagepath,omitempty"`      // 小程序页面路径
	Title     string `json:"title,omitempty"`        // 引用文献样式标题
	QuoteText string `json:"quote_text,omitempty"`   // 引用文献样式的引用文案
}

// TemplateCardHorizontalContent 二级标题+文本
type TemplateCardHorizontalContent struct {
	Type    int    `json:"type,omitempty"`     // 链接类型: 0非链接 1跳转URL 2下载附件 3成员详情
	KeyName string `json:"keyname"`             // 二级标题，建议不超过5个字
	Value   string `json:"value,omitempty"`     // 二级文本
	URL     string `json:"url,omitempty"`       // 跳转URL，type=1时必填
	MediaID string `json:"media_id,omitempty"`  // 附件media_id，type=2时必填
	UserID  string `json:"userid,omitempty"`    // 成员userid，type=3时必填
}

// TemplateCardJump 跳转指引
type TemplateCardJump struct {
	Type    int    `json:"type,omitempty"`     // 跳转类型: 0非链接 1跳转URL 2跳转小程序
	Title   string `json:"title"`               // 跳转链接文案，建议不超过18个字
	URL     string `json:"url,omitempty"`       // 跳转URL，type=1时必填
	AppID   string `json:"appid,omitempty"`     // 小程序appid，type=2时必填
	PagePath string `json:"pagepath,omitempty"`  // 小程序页面路径
}

// TemplateCardAction 卡片整体点击事件
type TemplateCardAction struct {
	Type     int    `json:"type,omitempty"`      // 跳转类型: 0非链接 1跳转URL 2打开小程序
	URL      string `json:"url,omitempty"`       // 跳转URL，type=1时必填
	AppID    string `json:"appid,omitempty"`      // 小程序appid，type=2时必填
	PagePath string `json:"pagepath,omitempty"`   // 小程序页面路径
}

// TemplateCardButtonSelection 下拉选择器
type TemplateCardButtonSelection struct {
	QuestionKey string                    `json:"question_key"`           // 选择器key，用户提交后回调时返回
	Title       string                    `json:"title,omitempty"`        // 选择器左边标题
	OptionList  []TemplateCardSelectOption `json:"option_list"`            // 选项列表，最多10个
	SelectedID  string                    `json:"selected_id,omitempty"`  // 默认选定id
}

// TemplateCardSelectOption 下拉选择器选项
type TemplateCardSelectOption struct {
	ID   string `json:"id"`   // 选项id，用户提交后回调时返回
	Text string `json:"text"` // 选项文案，建议不超过16个字
}

// TemplateCardButton 按钮
type TemplateCardButton struct {
	Type  int    `json:"type,omitempty"`    // 点击事件类型: 0回调事件 1跳转URL，默认0
	Text  string `json:"text"`              // 按钮文案，建议不超过10个字
	Style int    `json:"style,omitempty"`   // 按钮样式 1~4
	Key   string `json:"key,omitempty"`     // 按钮key，type=0时必填
	URL   string `json:"url,omitempty"`     // 跳转URL，type=1时必填
}

// TemplateCardContent 模板卡片消息内容
type TemplateCardContent struct {
	CardType               string                        `json:"card_type,omitempty"`                // 卡片类型
	Source                 *TemplateCardSource            `json:"source,omitempty"`                   // 卡片来源样式
	ActionMenu             *TemplateCardActionMenu        `json:"action_menu,omitempty"`              // 右上角菜单
	MainTitle              TemplateCardMainTitle          `json:"main_title"`                         // 一级标题
	QuoteArea              *TemplateCardQuoteArea         `json:"quote_area,omitempty"`               // 引用文献样式
	SubTitleText           string                        `json:"sub_title_text,omitempty"`            // 二级普通文本，建议不超过160个字
	HorizontalContentList  []TemplateCardHorizontalContent `json:"horizontal_content_list,omitempty"`  // 二级标题+文本列表，最多6个
	JumpList               []TemplateCardJump             `json:"jump_list,omitempty"`                 // 跳转指引列表，最多3个
	CardAction             *TemplateCardAction            `json:"card_action,omitempty"`               // 整体卡片点击事件
	TaskID                 string                        `json:"task_id,omitempty"`                  // 任务id，用于后续更新卡片接口调用
	ButtonSelection        *TemplateCardButtonSelection   `json:"button_selection,omitempty"`         // 下拉选择器
	ButtonList             []TemplateCardButton           `json:"button_list,omitempty"`              // 按钮列表，最多6个
}

// SendTemplateCardMessage 发送模板卡片消息
//
//	收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
//	否则为单纯的【发送应用消息】接口调用。
//
//	cardType 可选值: text_notice, news_notice, button_interaction, vote_interaction, multiple_interaction
//	button_interaction/vote_interaction/multiple_interaction 类型必须提供 taskID
//
//	参考: https://developer.work.weixin.qq.com/document/path/90236#模板卡片消息
func (a *API) SendTemplateCardMessage(recipient *Recipient, cardType string, content *TemplateCardContent, isSafe bool) error {
	content.CardType = cardType
	// 转换为 map[string]any 以适配 sendMessage
	contentMap, err := structToMap(content)
	if err != nil {
		return err
	}
	return a.sendMessage(recipient, "template_card", contentMap, isSafe)
}

// SendTemplateCardButtonInteraction 发送按钮交互型模板卡片消息
//
//	收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
//	否则为单纯的【发送应用消息】接口调用。
//
//	按钮交互型卡片支持用户点击按钮触发回调事件，需配合接收 template_card_event 使用。
//	taskID 用于后续更新卡片接口调用，72小时内有效，且只能使用一次。
//
//	参考: https://developer.work.weixin.qq.com/document/path/90236#按钮交互型
func (a *API) SendTemplateCardButtonInteraction(
	recipient *Recipient,
	title string,
	desc string,
	buttons []TemplateCardButton,
	taskID string,
	isSafe bool,
) error {
	content := &TemplateCardContent{
		CardType: "button_interaction",
		MainTitle: TemplateCardMainTitle{
			Title: title,
			Desc:  desc,
		},
		TaskID:     taskID,
		ButtonList: buttons,
	}
	contentMap, err := structToMap(content)
	if err != nil {
		return err
	}
	return a.sendMessage(recipient, "template_card", contentMap, isSafe)
}

// structToMap converts a struct to map[string]any using JSON marshal/unmarshal
func structToMap(v any) (map[string]any, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}
