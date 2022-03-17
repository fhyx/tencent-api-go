package receiv

import (
	"encoding/xml"
	"fmt"
)

type Event struct {
	XMLName      xml.Name    `xml:"xml"`
	ToUserName   string      `xml:"ToUserName"`   // 企业微信CorpID
	FromUserName string      `xml:"FromUserName"` // 此事件该值固定为sys，表示该消息由系统生成
	CreateTime   int64       `xml:"CreateTime"`   // 消息创建时间 （整型）
	MsgType      MessageType `xml:"MsgType"`      // 消息的类型，此时固定为event
	EvnType      EventType   `xml:"Event"`
}

func (e *Event) String() string {
	return fmt.Sprintf("%s %s", e.MsgType, e.EvnType)
}

type EventChangeContact struct {
	Event
	ChangeType ChangeType `xml:"ChangeType"`
}

func (e *EventChangeContact) String() string {
	return fmt.Sprintf("%s %s %s", e.MsgType, e.EvnType, e.ChangeType)
}

type eventChangeContactUser struct {
	Name           string `xml:"Name"`
	Department     string `xml:"Department"`     // 1,2,3
	MainDepartment int32  `xml:"MainDepartment"` // 主部门
	IsLeaderInDept string `xml:"IsLeaderInDept"` // 1,0,0 是否为部门负责人，0-否，1-是，顺序与Department字段的部门逐一对
	DirectLeader   string `xml:"DirectLeader"`   // 直属上级UserID，最多5个，逗号分隔
	Position       string `xml:"Position"`
	Mobile         string `xml:"Mobile"`
	Gender         int8   `xml:"Gender"`  // 性别，1表示男性，2表示女性
	Email          string `xml:"Email"`   // 邮箱
	BizMail        string `xml:"BizMail"` // 企业邮箱
	Status         int32  `xml:"Status"`  // 激活状态：1=已激活 2=已禁用 4=未激活 已激活代表已激活企业微信或已关注微信插件（原企业号）5=成员退出
	Avatar         string `xml:"Avatar"`  // 头像url。注：如果要获取小图将url最后的”/0”改成”/100”即可。
	Alias          string `xml:"Alias"`
	Telephone      string `xml:"Telephone"`
	Address        string `xml:"Address"`
}

type eventChangeContactExtAttr struct {
	Text  string `xml:",chardata" json:"-"`
	Attrs []struct {
		Chardata string `xml:",chardata" json:"-"`
		Name     string `xml:"Name" json:"name,omitempty"`
		Type     string `xml:"Type" json:"type,omitempty"`
		Text     struct {
			Text  string `xml:",chardata" json:"-"`
			Value string `xml:"Value" json:"value,omitempty"`
		} `xml:"Text" json:"text"`
		Web struct {
			Text  string `xml:",chardata" json:"-"`
			Title string `xml:"Title" json:"title,omitempty"`
			URL   string `xml:"Url" json:"url,omitempty"`
		} `xml:"Web" json:"web"`
	} `xml:"Item" json:"attrs,omitempty"`
}

