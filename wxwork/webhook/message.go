package webhook

type Message struct {
	// 消息类型 text/markdown/image/news
	MsgType string `json:"msgtype"`

	// 内容 text/markdown
	Content string `json:"content,omitempty"`

	// userid的列表
	MentionedList []string `json:"mentioned_list,omitempty"`

	// 手机号列表
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`

	// Image
	Image *MessageImage `json:"image,omitempty"`

	// Articles
	Articles []MessageArticle `json:"articles"` // 图文消息，一个图文消息支持1到8条图文
}

// MessageImage 图片
type MessageImage struct {
	Base64 string `json:"base64"` // 图片内容的base64编码
	MD5    string `json:"md5"`    // 图片内容（base64编码前）的md5值
}

// Article 图文
type MessageArticle struct {
	Title       string `json:"title"`                 // 标题，不超过128个字节，超过会自动截断
	Description string `json:"description,omitempty"` // 描述，不超过512个字节，超过会自动截断
	URL         string `json:"url"`                   // 点击后跳转的链接
	PicURL      string `json:"picurl,omitempty"`      // 图文消息的图片链接，支持JPG、PNG格式，较好的效果为大图 1068*455，小图150*150
}

func NewTextMessage(content string) *Message {
	return &Message{
		MsgType: "text",
		Content: content,
	}
}

func NewMarkdownMessage(content string) *Message {
	return &Message{
		MsgType: "markdown",
		Content: content,
	}
}

// TODO: more msg type
