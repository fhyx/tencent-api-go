package receiv

// EventType 事件类型
type EventType string

// MessageType 消息类型
type MessageType string

// ChangeType 变更类型
type ChangeType string

// Message 接收到的消息
type Message struct {
	// ToUserName 企业微信 CorpID
	ToUserName string `xml:"ToUserName"`
	// FromUserName 成员 UserID
	FromUserName string `xml:"FromUserName"`
	// CreateTime 消息创建时间（整型）
	CreateTime int64 `xml:"CreateTime"`
	// MsgType 消息类型
	MsgType MessageType `xml:"MsgType"`
	// MsgID 消息 id，64 位整型
	MsgID int64 `xml:"MsgId"`
	// AgentID 企业应用的 id，整型。可在应用的设置页面查看
	AgentID int64 `xml:"AgentID"`
	// Event 事件类型 MsgType 为 event
	EvnType EventType `xml:"Event"`
	// ChangeType 变更类型 Event 为 change_contact 等的值：
	// create_user,update_user,delete_user,
	// create_party,update_party,delete_party,
	ChangeType ChangeType `xml:"ChangeType"`
}

type MessageText struct {
	Message
	Content string
}

type MessageImage struct {
	Message
	PicUrl  string
	MediaId string
}

type MessageVoice struct {
	Message
	Format      string
	Recognition string
}

type MessageVideo struct {
	Message
	MediaId      string
	ThumbMediaId string
}

type MessageLocation struct {
	Message
	Location_X string
	Location_Y string
	Scale      string
	Label      string
}

type MessageLink struct {
	Message
	Title       string
	Description string
	Url         string
}

type IDGetter interface {
	GetID() string
}

type NameGetter interface {
	GetName() string
}

type MessageGetter interface {
	GetMessage() string
}
