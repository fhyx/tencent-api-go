package models

import (
	"fmt"
)

type ErrorCoder interface {
	Error() string
	GetErrorCode() int
	GetErrorMsg() string
}

type WcError struct {
	ErrCode int    `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

func (e WcError) Error() string {
	return fmt.Sprintf("errcode: %d, errmsg: %s", e.ErrCode, e.ErrMsg)
}

func (e WcError) GetErrorMsg() string {
	return e.ErrMsg
}

func (e WcError) GetErrorCode() int {
	return e.ErrCode
}
