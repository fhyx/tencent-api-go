package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Error
}

type TokenHolder struct {
	currToken  *Token
	base       string
	method     string
	apiAuths   string
	corpId     string
	corpSecret string
	expiresAt  int64
}

var (
	errEmptyAuths = errors.New("empty auth string or corpId and corpSecret")
)

func NewTokenHolder(baseUrl string) *TokenHolder {
	return &TokenHolder{
		base:   baseUrl,
		method: "POST",
	}
}

func (th *TokenHolder) SetAuth(auths string) {
	th.apiAuths = auths
}

func (th *TokenHolder) SetCorp(id, secret string) {
	th.corpId = id
	th.corpSecret = secret
}

func (th *TokenHolder) Expired() bool {
	return th.expiresAt < time.Now().Unix()
}

func (th *TokenHolder) Valid() bool {
	if th.currToken == nil {
		return false
	}
	return !th.Expired()
}

func (th *TokenHolder) GetAuthToken() (token string, err error) {
	if !th.Valid() {
		debug("token is nil or expired, refreshing it")
		th.currToken, err = th.requestToken()
		if err != nil {
			return "", err
		}
		// log.Print("got token", th.currToken)
		th.expiresAt = time.Now().Unix() + th.currToken.ExpiresIn
	}
	token = th.currToken.AccessToken
	return
}

func (th *TokenHolder) requestToken() (token *Token, err error) {
	var resp []byte
	if th.apiAuths != "" { // for ExMail Old API
		body_str := "grant_type=client_credentials"
		resp, err = DoHTTP("POST", th.base, th.apiAuths, bytes.NewBufferString(body_str))
	} else if th.corpId != "" && th.corpSecret != "" { // for ExWechat and ExMail
		uri := fmt.Sprintf("%s?corpid=%s&corpsecret=%s", th.base, th.corpId, th.corpSecret)
		resp, err = DoHTTP("GET", uri, "", nil)
	} else {
		err = errEmptyAuths
	}

	if err != nil {
		log.Printf(" err %s", err)
		return
	}

	obj := &Token{}
	err = json.Unmarshal(resp, obj)
	if err != nil {
		log.Printf("unmarshal err %s", err)
		return
	}
	token = obj

	return
}
