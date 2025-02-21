package receiv

const (
	// MessageTypeText 文本消息
	MessageTypeText MessageType = "text"
	// MessageTypeImage 图片消息
	MessageTypeImage MessageType = "image"

	// MessageTypeVoice 语音消息
	MessageTypeVoice MessageType = "voice"

	// MessageTypeVideo 视频消息
	MessageTypeVideo MessageType = "video"

	// MessageTypeLocation 位置消息
	MessageTypeLocation MessageType = "location"

	// MessageTypeLink 链接消息
	MessageTypeLink MessageType = "link"

	// MessageTypeEvent 事件消息
	MessageTypeEvent MessageType = "event"
)

const (
	// EventTypeChangeContact 企业成员变更事件
	EventTypeChangeContact EventType = "change_contact"

	// EventTypeSysApprovalChange 审批申请状态变化回调通知
	EventTypeSysApprovalChange EventType = "sys_approval_change"
	// 审批状态通知事件
	EventTypeOpenApprovalChange EventType = "open_approval_change"
	// 成员关注及取消关注事件
	EventTypeSubscribe EventType = "subscribe"
	// 进入应用
	EventTypeEnterAgent EventType = "enter_agent"
	// 异步任务完成事件推送
	EventTypeBatchJobResult EventType = "batch_job_result"
	// 上报地理位置
	EventTypeLocation EventType = "LOCATION"
	// 点击菜单拉取消息
	EventTypeClick EventType = "click"
	// 点击菜单跳转链接
	EventTypeView EventType = "view"
	// 点击菜单跳转小程序
	EventTypeViewMiniProgram EventType = "view_miniprogram"
	// 扫码推事件的事件推送
	EventTypeScanCodePush EventType = "scancode_push"
	// 扫码推事件且弹出“消息接收中”提示框的事件推送
	EventTypeScanCodeWait EventType = "scancode_waitmsg"
)

// create_user,update_user,delete_user,
// create_party,update_party,delete_party,
const (
	ChangeTypeCreateUser  ChangeType = "create_user"  // 新增成员事件
	ChangeTypeUpdateUser  ChangeType = "update_user"  // 更新成员事件
	ChangeTypeDeleteUser  ChangeType = "delete_user"  // 删除成员事件
	ChangeTypeCreateParty ChangeType = "create_party" // 新增部门事件
	ChangeTypeUpdateParty ChangeType = "update_party" // 更新部门事件
	ChangeTypeDeleteParty ChangeType = "delete_party" // 删除部门事件

	ChangeTypeUpdateTag ChangeType = "update_tag" // 标签成员变更事件
)
