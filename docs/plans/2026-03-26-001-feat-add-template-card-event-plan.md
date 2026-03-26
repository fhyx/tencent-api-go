---
title: 新增 template_card_event 模板卡片事件接收
type: feat
status: active
date: 2026-03-26
deepened: true
deepened_date: 2026-03-26
---

# 新增 template_card_event 模板卡片事件接收

## Enhancement Summary

**Deepened on:** 2026-03-26
**Sections enhanced:** All sections with research insights
**Research agents used:** architecture-strategist, code-simplicity-reviewer, repo-research-analyst, best-practices-researcher
**External sources:** 企业微信模板卡片更新接口文档 (path/94888)

### Key Improvements
1. 确认了 `SelectedItems` 可能为空，使用指针切片处理
2. 确定了需要实现 `IDGetter`, `NameGetter` 接口
3. 补充了 `ResponseCode` 字段用于后续更新卡片接口调用
4. 确认了 CardType 五种类型

### New Considerations Discovered
- 模板卡片事件与发送消息紧密关联，ResponseCode 是更新卡片的关键凭证
- SelectedItems 是嵌套数组，需使用 `>` 路径语法处理
- 事件中的 EventKey 与发送时按钮的 key 对应

---

## Overview

在 `wxwork/receiv` 包中新增模板卡片事件 (`template_card_event`) 的接收支持。当用户点击应用下发的模板卡片消息按钮时，企业微信会推送此事件。

## Problem Statement / Motivation

企业微信支持发送模板卡片消息，用户点击按钮后需要能够接收并处理相应的回调事件。目前项目缺少对 `template_card_event` 事件的支持。

**业务价值：** 接收模板卡片事件后才能调用更新卡片接口，实现完整的交互流程（如投票、评分等）。

## Proposed Solution

按照现有事件处理模式，新增：
1. 事件类型常量
2. 事件结构体（包含 SelectedItems 嵌套结构）
3. 在 `parseEvent` 中注册事件解析
4. 实现 `IDGetter`, `NameGetter` 接口

### Research Insights

**Best Practices:**
- XML 嵌套数组使用 `>` 路径语法：`xml:"SelectedItems>SelectedItem"`
- 可选嵌套结构使用指针：`*[]SelectedItem` 或 `omitempty`
- 与项目现有模式保持一致：嵌入式 `Event` 基类

**Architecture Considerations:**
- `TaskId` 可作为业务主键，实现 `IDGetter` 接口
- `CardType` 作为卡片类型标识，实现 `NameGetter` 接口
- `ResponseCode` 字段必须保留，用于调用更新卡片接口

## Technical Considerations

- **XML 嵌套结构**: `SelectedItems` 包含 `SelectedItem` 数组，每个 `SelectedItem` 包含 `QuestionKey` 和 `OptionIds`
- **可选字段**: `SelectedItems` 在不同卡片类型中可能为空，需使用指针处理
- **接口实现**: 参考现有事件实现 `GetID()`, `GetName()` 等接口
- **ResponseCode**: 72小时内有效，用于调用更新卡片接口

## Acceptance Criteria

- [ ] 在 `const.go` 中添加 `EventTypeTemplateCard` 常量
- [ ] 在 `types_event.go` 中新增 `EventTemplateCard` 结构体及嵌套类型
- [ ] 在 `handler.go` 的 `parseEvent` 函数中添加事件解析 case
- [ ] 实现 `GetID()` 和 `GetName()` 接口
- [ ] 编写单元测试验证 XML 解析正确性

## Context

### 官方事件 XML 示例

```xml
<xml>
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
</xml>
```

### CardType 枚举值

| 值 | 说明 | 备注 |
|----|------|------|
| `text_notice` | 文本通知型 | |
| `news_notice` | 图文通知型 | |
| `button_interaction` | 按钮交互型 | 支持按钮点击回调 |
| `vote_interaction` | 投票交互型 | 支持投票选择 |
| `multiple_interaction` | 多项选择型 | 支持多项选择 |

### 字段说明

