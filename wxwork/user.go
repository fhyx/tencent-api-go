package wxwork

import (
	"fhyx.online/tencent-api-go/client"
	"fhyx.online/tencent-api-go/gender"
)

type Status uint8

const (
	SNone     Status = 0
	SActived  Status = 1
	SInactive Status = 2
	SUnlit    Status = 4
)

func (s Status) String() string {
	switch s {
	case SActived:
		return "actived"
	case SInactive:
		return "inactive"
	case SUnlit:
		return "unlit"
	case SNone:
		return "none"
	default:
		return "unknown"
	}
}

// Attribute 为用户扩展信息
type Attribute struct {
	Type    int32        `json:"type,omitempty"`
	Name    string       `json:"name"`
	Text    *attrText    `json:"text,omitempty"`
	Web     *attrWeb     `json:"web,omitempty"`
	MiniApp *attrMiniApp `json:"miniprogram,omitempty"`
}

type attrText struct {
	Value string `json:"value,omitempty"`
}

type attrWeb struct {
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`
}

type attrMiniApp struct {
	AppID    string `json:"appid,omitempty"`
	PagePath string `json:"pagepath,omitempty"`
	Title    string `json:"title,omitempty"`
}

// Attributes 扩展信息列表
type Attributes struct {
	Attrs []Attribute `json:"attrs,omitempty"`
}

type externalProfile struct {
	ExternalCorpName string      `json:"external_corp_name,omitempty"`
	ExternalAttrs    []Attribute `json:"external_attr,omitempty"`
}

// User 为企业用户信息
// 参数	说明
// userid	成员UserID。对应管理端的帐号，企业内必须唯一。不区分大小写，长度为1~64个字节
// name	成员名称；第三方不可获取，调用时返回userid以代替name；代开发自建应用需要管理员授权才返回；对于非第三方创建的成员，第三方通讯录应用也不可获取；未返回name的情况需要通过通讯录展示组件来展示名字
// mobile	手机号码，代开发自建应用需要管理员授权才返回；第三方仅通讯录应用可获取；对于非第三方创建的成员，第三方通讯录应用也不可获取；上游企业不可获取下游企业成员该字段
// department	成员所属部门id列表，仅返回该应用有查看权限的部门id；成员授权模式下，固定返回根部门id，即固定为1。对授权了“组织架构信息”权限的第三方应用，返回成员所属的全部部门id
// order	部门内的排序值，默认为0。数量必须和department一致，数值越大排序越前面。值范围是0, 2^32)。[成员授权模式下不返回该字段
// position	职务信息；代开发自建应用需要管理员授权才返回；第三方仅通讯录应用可获取；对于非第三方创建的成员，第三方通讯录应用也不可获取；上游企业不可获取下游企业成员该字段
// gender	性别。0表示未定义，1表示男性，2表示女性。第三方仅通讯录应用可获取；对于非第三方创建的成员，第三方通讯录应用也不可获取；上游企业不可获取下游企业成员该字段。注：不可获取指返回值0
// email	邮箱，代开发自建应用需要管理员授权才返回；第三方仅通讯录应用可获取；对于非第三方创建的成员，第三方通讯录应用也不可获取；上游企业不可获取下游企业成员该字段
// biz_mail	企业邮箱，代开发自建应用不返回；第三方仅通讯录应用可获取；对于非第三方创建的成员，第三方通讯录应用也不可获取；上游企业不可获取下游企业成员该字段
// is_leader_in_dept	表示在所在的部门内是否为部门负责人，数量与department一致；第三方通讯录应用或者授权了“组织架构信息-应用可获取企业的部门组织架构信息-部门负责人”权限的第三方应用可获取；对于非第三方创建的成员，第三方通讯录应用不可获取；上游企业不可获取下游企业成员该字段
// direct_leader	直属上级UserID，返回在应用可见范围内的直属上级列表，最多有五个直属上级；第三方通讯录应用或者授权了“组织架构信息-应用可获取可见范围内成员组织架构信息-直属上级”权限的第三方应用可获取；对于非第三方创建的成员，第三方通讯录应用不可获取；上游企业不可获取下游企业成员该字段
// avatar	头像url。 第三方仅通讯录应用可获取；对于非第三方创建的成员，第三方通讯录应用也不可获取；上游企业不可获取下游企业成员该字段
// thumb_avatar	头像缩略图url。第三方仅通讯录应用可获取；对于非第三方创建的成员，第三方通讯录应用也不可获取；上游企业不可获取下游企业成员该字段
// telephone	座机。代开发自建应用需要管理员授权才返回；第三方仅通讯录应用可获取；对于非第三方创建的成员，第三方通讯录应用也不可获取；上游企业不可获取下游企业成员该字段
// alias	别名；第三方仅通讯录应用可获取；对于非第三方创建的成员，第三方通讯录应用也不可获取；上游企业不可获取下游企业成员该字段
// extattr	扩展属性，代开发自建应用需要管理员授权才返回；第三方仅通讯录应用可获取；对于非第三方创建的成员，第三方通讯录应用也不可获取；上游企业不可获取下游企业成员该字段
// status	激活状态: 1=已激活，2=已禁用，4=未激活，5=退出企业。
// qr_code	员工个人二维码，扫描可添加为外部联系人(注意返回的是一个url，可在浏览器上打开该url以展示二维码)；第三方仅通讯录应用可获取；对于非第三方创建的成员，第三方通讯录应用也不可获取；上游企业不可获取下游企业成员该字段
// external_profile	成员对外属性，字段详情见对外属性；代开发自建应用需要管理员授权才返回；第三方仅通讯录应用可获取；对于非第三方创建的成员，第三方通讯录应用也不可获取；上游企业不可获取下游企业成员该字段
// external_position	对外职务，如果设置了该值，则以此作为对外展示的职务，否则以position来展示。代开发自建应用需要管理员授权才返回；第三方仅通讯录应用可获取；对于非第三方创建的成员，第三方通讯录应用也不可获取；上游企业不可获取下游企业成员该字段
// address	地址。代开发自建应用需要管理员授权才返回；第三方仅通讯录应用可获取；对于非第三方创建的成员，第三方通讯录应用也不可获取；上游企业不可获取下游企业成员该字段
// open_userid	全局唯一。对于同一个服务商，不同应用获取到企业内同一个成员的open_userid是相同的，最多64个字节。仅第三方应用可获取
// main_department	主部门，仅当应用对主部门有查看权限时返回。
type User struct {
	UID           string        `json:"userid"`
	Name          string        `json:"name,omitempty"`
	Alias         string        `json:"alias,omitempty"`
	OpenID        string        `json:"open_userid,omitempty"`
	DepartmentIds []int         `json:"department,omitempty"`
	Title         string        `json:"position,omitempty"`
	Mobile        string        `json:"mobile,omitempty"`
	Email         string        `json:"email,omitempty"`
	BizMail       string        `json:"biz_mail,omitempty"`
	Tel           string        `json:"telephone,omitempty"`
	Gender        gender.Gender `json:"gender,omitempty"`
	Status        Status        `json:"status,omitempty"`
	Enabled       int8          `json:"enable,omitempty"`
	Avatar        string        `json:"avatar,omitempty"`
	QRCode        string        `json:"qr_code,omitempty"`
	Address       string        `json:"address,omitempty"`
	LeaderInDepts []int         `json:"is_leader_in_dept,omitempty"`
	OrderInDepts  []int         `json:"order,omitempty"`
	ExtAttr       *Attributes   `json:"extattr,omitempty"`

	MainDepartment int32    `json:"main_department,omitempty"`
	DirectLeader   []string `json:"direct_leader,omitempty"`

	ExternalPosition string           `json:"external_position,omitempty"`
	ExternalProfile  *externalProfile `json:"external_profile,omitempty"`

	client.Error
}

func (u User) IsActived() bool {
	return u.Status == SActived
}

func (u User) IsEnabled() bool {
	return u.Enabled == 1
}

// Users ...
type Users []User

type UserUp = User

type usersResponse struct {
	client.Error

	UserList Users `json:"userlist"`
}

func (ur *usersResponse) Users() Users {
	return ur.UserList
}

// OAuth2UserInfo 为用户 OAuth2 验证登录后的简单信息
type OAuth2UserInfo struct {
	UserID     string `json:"UserId,omitempty"`
	DeviceID   string `json:"DeviceId,omitempty"`
	UserTicket string `json:"user_ticket,omitempty"`
	OpenId     string `json:"OpenId,omitempty"` // 非企业成员
	client.Error
}