/*
<xml>
	<ToUserName><![CDATA[toUser]]></ToUserName>
	<FromUserName><![CDATA[sys]]></FromUserName>
	<CreateTime>1403610513</CreateTime>
	<MsgType><![CDATA[event]]></MsgType>
	<Event><![CDATA[change_contact]]></Event>
	<ChangeType>create_user</ChangeType>
	<UserID><![CDATA[zhangsan]]></UserID>
	<Name><![CDATA[张三]]></Name>
	<Department><![CDATA[1,2,3]]></Department>
	<MainDepartment>1</MainDepartment>
	<IsLeaderInDept><![CDATA[1,0,0]]></IsLeaderInDept>
	<DirectLeader><![CDATA[lisi,wangwu]]></DirectLeader>
	<Position><![CDATA[产品经理]]></Position>
	<Mobile>13800000000</Mobile>
	<Gender>1</Gender>
	<Email><![CDATA[zhangsan@gzdev.com]]></Email>
	<BizMail><![CDATA[zhangsan@qyycs2.wecom.work]]></BizMail>
	<Status>1</Status>
	<Avatar><![CDATA[http://wx.qlogo.cn/mmopen/ajNVdqHZLLA3WJ6DSZUfiakYe37PKnQhBIeOQBO4czqrnZDS79FH5Wm5m4X69TBicnHFlhiafvDwklOpZeXYQQ2icg/0]]></Avatar>
	<Alias><![CDATA[zhangsan]]></Alias>
	<Telephone><![CDATA[020-123456]]></Telephone>
	<Address><![CDATA[广州市]]></Address>
	<ExtAttr>
		<Item>
		<Name><![CDATA[爱好]]></Name>
		<Type>0</Type>
		<Text>
			<Value><![CDATA[旅游]]></Value>
		</Text>
		</Item>
		<Item>
		<Name><![CDATA[卡号]]></Name>
		<Type>1</Type>
		<Web>
			<Title><![CDATA[企业微信]]></Title>
			<Url><![CDATA[https://work.weixin.qq.com]]></Url>
		</Web>
		</Item>
	</ExtAttr>
</xml>
*/
type EventChangeContactCreateUser struct {
	EventChangeContact

	UserID string `xml:"UserID"` // 成员UserID

	eventChangeContactUser

	ExtAttr eventChangeContactExtAttr `xml:"ExtAttr" json:"extAttr"`
}

func (e *EventChangeContactCreateUser) GetID() string {
	return e.UserID
}

func (e *EventChangeContactCreateUser) GetName() string {
	return e.Name
}

/*
<xml>
	<ToUserName><![CDATA[toUser]]></ToUserName>
	<FromUserName><![CDATA[sys]]></FromUserName>
	<CreateTime>1403610513</CreateTime>
	<MsgType><![CDATA[event]]></MsgType>
	<Event><![CDATA[change_contact]]></Event>
	<ChangeType>update_user</ChangeType>
	<UserID><![CDATA[zhangsan]]></UserID>
	<NewUserID><![CDATA[zhangsan001]]></NewUserID>
	<Name><![CDATA[张三]]></Name>
	<Department><![CDATA[1,2,3]]></Department>
	<MainDepartment>1</MainDepartment>
	<IsLeaderInDept><![CDATA[1,0,0]]></IsLeaderInDept>
	<Position><![CDATA[产品经理]]></Position>
	<Mobile>13800000000</Mobile>
	<Gender>1</Gender>
	<Email><![CDATA[zhangsan@gzdev.com]]></Email>
	<Status>1</Status>
	<Avatar><![CDATA[http://wx.qlogo.cn/mmopen/ajNVdqHZLLA3WJ6DSZUfiakYe37PKnQhBIeOQBO4czqrnZDS79FH5Wm5m4X69TBicnHFlhiafvDwklOpZeXYQQ2icg/0]]></Avatar>
	<Alias><![CDATA[zhangsan]]></Alias>
	<Telephone><![CDATA[020-123456]]></Telephone>
	<Address><![CDATA[广州市]]></Address>
	<ExtAttr>
		<Item>
		<Name><![CDATA[爱好]]></Name>
		<Type>0</Type>
		<Text>
			<Value><![CDATA[旅游]]></Value>
		</Text>
		</Item>
		<Item>
		<Name><![CDATA[卡号]]></Name>
		<Type>1</Type>
		<Web>
			<Title><![CDATA[企业微信]]></Title>
			<Url><![CDATA[https://work.weixin.qq.com]]></Url>
		</Web>
		</Item>
	</ExtAttr>
</xml>
*/
type EventChangeContactUpdateUser struct {
	EventChangeContact

	UserID    string `xml:"UserID"`    // 变更信息的成员UserID
	NewUserID string `xml:"NewUserID"` // 新的UserID，变更时推送（userid由系统生成时可更改一次）

	eventChangeContactUser

	ExtAttr eventChangeContactExtAttr `xml:"ExtAttr"`
}

func (e *EventChangeContactUpdateUser) GetID() string {
	if len(e.NewUserID) > 0 {
		return e.UserID + "|" + e.NewUserID
	}
	return e.UserID
}

func (e *EventChangeContactUpdateUser) GetName() string {
	return e.Name
}

