---
title: 新增模板卡片消息发送 (按钮交互型)
type: feat
status: active
date: 2026-03-27
---

# 新增模板卡片消息发送 (按钮交互型)

## Overview

在 `wxwork/api_message.go` 中新增模板卡片消息发送支持，优先实现**按钮交互型 (button_interaction)**。这是现有 `SendTextCardMessage` 的增强版，支持更丰富的交互能力（按钮点击回调、下拉选择等）。

## Problem Statement / Motivation

企业微信模板卡片消息是一种更强大的消息卡片类型，支持：
- 按钮交互（用户点击按钮触发回调）
- 投票选择
- 下拉选择器

已有 plan (`2026-03-26-001-feat-add-template-card-event-plan`) 实现了**接收**模板卡片事件，本 plan 实现**发送**模板卡片消息，形成完整的交互闭环。

## Proposed Solution

参考现有 `SendTextCardMessage` 的实现模式，新增:

1. 模板卡片相关结构体（嵌套结构）
2. `SendTemplateCardMessage` 方法（支持按钮交互型）
3. 可选: 其他卡片类型的便捷方法

### 按钮交互型卡片结构

```go
// 模板卡片消息内容
type TemplateCardContent struct {
    CardType string `json:"card_type"` // 例如 "button_interaction"

    // 卡片来源样式
    Source *TemplateCardSource `json:"source,omitempty"`

    // 右上角菜单
    ActionMenu *TemplateCardActionMenu `json:"action_menu,omitempty"`

    // 重要: main_title 在 button_interaction 类型中是必填的
    MainTitle TemplateCardMainTitle `json:"main_title"`

    // 引用样式
    QuoteArea *TemplateCardQuoteArea `json:"quote_area,omitempty"`

    SubTitleText string `json:"sub_title_text,omitempty"`

    // 二级标题+文本列表
    HorizontalContentList []TemplateCardHorizontalContent `json:"horizontal_content_list,omitempty"`

    // 跳转指引
    JumpList []TemplateCardJump `json:"jump_list,omitempty"`

    // 整体卡片点击事件
    CardAction *TemplateCardAction `json:"card_action,omitempty"`

    // 【关键】task_id 是 button_interaction 类型的必填字段
    TaskID string `json:"task_id"`

    // 下拉选择器
    ButtonSelection *TemplateCardButtonSelection `json:"button_selection,omitempty"`

    // 【关键】按钮列表是 button_interaction 类型的必填字段
    ButtonList []TemplateCardButton `json:"button_list"`
}
```

## Technical Considerations

### 复用现有架构

- `reqMessage.MarshalJSON` 已支持 `msgtype` 多态: `out[x.MsgType] = x.Content`
- 只需构建 `TemplateCardContent` 作为 `Content` 参数传入 `sendMessage`
- 复用现有 `sendMessage` 的群聊/单聊判断逻辑

### msgtype vs card_type 区别

| 字段 | 值 | 说明 |
|------|-----|------|
| `msgtype` | `"template_card"` | 消息类型，固定值 |
| `card_type` | `"button_interaction"` | 卡片类型，表示卡片的子类型 |

### task_id 要求

- button_interaction/vote_interaction/multiple_interaction 类型 **必须** 提供 task_id
- 用于后续更新卡片接口调用
- 格式: 数字、字母和"`_-@`"组成，最长128字节

### button_list 按钮结构

```go
type TemplateCardButton struct {
    Type  int    `json:"type,omitempty"`   // 0=回调事件, 1=跳转URL，默认0
    Text  string `json:"text"`             // 按钮文案，建议不超过10个字
    Style int    `json:"style,omitempty"`  // 按钮样式 1~4
    Key   string `json:"key,omitempty"`   // 回调key，type=0时必填
    URL   string `json:"url,omitempty"`    // 跳转URL，type=1时必填
}
```

### horizontal_content_list type 类型

| type值 | 说明 | 必填字段 |
|--------|------|----------|
| 0或不填 | 非链接 | - |
| 1 | 跳转URL | `url` |
| 2 | 下载附件 | `media_id` |
| 3 | 成员详情 | `userid` |

## Acceptance Criteria

- [ ] 新增 `TemplateCardContent` 及相关嵌套结构体
- [ ] 新增 `SendTemplateCardMessage(recipient, *TemplateCardContent, isSafe) error` 方法
- [ ] 方法签名遵循现有 `SendXxxMessage` 模式
- [ ] `task_id` 字段正确处理（按钮交互型必填）
- [ ] `button_list` 支持多按钮（最多6个）
- [ ] 可选便捷方法: `SendTemplateCardButtonInteraction(...)` 简化常见用法
- [ ] 文档注释遵循项目风格（中文，含官方文档链接）

