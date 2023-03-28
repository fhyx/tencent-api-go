package models

import (
	"fmt"
)

type ErrorCoder interface {
	GetErr() error
	GetErrorCode() int
	GetErrorMsg() string
}

type WcError struct {
	ErrCode int    `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

func (e WcError) GetErr() error {
	if e.ErrCode == 0 {
		return nil
	}

	return fmt.Errorf("errcode: %d, errmsg: %s", e.ErrCode, e.ErrMsg)
}

func (e WcError) GetErrorMsg() string {
	return e.ErrMsg
}

func (e WcError) GetErrorCode() int {
	return e.ErrCode
}