/*
<xml>
	<ToUserName><![CDATA[toUser]]></ToUserName>
	<FromUserName><![CDATA[sys]]></FromUserName>
	<CreateTime>1403610513</CreateTime>
	<MsgType><![CDATA[event]]></MsgType>
	<Event><![CDATA[change_contact]]></Event>
	<ChangeType>delete_user</ChangeType>
	<UserID><![CDATA[zhangsan]]></UserID>
</xml>
*/
type EventChangeContactDeleteUser struct {
	EventChangeContact
	UserID string `xml:"UserID"`
}

func (e *EventChangeContactDeleteUser) GetID() string {
	return e.UserID
}

/*
<xml>
    <ToUserName><![CDATA[toUser]]></ToUserName>
    <FromUserName><![CDATA[sys]]></FromUserName>
    <CreateTime>1403610513</CreateTime>
    <MsgType><![CDATA[event]]></MsgType>
    <Event><![CDATA[change_contact]]></Event>
    <ChangeType>create_party</ChangeType>
    <Id>2</Id>
    <Name><![CDATA[张三]]></Name>
    <ParentId><![CDATA[1]]></ParentId>
    <Order>1</Order>
</xml>
*/
type EventChangeContactCreateParty struct {
	EventChangeContact
	ID       string `xml:"Id"`
	Name     string `xml:"Name"`
	ParentId string `xml:"ParentId"`
	Order    int32  `xml:"Order"`
}

func (e *EventChangeContactCreateParty) GetID() string {
	return e.ID
}

func (e *EventChangeContactCreateParty) GetName() string {
	return e.Name
}

/*
<xml>
    <ToUserName><![CDATA[toUser]]></ToUserName>
    <FromUserName><![CDATA[sys]]></FromUserName>
    <CreateTime>1403610513</CreateTime>
    <MsgType><![CDATA[event]]></MsgType>
    <Event><![CDATA[change_contact]]></Event>
    <ChangeType>update_party</ChangeType>
    <Id>2</Id>
    <Name><![CDATA[张三]]></Name>
    <ParentId><![CDATA[1]]></ParentId>
</xml>
*/
type EventChangeContactUpdateParty struct {
	EventChangeContact
	ID       string `xml:"Id"`
	Name     string `xml:"Name"`
	ParentId string `xml:"ParentId"`
	// Order    int32  `xml:"Order"`
}

func (e *EventChangeContactUpdateParty) GetID() string {
	return e.ID
}

func (e *EventChangeContactUpdateParty) GetName() string {
	return e.Name
}

/*
<xml>
    <ToUserName><![CDATA[toUser]]></ToUserName>
    <FromUserName><![CDATA[sys]]></FromUserName>
    <CreateTime>1403610513</CreateTime>
    <MsgType><![CDATA[event]]></MsgType>
    <Event><![CDATA[change_contact]]></Event>
    <ChangeType>delete_party</ChangeType>
    <Id>2</Id>
</xml>
*/
type EventChangeContactDeleteParty struct {
	EventChangeContact
	ID string `xml:"Id"`
}

func (e *EventChangeContactDeleteParty) GetID() string {
	return e.ID
}

/*
<xml>
    <ToUserName><![CDATA[toUser]]></ToUserName>
    <FromUserName><![CDATA[sys]]></FromUserName>
    <CreateTime>1403610513</CreateTime>
    <MsgType><![CDATA[event]]></MsgType>
    <Event><![CDATA[change_contact]]></Event>
    <ChangeType><![CDATA[update_tag]]></ChangeType>
    <TagId>1</TagId>
    <AddUserItems><![CDATA[zhangsan,lisi]]></AddUserItems>
    <DelUserItems><![CDATA[zhangsan1,lisi1]]></DelUserItems>
    <AddPartyItems><![CDATA[1,2]]></AddPartyItems>
    <DelPartyItems><![CDATA[3,4]]></DelPartyItems>
</xml>
*/
type EventChangeContactUpdateTag struct {
	EventChangeContact
	TagId         string `xml:"TagId"`
	AddUserItems  string `xml:"AddUserItems"`
	DelUserItems  string `xml:"DelUserItems"`
	AddPartyItems string `xml:"AddPartyItems"`
	DelPartyItems string `xml:"DelPartyItems"`
}
