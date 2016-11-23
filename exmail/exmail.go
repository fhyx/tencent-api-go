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
)

func init() {
	holder = client.NewTokenHolder(urlToken)
	auths := os.Getenv("EXMAIL_API_AUTHS")
	if auths != "" {
		holder.SetAuth("Basic " + auths)
	}
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

type newCount struct {
	Alias    string
	NewCount json.Number
}

func CountNewMail(alias string) (int, error) {
	token, err := holder.GetAuthToken()
	if err != nil {
		return 0, err
	}
	auths := "Bearer " + token
	resp, err := client.DoHTTP("POST", urlNewCount, auths, bytes.NewBufferString("alias="+alias))
	if err != nil {
		log.Printf("doHTTP err %s", err)
		return 0, err
	}

	log.Printf("resp: %s", resp)

	obj := &newCount{}

	if e := json.Unmarshal(resp, obj); e != nil {
		log.Printf("unmarshal user err %s", e)
		return 0, e
	}

	count, err := obj.NewCount.Int64()
	if err != nil {
		log.Print(err)
	}
	return int(count), nil
}

func GetUser(alias string) (*User, error) {
	token, err := holder.GetAuthToken()
	if err != nil {
		return nil, err
	}
	auths := "Bearer " + token
	resp, err := client.DoHTTP("POST", urlUserGet, auths, bytes.NewBufferString("alias="+alias))
	if err != nil {
		log.Printf("doHTTP err %s", err)
		return nil, err
	}

	log.Printf("resp: %s", resp)

	exErr := &apiError{Arg: alias}
	if e := json.Unmarshal(resp, exErr); e != nil {
		log.Printf("unmarshal api err %s", e)
		return nil, e
	}

	if exErr.ErrCode != "" {
		log.Printf("apiError %s", exErr)
		return nil, exErr
	}

	obj := &User{}
	if e := json.Unmarshal(resp, obj); e != nil {
		log.Printf("unmarshal user err %s", e)
		return nil, e
	}

	return obj, nil
}
