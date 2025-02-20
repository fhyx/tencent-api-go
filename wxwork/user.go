package wxwork

import (
	"daxv.cn/gopak/tencent-api-go/models"
)

type User = models.User
type Users = models.Users
type UserUp = models.User

type usersResponse struct {
	models.WcError

	UserList Users `json:"userlist"`
}

func (ur *usersResponse) Users() Users {
	return ur.UserList
}

// OAuth2UserInfo 为用户 OAuth2 验证登录后的简单信息
type OAuth2UserInfo struct {
	models.WcError

	// 成员UserID。若需要获得用户详情信息，可调用通讯录接口：读取成员。
	// 如果是互联企业/企业互联/上下游，则返回的UserId格式如：CorpId/userid
	UserID string `json:"userid,omitempty"`
	// ?
	DeviceID string `json:"DeviceId,omitempty"`
	// 成员票据，最大为512字节，有效期为1800s。
	// scope为snsapi_privateinfo，且用户在应用可见范围之内时返回此参数。
	// 后续利用该参数可以获取用户信息或敏感信息，参见"获取访问用户敏感信息"。暂时不支持上下游或/企业互联场景
	UserTicket string `json:"user_ticket,omitempty"`

	// 非企业成员的标识，对当前企业唯一。不超过64字节
	OpenId string `json:"openid,omitempty"`
	// 外部联系人id，当且仅当用户是企业的客户，且跟进人在应用的可见范围内时返回。
	// 如果是第三方应用调用，针对同一个客户，同一个服务商不同应用获取到的id相同
	ExternalUserID string `json:"external_userid,omitempty"`
}

type OAuth2UserDetailReq struct {
	UserTicket string `json:"user_ticket"`
}
