package exwechat

import (
	"github.com/wealthworks/go-tencent-api/client"
	"github.com/wealthworks/go-tencent-api/gender"
)

type Status uint8

const (
	SNone     Status = 0
	SEnabled  Status = 1
	SDisabled Status = 2
	SUnlit    Status = 4
)

func (s Status) String() string {
	switch s {
	case SEnabled:
		return "enabled"
	case SDisabled:
		return "disabled"
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
type User struct {
	UID           string         `json:"userid"`
	Name          string         `json:"name,omitempty"`
	Alias         string         `json:"alias,omitempty"`
	EnglishName   string         `json:"english_name,omitempty"`
	DepartmentIds []int          `json:"department,omitempty"`
	Title         string         `json:"position,omitempty"`
	Mobile        string         `json:"mobile,omitempty"`
	Email         string         `json:"email,omitempty"`
	Tel           string         `json:"telephone,omitempty"`
	Gender        gender.Gender  `json:"gender,omitempty"`
	Status        Status         `json:"enable,omitempty"`
	Avatar        string         `json:"avatar,omitempty"`
	IsLeader      uint8          `json:"isleader,omitempty"`
	ExtAttr       UserAttributes `json:"extattr,omitempty"`
	client.Error
}

type usersResponse struct {
	client.Error

	Users []User `json:"userlist"`
}

// OAuth2UserInfo 为用户 OAuth2 验证登录后的简单信息
type OAuth2UserInfo struct {
	UserID     string `json:"UserId,omitempty"`
	DeviceID   string `json:"DeviceId,omitempty"`
	UserTicket string `json:"user_ticket,omitempty"`
	OpenId     string `json:"OpenId,omitempty"` // 非企业成员
	client.Error
}