## Implementation

### 新增方法签名示例

```go
// SendTemplateCardMessage 发送模板卡片消息
//
// 收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
// 否则为单纯的【发送应用消息】接口调用。
//
//	cardType 可选值: button_interaction, vote_interaction, multiple_interaction
//	taskID 按钮交互型/投票选择型/多项选择型卡片必填，用于后续更新卡片接口调用
//
//	参考: https://developer.work.weixin.qq.com/document/path/90236#模板卡片消息
func (a *API) SendTemplateCardMessage(recipient *Recipient, cardType string, content *TemplateCardContent, isSafe bool) error {
    // 构建 content map，注入 card_type
    m := map[string]any{
        "card_type": cardType,
    }
    // ... 反射或手动复制 content 字段到 m
    return a.sendMessage(recipient, "template_card", m, isSafe)
}

// 便捷方法: 按钮交互型
func (a *API) SendTemplateCardButtonInteraction(
    recipient *Recipient,
    title, desc string,
    buttons []TemplateCardButton,
    taskID string,
    isSafe bool,
) error {
    content := &TemplateCardContent{
        CardType: "button_interaction",
        MainTitle: TemplateCardMainTitle{
            Title: title,
            Desc:  desc,
        },
        TaskID:     taskID,
        ButtonList: buttons,
    }
    return a.SendTemplateCardMessage(recipient, "button_interaction", content, isSafe)
}
```

### 复用 vs 新结构

考虑到 `TemplateCardContent` 结构体较复杂，有两种实现策略:

| 策略 | 优点 | 缺点 |
|------|------|------|
| A. 使用 `map[string]any` | 灵活，简单 | 失去类型安全 |
| B. 定义完整结构体 | 类型安全，文档清晰 | 代码量大 |

**推荐策略 A**: 保持与现有消息类型一致的模式（`map[string]any` content），但提供辅助函数/方法简化构建过程。

## Context

### button_interaction 请求示例

```json
{
    "touser": "UserID1|UserID2",
    "toparty": "PartyID1",
    "msgtype": "template_card",
    "agentid": 1,
    "template_card": {
        "card_type": "button_interaction",
        "source": {
            "icon_url": "图片url",
            "desc": "企业微信",
            "desc_color": 1
        },
        "main_title": {
            "title": "欢迎使用企业微信",
            "desc": "您的好友正在邀请您加入企业微信"
        },
        "quote_area": {
            "type": 1,
            "url": "https://work.weixin.qq.com",
            "title": "引用样式",
            "quote_text": "企业微信真好用"
        },
        "sub_title_text": "下载企业微信还能抢红包！",
        "horizontal_content_list": [
            {"keyname": "邀请人", "value": "张三"},
            {"type": 1, "keyname": "官网", "value": "点击访问", "url": "https://work.weixin.qq.com"}
        ],
        "card_action": {
            "type": 2,
            "url": "https://work.weixin.qq.com",
            "appid": "小程序的appid",
            "pagepath": "/index.html"
        },
        "task_id": "task_id_123",
        "button_selection": {
            "question_key": "btn_question_key1",
            "title": "企业微信评分",
            "option_list": [
                {"id": "btn_selection_id1", "text": "100分"},
                {"id": "btn_selection_id2", "text": "101分"}
            ],
            "selected_id": "btn_selection_id1"
        },
        "button_list": [
            {"text": "按钮1", "style": 1, "key": "button_key_1"},
            {"text": "按钮2", "style": 2, "key": "button_key_2"}
        ]
    }
}
```

## Dependencies & Risks

| 依赖/风险 | 说明 |
|----------|------|
| 现有 sendMessage 架构 | 复用现有逻辑，无需修改核心架构 |
| task_id 业务生成 | 需要调用方提供，唯一性由业务保证 |
| 企业微信版本 | 按钮交互型需 3.1.6+，投票/多项选择需 3.1.12+ |

## Sources

- [企业微信模板卡片消息文档](https://developer.work.weixin.qq.com/document/path/90236#模板卡片消息)
- [现有 api_message.go 实现](wxwork/api_message.go)
- [相关: 模板卡片事件接收 plan](./2026-03-26-001-feat-add-template-card-event-plan.md)
