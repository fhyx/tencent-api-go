package exmail

import (
	"encoding/json"
	"fmt"
	"log"

	"daxv.cn/gopak/tencent-api-go/client"
	"daxv.cn/gopak/tencent-api-go/gender"
)

const (
	urlToken    = "https://api.exmail.qq.com/cgi-bin/gettoken"
	urlGetLogin = "https://api.exmail.qq.com/cgi-bin/service/get_login_url"
	urlUserGet  = "https://api.exmail.qq.com/cgi-bin/user/get"
	urlNewCount = "https://api.exmail.qq.com/cgi-bin/mail/newcount"
)

type OpenType uint8

const (
	OTIgnore   OpenType = 0
	OTEnabled  OpenType = 1
	OTDisabled OpenType = 2
)

/*
{
   "errcode": 0,
   "errmsg": "ok",
   "userid": " zhangsan@gzdev.com ",
   "name": "李四",
   "department": [1, 2],
   "position": "后台工程师",
   "mobile": "15913215421",
   "gender": "1",
   "enable": "1",
   "slaves":[ zhangsan@gz.com, zhangsan@bjdev.com],
   "cpwd_login":0
}
*/
type User struct {
	client.Error
	Alias      string        `json:"userid"` // main email
	Name       string        `json:"name"`
	Gender     gender.Gender `json:"gender,omitempty"`
	Title      string        `json:"position,omitempty"`
	ExtId      string        `json:"extid,omitempty"`
	Tel        string        `json:"tel,omitempty"`
	Mobile     string        `json:"mobile,omitempty"`
	Slaves     []string      `json:"slaves"` // email aliases
	Department []int         `json:"department,omitempty"`
	Enable     uint8         `json:"enable,omitempty"`
	// OpenType   OpenType      `json:"OpenType"`
	// TODO: PartyList
}

type loginUrl struct {
	client.Error
	LoginUrl  string `json:"login_url,omitempty"`
	ExpiresIn int64  `json:"expires_in,omitempty"`
}

func GetLoginURL(alias string) (s string, err error) {
	obj := &loginUrl{}
	u := fmt.Sprintf("%s?userid=%s", urlGetLogin, alias)
	err = ApiLogin().c.GetJSON(u, obj)
	if err != nil {
		log.Print(err)
		return "", err
	}
	logger().Debugw("GetLoginURL", "url", u)

	s = obj.LoginUrl

	return
	// return fmt.Sprintf(urlLogin, agent, alias, ticket), nil
}

type newCount struct {
	Alias    string      `json:"alias,omitempty"`
	NewCount json.Number `json:"count,omitempty"`
}

func CountNewMail(alias string) (c int, err error) {
	obj := &newCount{}

	u := fmt.Sprintf("%s?userid=%s", urlNewCount, alias)
	err = ApiCheck().c.GetJSON(u, obj)
	if err != nil {
		return 0, err
	}
	logger().Debugw("CountNewMail %s", u)

	count, err := obj.NewCount.Int64()
	if err != nil {
		log.Print(err)
	}

	c = int(count)
	return
}

func GetUser(alias string) (*User, error) {
	obj := &User{}
	err := ApiContact().c.GetJSON(urlUserGet+"?userid="+alias, obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}
