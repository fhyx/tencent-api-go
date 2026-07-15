package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type Token struct {
	Error

	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
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
		logger().LogAttrs(context.Background(), slog.LevelInfo, "token is nil or expired, refreshing it")
		th.currToken, err = th.requestToken()
		if err != nil {
			logger().LogAttrs(context.Background(), slog.LevelInfo, "GetAuthToken fail",
				slog.Any("err", err),
			)
			return "", err
		}
		logger().LogAttrs(context.Background(), slog.LevelInfo, "GetAuthToken ok",
			slog.Any("token", th.currToken),
		)
		th.expiresAt = time.Now().Unix() + th.currToken.ExpiresIn
	}
	token = th.currToken.AccessToken
	return
}

func (th *TokenHolder) requestToken() (token *Token, err error) {
	var resp []byte
	var uri string
	if th.apiAuths != "" { // for ExMail Old API
		bodyStr := "grant_type=client_credentials"
		uri = th.base
		resp, err = DoHTTPData("POST", uri, th.apiAuths, bytes.NewBufferString(bodyStr))
	} else if th.corpId != "" && th.corpSecret != "" { // for ExWechat and ExMail
		uri = fmt.Sprintf("%s?corpid=%s&corpsecret=%s", th.base, th.corpId, th.corpSecret)
		resp, err = DoHTTPData("GET", uri, "", nil)
	} else {
		err = errEmptyAuths
	}

	if err != nil {
		logger().LogAttrs(context.Background(), slog.LevelInfo, "doHTTP fail",
			slog.Any("err", err),
		)
		return
	}

	obj := &Token{}
	err = json.Unmarshal(resp, obj)
	if err != nil {
		logger().LogAttrs(context.Background(), slog.LevelInfo, "json.unmarshal fail",
			slog.Any("err", err),
		)
		return
	}
	if obj.ErrCode != 0 {
		err = obj.GetErr()
		logger().LogAttrs(context.Background(), slog.LevelInfo, "request token fail",
			slog.String("uri", uri),
			slog.Any("err", err),
		)
		return
	}

	token = obj

	return
}
