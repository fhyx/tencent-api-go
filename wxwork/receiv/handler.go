package receiv

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"daxv.cn/gopak/tencent-api-go/wxbizmsgcrypt"
	"daxv.cn/gopak/tencent-api-go/wxwork/webhook"
)

type Config struct {
	AppID       string
	Token       string
	EncodingKey string

	NotifyURI string

	MsgHandler ReceiveHandler // second
	ReqHandler ReauestHandler // first
}

type ReceiveHandler interface {
	OnReceived(ctx context.Context, msg interface{})
}

type ReauestHandler interface {
	OnRequest(req *http.Request, msg interface{})
}

type Receiver interface {
	http.Handler
	SetReceiveHandler(mh ReceiveHandler)
	SetReauestHandler(reqhdl ReauestHandler)
}

type server struct {
	cpt *wxbizmsgcrypt.WXBizMsgCrypt
	mh  ReceiveHandler
	rh  ReauestHandler
	nh  webhook.Notifier
}

var _ http.Handler = (*server)(nil)

func NewHandler(cfg Config) Receiver {
	s := &server{
		cpt: wxbizmsgcrypt.NewWXBizMsgCrypt(
			cfg.Token, cfg.EncodingKey,
			cfg.AppID, wxbizmsgcrypt.XmlType),
		mh: cfg.MsgHandler,
		rh: cfg.ReqHandler,
	}
	if len(cfg.NotifyURI) > 0 {
		s.nh = webhook.NewClient(cfg.NotifyURI)
	}
	return s
}

func (s *server) SetReceiveHandler(mh ReceiveHandler) {
	s.mh = mh
}

func (s *server) SetReauestHandler(reqhdl ReauestHandler) {
	s.rh = reqhdl
}

func (s *server) echoTestHandler(rw http.ResponseWriter, req *http.Request) {
	msgSign := req.URL.Query().Get("msg_signature")
	timestamp := req.URL.Query().Get("timestamp")
	nonce := req.URL.Query().Get("nonce")
	echoStr := req.URL.Query().Get("echostr")
	text, cryptErr := s.cpt.VerifyURL(msgSign, timestamp, nonce, echoStr)
	if cryptErr != nil {
		logger().Infow("verifyURL fail", "err", cryptErr,
			"timestamp", timestamp, "nonce", nonce, "echoStr", echoStr,
			"query", req.URL.RawQuery)
		rw.WriteHeader(http.StatusBadRequest)
		text = []byte(fmt.Sprintf("error#%04d: %s", cryptErr.ErrCode, cryptErr.ErrMsg))
	}
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = rw.Write([]byte(text))
}

func (s *server) eventHandler(rw http.ResponseWriter, req *http.Request) {
	msgSign := req.FormValue("msg_signature")
	timestamp := req.FormValue("timestamp")
	nonce := req.FormValue("nonce")
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}
	decrypted, cryptErr := s.cpt.DecryptMsg(msgSign, timestamp, nonce, body)
	if nil != cryptErr {
		logger().Infow("decrypt fail", "err", cryptErr)
		s.notifyText("error: decrypt fail: " + cryptErr.Error())
		return
	}

	msg, err := s.parseMsg(decrypted)
	if err != nil {
		s.notifyText("error: parseMsg fail: " + err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if s.rh != nil {
		s.rh.OnRequest(req, msg)
	} else if s.mh != nil {
		s.mh.OnReceived(req.Context(), msg)
	} else {
		logger().Infow("without handler")
	}
	s.notifyMsg(msg)

}

func (s *server) notifyText(msg string) {
	if s.nh != nil {
		_ = s.nh.Notify(webhook.NewTextMessage(msg))
	}
}

func (s *server) notifyImage(msg, uri string) {
	if s.nh != nil {
		_ = s.nh.Notify(webhook.NewMarkdownMessage(fmt.Sprintf("![image](%s)\n>%s", uri, msg)))
	}
}

func (s *server) notifyMsg(m interface{}) {
	if s.nh != nil {
		if v, ok := m.(fmt.Stringer); ok {
			var sb strings.Builder
			fmt.Fprintf(&sb, "%s", v)
			if id, ok := m.(IDGetter); ok {
				fmt.Fprintf(&sb, " id=%s", id.GetID())
			}
			if v, ok := m.(NameGetter); ok {
				if s := v.GetName(); len(s) > 0 {
					fmt.Fprintf(&sb, " name='%s'", s)
				}
			}
			if v, ok := m.(MessageGetter); ok {
				if s := v.GetMessage(); len(s) > 0 {
					fmt.Fprintf(&sb, " msg='%s'", s)
				}
			}
			if v, ok := m.(ChangesGetter); ok {
				if cs := v.GetChanges(); len(cs) > 0 {
					fmt.Fprintf(&sb, " chg=%q", strings.Join(cs, ","))
				}
			}
			text := sb.String()
			if v, ok := m.(AvatarGetter); ok {
				if uri := v.GetAvatar(); len(uri) > 0 {
					s.notifyImage(text, v.GetAvatar())
					return
				}
			}
			s.notifyText(text)
		}
	}
}

func (s *server) parseMsg(body []byte) (interface{}, error) {
	msg := new(Message)
	err := xml.Unmarshal(body, msg)
	if nil != err {
		logger().Infow("xml decode failed", "body", string(body), "err", err)
	} else {
		logger().Infow("xml decode ok", "body", string(body), "msg", msg)
	}

	switch msg.MsgType {
	case MessageTypeText:
		var x MessageText
		err = xml.Unmarshal(body, &x)
		return &x, err
		// TODO: more types
	case MessageTypeEvent:
		// s.notifyText(fmt.Sprintf("got msg event: %s, change: %s", msg.Event, msg.ChangeType))
		return s.parseEvent(msg, body)
	default:
		logger().Infow("unknown msg", "MsgType", msg.MsgType)
		return nil, fmt.Errorf("unknown msg '%s'", msg.MsgType)
	}

}

func (s *server) parseEvent(msg *Message, body []byte) (interface{}, error) {
	switch msg.EvnType {
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
	case EventTypeSysApprovalChange:
		var ev EventSysApprovalChange
		err := xml.Unmarshal(body, &ev)
		if err != nil {
			return nil, err
		}
		return &ev, nil
	default:
		logger().Infow("unknown msg", "EvnType", msg.EvnType)
		return nil, fmt.Errorf("unknown event %s from %s", msg.EvnType, msg.FromUserName)
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