| 字段 | 说明 |
|------|------|
| EventKey | 与发送模板卡片消息时指定的按钮 btn:key 值相同 |
| TaskId | 与发送模板卡片消息时指定的 task_id 相同 |
| ResponseCode | 用于调用更新卡片接口，72小时内有效，且只能使用一次 |
| QuestionKey | 问题的 key 值 |
| OptionIds | 对应问题的选项列表 |

## Implementation

### const.go 修改

```go
// EventTypeTemplateCard 模板卡片事件（点击模板卡片按钮）
EventTypeTemplateCard EventType = "template_card_event"
```

### types_event.go 新增结构体

```go
// EventTemplateCardSelectedItem 模板卡片选中项
type EventTemplateCardSelectedItem struct {
	QuestionKey string   `xml:"QuestionKey"`           // 问题key
	OptionIds   []string `xml:"OptionIds>OptionId"`    // 选项ID列表
}

// EventTemplateCard 模板卡片事件
type EventTemplateCard struct {
	Event
	TaskId        string                          `xml:"TaskId"`          // 任务ID
	CardType      string                          `xml:"CardType"`        // 卡片类型
	ResponseCode  string                          `xml:"ResponseCode"`    // 更新卡片用的ResponseCode
	SelectedItems []EventTemplateCardSelectedItem `xml:"SelectedItems>SelectedItem"` // 选中项列表
}

func (e *EventTemplateCard) GetID() string {
	return e.TaskId
}

func (e *EventTemplateCard) GetName() string {
	return e.CardType
}
```

**Design Decisions:**
- `SelectedItems` 使用切片 + `xml:"SelectedItems>SelectedItem"` 路径语法处理嵌套数组
- 当无 SelectedItems 时，XML 解析得到空切片（而非 nil），不影响业务逻辑

### handler.go parseEvent 新增 case

```go
case EventTypeTemplateCard:
    var ev EventTemplateCard
    err := xml.Unmarshal(body, &ev)
    return &ev, err
```

## Integration with Update Card API

收到模板卡片事件后，可通过 `ResponseCode` 调用更新卡片接口：

```
POST https://qyapi.weixin.qq.com/cgi-bin/message/update_template_card?access_token=ACCESS_TOKEN
```

**请求示例（更新按钮为不可点击状态）：**
```json
{
    "userids": ["userid1", "userid2"],
    "agentid": 1,
    "response_code": "<事件中的ResponseCode>",
    "button": {
        "replace_name": "已处理"
    }
}
```

**请求示例（更新为新的卡片）：**
```json
{
    "userids": ["userid1"],
    "agentid": 1,
    "response_code": "<事件中的ResponseCode>",
    "template_card": {
        "card_type": "text_notice",
        "main_title": {
            "title": "已更新标题"
        }
    }
}
```

## Dependencies & Risks

| 依赖/风险 | 说明 |
|----------|------|
| 无外部依赖 | 基于现有 wxwork/receiv 包模式 |
| ResponseCode 有效期 | 72小时内有效，只能使用一次，需及时处理 |
| SelectedItems 可能为空 | 使用指针 + omitempty 处理 |

## Sources

- [企业微信模板卡片事件推送文档](https://developer.work.weixin.qq.com/document/path/90376)
- [企业微信模板卡片更新接口文档](https://developer.work.weixin.qq.com/document/path/94888)

---

## Appendix: Go XML 解析最佳实践

### 嵌套数组路径语法

```go
// 处理 SelectedItems>SelectedItem 嵌套
type FormData struct {
    SelectedItems []SelectedItem `xml:"SelectedItems>SelectedItem"`
}

// 处理 OptionIds>OptionId 嵌套
type SelectedItem struct {
    OptionIds []string `xml:"OptionIds>OptionId"`
}
```

### 可选嵌套结构处理

```go
// 整个 SelectedItems 可能不存在
type EventTemplateCard struct {
    SelectedItems *[]SelectedItem `xml:"SelectedItems,omitempty"`
}
```

| 方式 | 适用场景 |
|------|----------|
| `Field []Struct` | 元素存在但可能为空数组 |
| `Field *[]Struct` + `omitempty` | 整个元素可能不存在 |
| `Field *Struct` + `omitempty` | 可选嵌套结构 |
