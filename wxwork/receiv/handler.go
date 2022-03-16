package receiv

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"fhyx.online/tencent-api-go/wxbizmsgcrypt"
	"fhyx.online/tencent-api-go/wxwork/webhook"
)

type MessageHandler interface {
	OnReceived(ctx context.Context, msg interface{})
}

type Receiver interface {
	http.Handler
	SetMessageHandler(mh MessageHandler)
}

type server struct {
	cpt *wxbizmsgcrypt.WXBizMsgCrypt
	mh  MessageHandler
	nh  webhook.Notifier
}

var _ http.Handler = (*server)(nil)

func NewHandler(token, encKey string, appid, notifyUri string, mh MessageHandler) Receiver {
	s := &server{
		cpt: wxbizmsgcrypt.NewWXBizMsgCrypt(token, encKey, appid, wxbizmsgcrypt.XmlType),
		mh:  mh,
	}
	if len(notifyUri) > 0 {
		s.nh = webhook.NewClient(notifyUri)
	}
	return s
}

func (s *server) SetMessageHandler(mh MessageHandler) {
	s.mh = mh
}

func (s *server) echoTestHandler(rw http.ResponseWriter, req *http.Request) {
	msgSign := req.URL.Query().Get("msg_signature")
	timestamp := req.URL.Query().Get("timestamp")
	nonce := req.URL.Query().Get("nonce")
	echoStr := req.URL.Query().Get("echoStr")
	text, cryptErr := s.cpt.VerifyURL(msgSign, timestamp, nonce, echoStr)
	if cryptErr != nil {
		logger().Infow("VerifyURL fail", "err", cryptErr)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = rw.Write([]byte(text))
}

func (s *server) eventHandler(rw http.ResponseWriter, req *http.Request) {
	msgSign := req.FormValue("msg_signature")
	timestamp := req.FormValue("timestamp")
	nonce := req.FormValue("nonce")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return
	}
	decrypted, cryptErr := s.cpt.DecryptMsg(msgSign, timestamp, nonce, body)
	if nil != cryptErr {
		logger().Infow("decrypt fail", "err", cryptErr)
		s.notifyText("decrypt fail: " + cryptErr.Error())
		return
	}

	m, err := s.parseMsg(decrypted)
	if err != nil {
		s.notifyText("parseMsg fail: " + err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if s.mh != nil {
		s.mh.OnReceived(req.Context(), m)
	}

}

func (s *server) notifyText(msg string) {
	if s.nh != nil {
		_ = s.nh.Notify(webhook.NewTextMessage(msg))
	}
}

func (s *server) parseMsg(body []byte) (interface{}, error) {
	msg := new(Message)
	err := xml.Unmarshal(body, msg)
	if nil != err {
		logger().Infow("Unmarshal fail", "err", err)
	} else {
		logger().Infow("Unmarshal ok", "msg", msg)
	}

	switch msg.MsgType {
	case MessageTypeText:
		var x MessageText
		err = xml.Unmarshal(body, &x)
		return &x, err
		// TODO: more types
	case MessageTypeEvent:
		s.notifyText(fmt.Sprintf("got msg event: %s, change: %s", msg.Event, msg.ChangeType))
		return s.parseEvent(msg.Event, body)
	default:
		return nil, fmt.Errorf("unknown event '%s'", msg.Event)
	}

}

func (s *server) parseEvent(et EventType, body []byte) (interface{}, error) {
	switch et {
	case EventTypeChangeContact:
		var ec EventChangeContact
		err := xml.Unmarshal(body, &ec)
		if err != nil {
			return nil, err
		}
		switch ec.ChangeType {
		case ChangeTypeCreateUser:
			var obj EventChangeContactCreateUser
			err := xml.Unmarshal(body, &obj)
			return &obj, err
		case ChangeTypeUpdateUser:
			var obj EventChangeContactUpdateUser
			err := xml.Unmarshal(body, &obj)
			return &obj, err
		case ChangeTypeDeleteUser:
			var obj EventChangeContactDeleteUser
			err := xml.Unmarshal(body, &obj)
			return &obj, err
		case ChangeTypeCreateParty:
			var obj EventChangeContactCreateParty
			err := xml.Unmarshal(body, &obj)
			return &obj, err
		case ChangeTypeUpdateParty:
			var obj EventChangeContactUpdateParty
			err := xml.Unmarshal(body, &obj)
			return &obj, err
		case ChangeTypeDeleteParty:
			var obj EventChangeContactDeleteParty
			err := xml.Unmarshal(body, &obj)
			return &obj, err
		case ChangeTypeUpdateTag:
			var obj EventChangeContactUpdateTag
			err := xml.Unmarshal(body, &obj)
			return &obj, err
		default:
			return nil, fmt.Errorf("unknown event change contact '%s'", ec.ChangeType)
		}
	default:
		return nil, fmt.Errorf("unknown event type '%s", et)
	}
}

func (s *server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		// 测试回调模式请求
		s.echoTestHandler(rw, req)

	case http.MethodPost:
		// 回调事件
		s.eventHandler(rw, req)

	default:
		// unhandled request method
		rw.WriteHeader(http.StatusNotImplemented)
	}
}
