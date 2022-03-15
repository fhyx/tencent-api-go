package receiv

import (
	"encoding/xml"
)

type Event struct {
	XMLName      xml.Name    `xml:"xml"`
	ToUserName   string      `xml:"ToUserName"`
	FromUserName string      `xml:"FromUserName"`
	CreateTime   int64       `xml:"CreateTime"`
	MsgType      MessageType `xml:"MsgType"`
	Event        EventType   `xml:"Event"`
}

type EventChangeContact struct {
	Event
	ChangeType ChangeType `xml:"ChangeType"`
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
	UserID         string `xml:"UserID"`
	Name           string `xml:"Name"`
	Department     string `xml:"Department"`
	IsLeaderInDept string `xml:"IsLeaderInDept"`
	Position       string `xml:"Position"`
	Mobile         string `xml:"Mobile"`
	Gender         int8   `xml:"Gender"`
	Email          string `xml:"Email"`
	Status         int32  `xml:"Status"`
	Avatar         string `xml:"Avatar"`
	Alias          string `xml:"Alias"`
	Telephone      string `xml:"Telephone"`
	Address        string `xml:"Address"`
	ExtAttr        struct {
		Text string `xml:",chardata"`
		Item []struct {
			Chardata string `xml:",chardata"`
			Name     string `xml:"Name"`
			Type     string `xml:"Type"`
			Text     struct {
				Text  string `xml:",chardata"`
				Value string `xml:"Value"`
			} `xml:"Text"`
			Web struct {
				Text  string `xml:",chardata"`
				Title string `xml:"Title"`
				URL   string `xml:"Url"`
			} `xml:"Web"`
		} `xml:"Item"`
	} `xml:"ExtAttr"`
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
	UserID         string `xml:"UserID"`
	NewUserID      string `xml:"NewUserID"`
	Name           string `xml:"Name"`
	Department     string `xml:"Department"`
	IsLeaderInDept string `xml:"IsLeaderInDept"`
	Position       string `xml:"Position"`
	Mobile         string `xml:"Mobile"`
	Gender         int8   `xml:"Gender"`
	Email          string `xml:"Email"`
	Status         int32  `xml:"Status"`
	Avatar         string `xml:"Avatar"`
	Alias          string `xml:"Alias"`
	Telephone      string `xml:"Telephone"`
	Address        string `xml:"Address"`
	ExtAttr        struct {
		Text string `xml:",chardata"`
		Item []struct {
			Chardata string `xml:",chardata"`
			Name     string `xml:"Name"`
			Type     string `xml:"Type"`
			Text     struct {
				Text  string `xml:",chardata"`
				Value string `xml:"Value"`
			} `xml:"Text"`
			Web struct {
				Text  string `xml:",chardata"`
				Title string `xml:"Title"`
				URL   string `xml:"Url"`
			} `xml:"Web"`
		} `xml:"Item"`
	} `xml:"ExtAttr"`
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
