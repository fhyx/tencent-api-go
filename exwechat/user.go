package exwechat

import (
	"github.com/wealthworks/go-tencent-api/client"
	"github.com/wealthworks/go-tencent-api/gender"
)

type Status uint8

const (
	Disabled Status = 0
	Enabled  Status = 1
)

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
	EnglishName   string         `json:"english_name,omitempty"`
	DepartmentIds []int          `json:"department,omitempty"`
	Title         string         `json:"position,omitempty"`
	Mobile        string         `json:"mobile,omitempty"`
	Email         string         `json:"email,omitempty"`
	Gender        gender.Gender  `json:"gender,omitempty"`
	Status        Status         `json:"enable,omitempty"`
	Avatar        string         `json:"avatar,omitempty"`
	ExtAttr       UserAttributes `json:"extattr,omitempty"`
	client.Error
}

type usersResponse struct {
	client.Error

	Users []User `json:"userlist"`
}
