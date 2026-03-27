package wxwork

import (
	"testing"
)

func TestTemplateCardButtonInteractionJSON(t *testing.T) {
	buttons := []TemplateCardButton{
		{Text: "按钮1", Style: 1, Key: "btn_key_1"},
		{Text: "按钮2", Style: 2, Key: "btn_key_2"},
	}

	content := &TemplateCardContent{
		CardType: "button_interaction",
		MainTitle: TemplateCardMainTitle{
			Title: "欢迎使用企业微信",
			Desc:  "您的好友正在邀请您加入企业微信",
		},
		SubTitleText: "下载企业微信还能抢红包！",
		TaskID:      "task_id_123",
		ButtonList:  buttons,
	}

	// 验证结构体内容
	if content.CardType != "button_interaction" {
		t.Errorf("expected card_type button_interaction, got %s", content.CardType)
	}
	if content.MainTitle.Title != "欢迎使用企业微信" {
		t.Errorf("expected title '欢迎使用企业微信', got %s", content.MainTitle.Title)
	}
	if len(content.ButtonList) != 2 {
		t.Errorf("expected 2 buttons, got %d", len(content.ButtonList))
	}
	if content.ButtonList[0].Key != "btn_key_1" {
		t.Errorf("expected first button key 'btn_key_1', got %s", content.ButtonList[0].Key)
	}

	// 验证 JSON 序列化 (确保 MarshalJSON 能正常工作)
	// 与 SendTemplateCardMessage 一样，使用 structToMap 转换
	contentMap, err := structToMap(content)
	if err != nil {
		t.Fatalf("structToMap failed: %v", err)
	}
	req := reqMessage{
		ToUser:  []string{"user1"},
		AgentID: 1,
		MsgType: "template_card",
		Content: contentMap,
		IsSafe:  false,
	}

	data, err := req.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	t.Logf("JSON output: %s", string(data))

	// 验证关键字段存在于 JSON 中
	jsonStr := string(data)
	if !contains(jsonStr, `"card_type":"button_interaction"`) {
		t.Errorf("JSON should contain card_type button_interaction, got: %s", jsonStr)
	}
	if !contains(jsonStr, `"task_id":"task_id_123"`) {
		t.Errorf("JSON should contain task_id, got: %s", jsonStr)
	}
	if !contains(jsonStr, `"button_list"`) {
		t.Errorf("JSON should contain button_list, got: %s", jsonStr)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
