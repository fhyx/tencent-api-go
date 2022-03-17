package wxwork

import (
	"fhyx.online/tencent-api-go/models"
)

type User = models.User
type Users = models.Users
type UserUp = models.User

type usersResponse struct {
	models.Error

	UserList Users `json:"userlist"`
}

func (ur *usersResponse) Users() Users {
	return ur.UserList
}

// OAuth2UserInfo 为用户 OAuth2 验证登录后的简单信息
type OAuth2UserInfo struct {
	models.Error

	UserID     string `json:"UserId,omitempty"`
	DeviceID   string `json:"DeviceId,omitempty"`
	UserTicket string `json:"user_ticket,omitempty"`
	OpenId     string `json:"OpenId,omitempty"` // 非企业成员
}
