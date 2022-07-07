package wxwork

import (
	"encoding/json"
	"fmt"
	"os"

	"daxv.cn/gopak/tencent-api-go/client"
)

var (
	UriPrefix = "https://qyapi.weixin.qq.com/cgi-bin"
)

var (
	_ IClient = (*API)(nil)
)

// API ...
type API struct {
	c *client.Client

	corpID     string
	corpSecret string

	AgentID int
}

// return new API instance from corpID, corpSecret
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

	c := client.NewClient(UriPrefix + "/gettoken")
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
	user := new(User)
	err := a.c.GetJSON(UriPrefix+"/user/get?userid="+userId, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *API) AddUser(user *User) (err error) {
	var data []byte
	data, err = json.Marshal(user)
	if err != nil {
		return
	}

	_, err = a.c.Post(UriPrefix+"/user/create", data)
	return
}

func (a *API) DeleteUser(userId string) (err error) {
	_, err = a.c.Get(UriPrefix + "/user/delete?userid=" + userId)
	return
}

func (a *API) ListDepartment(id string) (data Departments, err error) {
	var ret departmentResponse
	err = a.c.GetJSON(UriPrefix+"/department/list?id="+id, &ret)

	if err == nil {
		data = ret.Departments
	}

	return
}

// ListUser ...
func (a *API) ListUser(lr ListReq) (ListResult, error) {
	var prefix = UriPrefix + "/user/list"
	if lr.IsSimple {
		prefix = UriPrefix + "/user/simplelist"
	}
	fc := "0"
	if lr.IncChild {
		fc = "1"
	}
	uri := fmt.Sprintf("%s?department_id=%s&fetch_child=%s", prefix, lr.DeptID, fc)

	var ret usersResponse
	err := a.c.GetJSON(uri, &ret)
	if err != nil {
		logger().Infow("getJSON fail", "uri", uri, "lr", lr, "err", err)
		return nil, err
	}

	return &ret, nil
}

func (a *API) GetOAuth2User(code string) (ou *OAuth2UserInfo, err error) {
	uri := fmt.Sprintf("%s/user/getuserinfo?code=%s", UriPrefix, code)

	ou = new(OAuth2UserInfo)
	err = a.c.GetJSON(uri, ou)

	return
}

type activeStatReq struct {
	Date string `json:"date"`
}

type activeStatRes struct {
	client.Error
	ActiveCount int `json:"active_cnt"`
}

// CountActivity ...
func (a *API) CountActivity(date string) (count int, err error) {
	var data []byte
	data, err = json.Marshal(&activeStatReq{Date: date})
	if err != nil {
		return
	}

	var res activeStatRes
	err = a.c.PostJSON(UriPrefix+"/user/get_active_stat", data, &res)
	if err != nil {
		logger().Infow("count activite fail", "date", date, "err", err)
		return
	}

	count = res.ActiveCount
	return
}

type IPListResult struct {
	IPList []string `json:"ip_list"`
}

func (a *API) GetCallbackIP() ([]string, error) {
	var res IPListResult
	err := a.c.GetJSON(UriPrefix+"/getcallbackip", &res)
	if err != nil {
		return nil, err
	}
	return res.IPList, nil
}
