//go:build integration

package wxwork

import (
	"os"
	"strconv"
	"testing"
)

// TestSendTemplateCardMessage 发送模板卡片消息集成测试
//
// 需要设置环境变量:
//   - WXWORK_CORP_ID: 企业ID
//   - WXWORK_APP_SECRET: 应用密钥
//   - WXWORK_TO_USER: 接收人用户ID (可选，默认 "@all")
func TestSendTemplateCardMessage(t *testing.T) {
	corpID := os.Getenv("WXWORK_CORP_ID")
	appSecret := os.Getenv("WXWORK_APP_SECRET")
	if corpID == "" || appSecret == "" {
		t.Skip("skip: WXWORK_CORP_ID or WXWORK_APP_SECRET not set")
	}

	toUser := os.Getenv("WXWORK_TO_USER")
	if toUser == "" {
		t.Skip("skip: WXWORK_TO_USER not set")
	}

	wa := NewAPI(corpID, appSecret)
	agentID := os.Getenv("WXWORK_AGENT_ID")
	if agentID == "" {
		t.Skip("skip: WXWORK_AGENT_ID not set")
	}
	wa.AgentID, _ = strconv.Atoi(agentID)
	recipient := &Recipient{
		UserIDs: []string{toUser},
	}

	buttons := []TemplateCardButton{
		{Text: "同意", Style: 1, Key: "btn_agree"},
		{Text: "拒绝", Style: 2, Key: "btn_reject"},
	}

	err := wa.SendTemplateCardButtonInteraction(
		recipient,
		"您有一个待审批流程",
		"请假申请 - 张三",
		buttons,
		"task_test_123",
		false,
	)
	if err != nil {
		t.Skipf("skip: SendTemplateCardButtonInteraction failed: %v", err)
	}
	t.Logf("SendTemplateCardButtonInteraction ok")
}
