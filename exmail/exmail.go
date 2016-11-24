package exmail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/wealthworks/go-tencent-api/client"
)

const (
	urlToken     = "https://exmail.qq.com/cgi-bin/token"
	apiBase      = "http://openapi.exmail.qq.com:12211/openapi/"
	urlAuthKey   = "http://openapi.exmail.qq.com:12211/openapi/mail/authkey"
	urlLogin     = "https://exmail.qq.com/cgi-bin/login?fun=bizopenssologin&method=bizauth&agent=%s&user=%s&ticket=%s"
	urlUserGet   = "http://openapi.exmail.qq.com:12211/openapi/user/get"
	urlNewCount  = "http://openapi.exmail.qq.com:12211/openapi/mail/newcount"
	urlPartyList = "http://openapi.exmail.qq.com:12211/openapi/party/list"
)

type OpenType uint8

const (
	OTIgnore   OpenType = 0
	OTEnabled  OpenType = 1
	OTDisabled OpenType = 2
)

var (
	holder *client.TokenHolder
	agent  string
)

func init() {
	holder = client.NewTokenHolder(urlToken)
	auths := os.Getenv("EXMAIL_API_AUTHS")
	if auths != "" {
		holder.SetAuth("Basic " + auths)
	}

	agent = os.Getenv("EXMAIL_LOGIN_AGENT")
}

/*{
"Alias": " test2@gzservice.com",
"Name": "鲍勃",
"Gender": 1,
"SlaveList": "bb@gzdev.com,bo@gzdev.com",
"Position": "工程师",
"Tel": "62394",
"Mobile": "",
"ExtId": "100",
"PartyList": {
	"Count": 3,
	"List": [{ "Value":"部门 a" }
		,{ "Value":"部门 B/部门 b" }
		,{"Value":"部门 c" }
}}*/
type User struct {
	Alias    string   `json:"Alias"`
	Name     string   `json:"Name"`
	Aliases  string   `json:"SlaveList"`
	Gender   uint8    `json:"Gender"`
	Title    string   `json:"Position"`
	ExtId    string   `json:"ExtId"`
	Tel      string   `json:"Tel"`
	Mobile   string   `json:"Mobile"`
	OpenType OpenType `json:"OpenType"`
	// TODO: PartyList
}

type apiError struct {
	Arg     string `json:"arg"`
	ErrCode string `json:"errcode"`
	ErrMsg  string `json:"error"`
}

func (e *apiError) Error() string {
	return fmt.Sprintf("%s:%s %q", e.ErrCode, e.ErrMsg, e.Arg)
}

type authTicket struct {
	Ticket string `json:"auth_key"`
}

func GetAuthTicket(alias string) (string, error) {
	obj := &authTicket{}
	err := request(urlAuthKey, "alias="+alias, obj)
	if err != nil {
		return "", err
	}
	return obj.Ticket, nil
}

func GetLoginURL(alias string) (string, error) {
	ticket, err := GetAuthTicket(alias)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(urlLogin, agent, alias, ticket), nil
}

type newCount struct {
	Alias    string
	NewCount json.Number
}

func CountNewMail(alias string) (int, error) {
	obj := &newCount{}

	err := request(urlNewCount, "alias="+alias, obj)
	if err != nil {
		return 0, err
	}

	count, err := obj.NewCount.Int64()
	if err != nil {
		log.Print(err)
	}
	return int(count), nil
}

func GetUser(alias string) (*User, error) {
	obj := &User{}
	err := request(urlUserGet, "alias="+alias, obj)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func request(url, body string, obj interface{}) error {
	token, err := holder.GetAuthToken()
	if err != nil {
		return err
	}
	auths := "Bearer " + token
	resp, err := client.DoHTTP("POST", url, auths, bytes.NewBufferString(body))
	if err != nil {
		log.Printf("doHTTP err %s", err)
		return err
	}

	log.Printf("resp: %s", resp)

	exErr := &apiError{Arg: body}
	if e := json.Unmarshal(resp, exErr); e != nil {
		log.Printf("unmarshal api err %s", e)
		return e
	}

	if exErr.ErrCode != "" {
		log.Printf("apiError %s", exErr)
		return exErr
	}

	if e := json.Unmarshal(resp, obj); e != nil {
		log.Printf("unmarshal user err %s", e)
		return e
	}

	return nil
}
