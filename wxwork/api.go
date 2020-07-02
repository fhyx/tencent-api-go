package wxwork

import (
	"encoding/json"
	"fmt"
	"os"

	"fhyx.online/tencent-api-go/client"
)

const (
	urlToken   = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
	urlGetUser = "https://qyapi.weixin.qq.com/cgi-bin/user/get"
	urlAddUser = "https://qyapi.weixin.qq.com/cgi-bin/user/create"
	urlDelUser = "https://qyapi.weixin.qq.com/cgi-bin/user/delete"

	urlListDept       = "https://qyapi.weixin.qq.com/cgi-bin/department/list"
	urlSimpleListUser = "https://qyapi.weixin.qq.com/cgi-bin/user/simplelist"
	urlListUser       = "https://qyapi.weixin.qq.com/cgi-bin/user/list"

	urlOAuth2GetUser = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo"
)

var (
	_ IClient = (*API)(nil)
)

// API ...
type API struct {
	corpID     string
	corpSecret string
	c          *client.Client
}

func NewAPI(strs ...string) *API {
	corpID := os.Getenv("WXWORK_CORP_ID")
	corpSecret := os.Getenv("WXWORK_CORP_SECRET")
	if len(strs) > 0 && len(strs[0]) > 0 {
		corpID = strs[0]
		if len(strs) > 1 && len(strs[1]) > 0 {
			corpSecret = strs[1]
		}
	}

	if corpID == "" || corpSecret == "" {
		logger().Infow("corpID or corpSecret are empty or not found")
	}

	c := client.NewClient(urlToken)
	c.SetContentType("application/json")
	c.SetCorp(corpID, corpSecret)
	return &API{
		corpID:     corpID,
		corpSecret: corpSecret,
		c:          c,
	}
}

func (a *API) CorpID() string {
	return a.corpID
}

func (a *API) GetUser(userId string) (*User, error) {
	token, err := a.c.GetAuthToken()
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s?access_token=%s&userid=%s", urlGetUser, token, userId)

	user := new(User)
	err = a.c.GetJSON(uri, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *API) AddUser(user *User) (err error) {
	var token string
	token, err = a.c.GetAuthToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s", urlAddUser, token)
	var data []byte
	data, err = json.Marshal(user)
	if err != nil {
		return
	}

	_, err = a.c.Post(uri, data)
	return
}

func (a *API) DeleteUser(userId string) (err error) {
	var token string
	token, err = a.c.GetAuthToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s&userid=%s", urlDelUser, token, userId)

	_, err = a.c.Get(uri)
	return
}

func (a *API) ListDepartment(id int) (data Departments, err error) {
	var token string
	token, err = a.c.GetAuthToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s&id=%d", urlListDept, token, id)

	var ret departmentResponse
	err = a.c.GetJSON(uri, &ret)

	if err == nil {
		data = ret.Departments
	}

	return
}

func (a *API) ListUser(lr ListReq) (data Users, err error) {
	var token string
	token, err = a.c.GetAuthToken()
	if err != nil {
		return
	}

	var prefix = urlListUser
	if lr.IsSimple {
		prefix = urlSimpleListUser
	}
	fc := "0"
	if lr.IncChild {
		fc = "1"
	}
	uri := fmt.Sprintf("%s?access_token=%s&department_id=%d&fetch_child=%s", prefix, token, lr.DeptID, fc)

	var ret usersResponse
	err = a.c.GetJSON(uri, &ret)

	if err == nil {
		data = ret.Users
	}

	return
}

func (a *API) GetOAuth2User(agentID int, code string) (ou *OAuth2UserInfo, err error) {
	var token string
	token, err = a.c.GetAuthToken()
	if err != nil {
		return
	}

	uri := fmt.Sprintf("%s?access_token=%s&agentid=%d&code=%s", urlOAuth2GetUser, token, agentID, code)

	ou = new(OAuth2UserInfo)
	err = a.c.GetJSON(uri, ou)

	return
}
