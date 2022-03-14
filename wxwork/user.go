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

// UserAttribute 为用户扩展信息
type UserAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// UserAttributes 为用户扩展信息列表
type UserAttributes struct {
	Attrs []*UserAttribute `json:"attrs,omitempty"`
}

// User 为企业用户信息
// 参数	说明
// errcode	返回码
// errmsg	对返回码的文本描述内容
// userlist	成员列表
// userid	成员UserID。对应管理端的帐号
// name	成员名称，此字段从2019年12月30日起，对新创建第三方应用不再返回，2020年6月30日起，对所有历史第三方应用不再返回，后续第三方仅通讯录应用可获取，第三方页面需要通过通讯录展示组件来展示名字
// mobile	手机号码，第三方仅通讯录应用可获取
// department	成员所属部门id列表，仅返回该应用有查看权限的部门id
// order	部门内的排序值，32位整数，默认为0。数量必须和department一致，数值越大排序越前面。
// position	职务信息；第三方仅通讯录应用可获取
// gender	性别。0表示未定义，1表示男性，2表示女性
// email	邮箱，第三方仅通讯录应用可获取
// is_leader_in_dept	表示在所在的部门内是否为上级；第三方仅通讯录应用可获取
// avatar	头像url。第三方仅通讯录应用可获取
// thumb_avatar	头像缩略图url。第三方仅通讯录应用可获取
// telephone	座机。第三方仅通讯录应用可获取
// alias	别名；第三方仅通讯录应用可获取
// status	激活状态: 1=已激活，2=已禁用，4=未激活，5=退出企业。
// 已激活代表已激活企业微信或已关注微工作台（原企业号）。未激活代表既未激活企业微信又未关注微工作台（原企业号）。
// extattr	扩展属性，第三方仅通讯录应用可获取
// qr_code	员工个人二维码，扫描可添加为外部联系人；第三方仅通讯录应用可获取
// external_profile	成员对外属性，字段详情见对外属性；第三方仅通讯录应用可获取
// external_position	对外职务。 第三方仅通讯录应用可获取
// address	地址，第三方仅通讯录应用可获取
// hide_mobile	是否隐藏手机号
// english_name	英文名
// open_userid	全局唯一。对于同一个服务商，不同应用获取到企业内同一个成员的open_userid是相同的，最多64个字节。仅第三方应用可获取
// main_department	主部门
type User struct {
	UID           string         `json:"userid"`
	Name          string         `json:"name,omitempty"`
	Alias         string         `json:"alias,omitempty"`
	EnglishName   string         `json:"english_name,omitempty"`
	OpenID        string         `json:"open_userid,omitempty"`
	DepartmentIds []int          `json:"department,omitempty"`
	Title         string         `json:"position,omitempty"`
	Mobile        string         `json:"mobile,omitempty"`
	Email         string         `json:"email,omitempty"`
	Tel           string         `json:"telephone,omitempty"`
	Gender        gender.Gender  `json:"gender,omitempty"`
	Status        Status         `json:"status,omitempty"`
	Enabled       int8           `json:"enable,omitempty"`
	Avatar        string         `json:"avatar,omitempty"`
	IsLeader      uint8          `json:"isleader,omitempty"`
	LeaderDepts   []int          `json:"is_leader_in_dept,omitempty"`
	ExtAttr       UserAttributes `json:"extattr,omitempty"`

	ExternalPosition string `json:"external_position,omitempty"`

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
