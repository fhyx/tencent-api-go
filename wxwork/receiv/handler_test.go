package receiv

import (
	"context"
	"os"
	"testing"

	"go.uber.org/zap"

	"daxv.cn/gopak/tencent-api-go/log"
	"daxv.cn/gopak/tencent-api-go/wxwork/webhook"
)

func TestMain(m *testing.M) {
	lgr, _ := zap.NewDevelopment()
	defer func() {
		_ = lgr.Sync() // flushes buffer, if any
	}()
	sugar := lgr.Sugar()
	log.SetLogger(sugar)

	os.Exit(m.Run())
}

type api struct{}

func (a *api) OnReceived(ctx context.Context, msg interface{}) {
	switch obj := msg.(type) {
	case *EventChangeContactCreateUser:
		logger().Infow("got EventChangeContactCreateUser", "obj", obj)
	case *EventChangeContactUpdateUser:
		logger().Infow("got EventChangeContactUpdateUser", "obj", obj)
	case *EventChangeContactDeleteUser:
		logger().Infow("got EventChangeContactDeleteUser", "obj", obj)
	case *EventChangeContactCreateParty:
		logger().Infow("got EventChangeContactCreateParty", "obj", obj)
	case *EventChangeContactUpdateParty:
		logger().Infow("got EventChangeContactUpdateParty", "obj", obj)
	case *EventSysApprovalChange:
		logger().Infow("got EventSysApprovalChange", "obj", obj)
	default:
		logger().Infow("unhandled", "msg", msg)
	}
}

func TestHandler(t *testing.T) {
	a := &api{}
	cfg := Config{
		AppID:      os.Getenv("WXWORK_CORP_ID"),
		NotifyURI:  os.Getenv("WXWORK_NOTIFY_URI"),
		MsgHandler: a,
	}

	body := "<xml><ToUserName><![CDATA[" + cfg.AppID + "]]></ToUserName><FromUserName><![CDATA[sys]]></FromUserName><CreateTime>1653577739</CreateTime><MsgType><![CDATA[event]]></MsgType><Event><![CDATA[sys_approval_change]]></Event><AgentID>3010040</AgentID><ApprovalInfo><SpNo>202205260001</SpNo><SpName><![CDATA[通用审批]]></SpName><SpStatus>1</SpStatus><TemplateId><![CDATA[C4NyFZG4u4aFe5HQzEVCJYSSVrVK8rokifSiPLXHi]]></TemplateId><ApplyTime>1653577739</ApplyTime><Applyer><UserId><![CDATA[liutao]]></UserId><Party><![CDATA[1]]></Party></Applyer><SpRecord><SpStatus>1</SpStatus><ApproverAttr>1</ApproverAttr><Details><Approver><UserId><![CDATA[liutao]]></UserId></Approver><Speech><![CDATA[]]></Speech><SpStatus>1</SpStatus><SpTime>0</SpTime></Details></SpRecord><Notifyer><UserId><![CDATA[XieChaoNing]]></UserId></Notifyer><StatuChangeEvent>1</StatuChangeEvent></ApprovalInfo></xml>"

	s := &server{}
	if len(cfg.NotifyURI) > 0 {
		s.nh = webhook.NewClient(cfg.NotifyURI)
	}
	msg, err := s.parseMsg([]byte(body))
	if err != nil {
		t.Errorf("parseMsg fail: %s", err)
		return
	}

	a.OnReceived(context.Background(), msg)
	s.notifyMsg(msg)

}
