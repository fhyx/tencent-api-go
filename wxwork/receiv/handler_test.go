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
	case *EventTemplateCard:
		logger().Infow("got EventTemplateCard", "obj", obj)
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

func TestTemplateCardEvent(t *testing.T) {
	body := `<xml>
    <ToUserName><![CDATA[toUser]]></ToUserName>
    <FromUserName><![CDATA[FromUser]]></FromUserName>
    <CreateTime>123456789</CreateTime>
    <MsgType><![CDATA[event]]></MsgType>
    <Event><![CDATA[template_card_event]]></Event>
    <EventKey><![CDATA[key111]]></EventKey>
    <TaskId><![CDATA[taskid111]]></TaskId>
    <CardType><![CDATA[text_notice]]></CardType>
    <ResponseCode><![CDATA[ResponseCode]]></ResponseCode>
    <AgentID>1</AgentID>
    <SelectedItems>
        <SelectedItem>
            <QuestionKey><![CDATA[QuestionKey1]]></QuestionKey>
            <OptionIds>
                <OptionId><![CDATA[OptionId1]]></OptionId>
                <OptionId><![CDATA[OptionId2]]></OptionId>
            </OptionIds>
        </SelectedItem>
        <SelectedItem>
            <QuestionKey><![CDATA[QuestionKey2]]></QuestionKey>
            <OptionIds>
                <OptionId><![CDATA[OptionId3]]></OptionId>
                <OptionId><![CDATA[OptionId4]]></OptionId>
            </OptionIds>
        </SelectedItem>
    </SelectedItems>
</xml>`

	s := &server{}
	msg, err := s.parseMsg([]byte(body))
	if err != nil {
		t.Errorf("parseMsg fail: %s", err)
		return
	}

	evt, ok := msg.(*EventTemplateCard)
	if !ok {
		t.Fatalf("expected *EventTemplateCard, got %T", msg)
	}

	if evt.TaskId != "taskid111" {
		t.Errorf("expected TaskId 'taskid111', got '%s'", evt.TaskId)
	}
	if evt.CardType != "text_notice" {
		t.Errorf("expected CardType 'text_notice', got '%s'", evt.CardType)
	}
	if evt.ResponseCode != "ResponseCode" {
		t.Errorf("expected ResponseCode 'ResponseCode', got '%s'", evt.ResponseCode)
	}
	if evt.EventKey != "key111" {
		t.Errorf("expected EventKey 'key111', got '%s'", evt.EventKey)
	}
	if evt.AgentID != 1 {
		t.Errorf("expected AgentID 1, got %d", evt.AgentID)
	}

	if evt.SelectedItems == nil || len(evt.SelectedItems) != 2 {
		t.Fatalf("expected 2 SelectedItems, got %v", evt.SelectedItems)
	}
	items := evt.SelectedItems
	if items[0].QuestionKey != "QuestionKey1" {
		t.Errorf("expected QuestionKey 'QuestionKey1', got '%s'", items[0].QuestionKey)
	}
	if len(items[0].OptionIds) != 2 {
		t.Errorf("expected 2 OptionIds, got %d", len(items[0].OptionIds))
	}
	if items[0].OptionIds[0] != "OptionId1" {
		t.Errorf("expected OptionId 'OptionId1', got '%s'", items[0].OptionIds[0])
	}

	// Test GetID and GetName
	if evt.GetID() != "taskid111" {
		t.Errorf("GetID() expected 'taskid111', got '%s'", evt.GetID())
	}
	if evt.GetName() != "text_notice" {
		t.Errorf("GetName() expected 'text_notice', got '%s'", evt.GetName())
	}

	t.Logf("EventTemplateCard parsed successfully: %+v", evt)
}

func TestTemplateCardEventWithoutSelectedItems(t *testing.T) {
	body := `<xml>
    <ToUserName><![CDATA[toUser]]></ToUserName>
    <FromUserName><![CDATA[FromUser]]></FromUserName>
    <CreateTime>123456789</CreateTime>
    <MsgType><![CDATA[event]]></MsgType>
    <Event><![CDATA[template_card_event]]></Event>
    <EventKey><![CDATA[key111]]></EventKey>
    <TaskId><![CDATA[taskid222]]></TaskId>
    <CardType><![CDATA[news_notice]]></CardType>
    <ResponseCode><![CDATA[ResponseCode222]]></ResponseCode>
    <AgentID>2</AgentID>
</xml>`

	s := &server{}
	msg, err := s.parseMsg([]byte(body))
	if err != nil {
		t.Errorf("parseMsg fail: %s", err)
		return
	}

	evt, ok := msg.(*EventTemplateCard)
	if !ok {
		t.Fatalf("expected *EventTemplateCard, got %T", msg)
	}

	if evt.TaskId != "taskid222" {
		t.Errorf("expected TaskId 'taskid222', got '%s'", evt.TaskId)
	}
	if evt.CardType != "news_notice" {
		t.Errorf("expected CardType 'news_notice', got '%s'", evt.CardType)
	}
	if evt.SelectedItems != nil {
		t.Errorf("expected nil SelectedItems, got %v", evt.SelectedItems)
	}

	t.Logf("EventTemplateCard without SelectedItems parsed successfully: %+v", evt)
}
