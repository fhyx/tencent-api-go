package wxbizmsgcrypt

import (
	"testing"
)

func TestEchoStr(t *testing.T) {

	token := "QDG6eK"
	receiverId := "wx5823bf96d3bd56c7"
	encodingAeskey := "jWmYm7qr5nMoAUwZRjGtBxmz3KA1tkAj3ykkR6q2B2C"
	wxcpt := NewWXBizMsgCrypt(token, encodingAeskey, receiverId, XmlType)

	verifyMsgSign := "5c45ff5e21c57e6ad56bac8758b79b1d9ac89fd3"
	verifyTimestamp := "1409659589"
	verifyNonce := "263014780"
	verifyEchoStr := "P9nAzCzyDtyTWESHep1vC5X9xho/qYX3Zpb4yKa9SKld1DsH3Iyt3tP3zNdtp+4RPcs8TgAE7OaBO+FZXvnaqQ=="
	echoStr, cryptErr := wxcpt.VerifyURL(verifyMsgSign, verifyTimestamp, verifyNonce, verifyEchoStr)
	if nil != cryptErr {
		t.Errorf("verifyUrl fail %s", cryptErr)
	} else {
		t.Logf("verifyUrl success echoStr %s", echoStr)
	}
}
