package models

import (
	"fmt"
)

type ErrorCoder interface {
	GetErrorCode() int
}

type Error struct {
	ErrCode int    `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("errcode: %d, errmsg: %s", e.ErrCode, e.ErrMsg)
}

func (e *Error) GetErrorCode() int {
	return e.ErrCode
}
