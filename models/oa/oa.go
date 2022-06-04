package oa

// ApplyEvent 提交审批申请
type ApplyEvent struct {
	// CreatorUserID 申请人userid，此审批申请将以此员工身份提交，申请人需在应用可见范围内
	CreatorUserID string `json:"creator_userid"`
	// TemplateID 模板id。可在“获取审批申请详情”、“审批状态变化回调通知”中获得，也可在审批模板的模板编辑页面链接中获得。暂不支持通过接口提交[打卡补卡][调班]模板审批单。
	TemplateID string `json:"template_id"`
	// UseTemplateApprover 审批人模式：0-通过接口指定审批人、抄送人（此时approver、notifyer等参数可用）; 1-使用此模板在管理后台设置的审批流程，支持条件审批。默认为0
	UseTemplateApprover uint8 `json:"use_template_approver"`
	// Approver 审批流程信息，用于指定审批申请的审批流程，支持单人审批、多人会签、多人或签，可能有多个审批节点，仅use_template_approver为0时生效。
	Approver []Approver `json:"approver"`
	// Notifier 抄送人节点userid列表，仅use_template_approver为0时生效。
	Notifier []string `json:"notifyer"`
	// NotifyType 抄送方式：1-提单时抄送（默认值）； 2-单据通过后抄送；3-提单和单据通过后抄送。仅use_template_approver为0时生效。
	NotifyType uint8 `json:"notify_type,omitempty"`
	// ApplyData 审批申请数据，可定义审批申请中各个控件的值，其中必填项必须有值，选填项可为空，数据结构同“获取审批申请详情”接口返回值中同名参数“apply_data”
	ApplyData Contents `json:"apply_data"`
	// SummaryList 摘要信息，用于显示在审批通知卡片、审批列表的摘要信息，最多3行
	SummaryList []SummaryList `json:"summary_list"`
}

// Approver 审批流程信息
type Approver struct {
	// Attr 节点审批方式：1-或签；2-会签，仅在节点为多人审批时有效
	Attr uint8 `json:"attr"`
	// UserID 审批节点审批人userid列表，若为多人会签、多人或签，需填写每个人的userid
	UserID []string `json:"userid"`
}

// Content 审批申请详情，由多个表单控件及其内容组成，其中包含需要对控件赋值的信息
type Content struct {
	// Control 控件类型：Text-文本；Textarea-多行文本；Number-数字；Money-金额；Date-日期/日期+时间；Selector-单选/多选；；Contact-成员/部门；Tips-说明文字；File-附件；Table-明细；
	Control Control `json:"control"`
	// ID 控件id：控件的唯一id，可通过“获取审批模板详情”接口获取
	ID string `json:"id"`
	// Title 控件名称 ，若配置了多语言则会包含中英文的控件名称
	Title []Text `json:"title"`
	// Value 控件值 ，需在此为申请人在各个控件中填写内容不同控件有不同的赋值参数，具体说明详见附录。模板配置的控件属性为必填时，对应value值需要有值。
	Value ContentValue `json:"value"`
}

func (c *Content) CheckTitle(args ...string) bool {
	if len(args) > 0 {
		for _, t := range c.Title {
			if len(args) == 1 && t.Text == args[0] {
				return true
			}
			if len(args) > 1 {
				if t.Text == args[0] && t.Lang == args[1] {
					return true
				}
			}
		}
	}
	return false
}

// Contents 审批申请详情，由多个表单控件及其内容组成，其中包含需要对控件赋值的信息
type Contents struct {
	// Contents 审批申请详情，由多个表单控件及其内容组成，其中包含需要对控件赋值的信息
	Contents []Content `json:"contents"`
}

// Text 通用文本信息
type Text struct {
	// Text 文字
	Text string `json:"text"`
	// Lang 语言
	Lang string `json:"lang"`
}

// SummaryList 摘要行信息，用于定义某一行摘要显示的内容
type SummaryList struct {
	// SummaryInfo 摘要行信息，用于定义某一行摘要显示的内容
	SummaryInfo []Text `json:"summary_info"`
}

// ContentValue 控件值 ，需在此为申请人在各个控件中填写内容不同控件有不同的赋值参数，具体说明详见附录。模板配置的控件属性为必填时，对应value值需要有值。
type ContentValue struct {
	// Text 文本/多行文本控件（control参数为Text或Textarea）
	Text string `json:"text"`
	// Number 数字控件（control参数为Number）
	Number string `json:"new_number"`
	// Money 金额控件（control参数为Money）
	Money string `json:"new_money"`
	// Date 日期/日期+时间控件（control参数为Date）
	Date ContentDate `json:"date"`
	// Selector 单选/多选控件（control参数为Selector）
	Selector ContentSelector `json:"selector"`
	// Members 成员控件（control参数为Contact，且value参数为members）
	Members []ContentMember `json:"members"`
	// Departments 部门控件（control参数为Contact，且value参数为departments）
	Departments []ContentDepartment `json:"departments"`
	// Files 附件控件（control参数为File，且value参数为files）
	Files []ContentFile `json:"files"`
	// Table 明细控件（control参数为Table）
	Table []ContentTableList `json:"children"`
	// Vacation 假勤组件-请假组件（control参数为Vacation）
	Vacation ContentVacation `json:"vacation"`
	// Location 位置控件（control参数为Location，且value参数为location）
	Location ContentLocation `json:"location"`
	// RelatedApproval 关联审批单控件（control参数为RelatedApproval，且value参数为related_approval）
	RelatedApproval []ContentRelatedApproval `json:"related_approval"`
	// Formula 公式控件（control参数为Formula，且value参数为formula）
	Formula ContentFormula `json:"formula"`
	// DateRange 时长组件（control参数为DateRange，且value参数为date_range）
	DateRange ContentDateRange `json:"date_range"`
}

// ContentDate 日期/日期+时间内容
type ContentDate struct {
	// Type 时间展示类型：day-日期；hour-日期+时间 ，和对应模板控件属性一致
	Type string `json:"type"`
	// Timestamp 时间戳-字符串类型，在此填写日期/日期+时间控件的选择值，以此为准
	Timestamp string `json:"s_timestamp"`
}

// ContentSelector 类型标志，单选/多选控件的config中会包含此参数
type ContentSelector struct {
	// Type 选择方式：single-单选；multi-多选
	Type string `json:"type"`
	// Options 多选选项，多选属性的选择控件允许输入多个
	Options []ContentSelectorOption `json:"options"`
}

// ContentSelectorOption 多选选项，多选属性的选择控件允许输入多个
type ContentSelectorOption struct {
	// Key 选项key，可通过“获取审批模板详情”接口获得
	Key string `json:"key"`
}

// ContentMember 所选成员内容，即申请人在此控件选择的成员，多选模式下可以有多个
type ContentMember struct {
	// UserID 所选成员的userid
	UserID string `json:"userid"`
	// Name 成员名
	Name string `json:"name"`
}

// ContentDepartment 所选部门内容，即申请人在此控件选择的部门，多选模式下可能有多个
type ContentDepartment struct {
	// OpenAPIID 所选部门id
	OpenAPIID string `json:"openapi_id"`
	// Name 所选部门名
	Name string `json:"name"`
}

// ContentFile 附件
type ContentFile struct {
	// FileID 文件id，该id为临时素材上传接口返回的的media_id，注：提单后将作为单据内容转换为长期文件存储；目前一个审批申请单，全局仅支持上传6个附件，否则将失败。
	FileID string `json:"file_id"`
}

// ContentTableList 子明细列表，在此填写子明细的所有子控件的值，子控件的数据结构同一般控件
type ContentTableList struct {
	// List 子明细列表，在此填写子明细的所有子控件的值，子控件的数据结构同一般控件
	List []Content `json:"list"`
}

// ContentVacation 请假内容，即申请人在此组件内选择的请假信息
type ContentVacation struct {
	// Selector 请假类型，所选选项与假期管理关联，为假期管理中的假期类型
	Selector ContentSelector `json:"selector"`
	// Attendance 假勤组件
	Attendance ContentVacationAttendance `json:"attendance"`
}

// ContentVacationAttendance 假勤组件
type ContentVacationAttendance struct {
	// DateRange 假勤组件时间选择范围
	DateRange ContentVacationAttendanceDateRange `json:"date_range"`
	// Type 假勤组件类型：1-请假；3-出差；4-外出；5-加班
	Type uint8 `json:"type"`
}

// ContentVacationAttendanceDateRange 假勤组件时间选择范围
type ContentVacationAttendanceDateRange struct {
	// Type 时间展示类型：day-日期；hour-日期+时间
	Type string `json:"type"`
	//  时长范围
	ContentDateRange
}

// ContentLocation 位置控件
type ContentLocation struct {
	// Latitude 纬度，精确到6位小数
	Latitude string `json:"latitude"`
	// Longitude 经度，精确到6位小数
	Longitude string `json:"longitude"`
	// Title 地点标题
	Title string `json:"title"`
	// Address 地点详情地址
	Address string `json:"address"`
	// Time 选择地点的时间
	Time int `json:"time"`
}

// ContentRelatedApproval 关联审批单控件
type ContentRelatedApproval struct {
	// SpNo 关联审批单的审批单号
	SpNo string `json:"sp_no"`
}

// ContentFormula 公式控件
type ContentFormula struct {
	// Value 公式的值，提交表单时无需填写，后台自动计算
	Value string `json:"value"`
}

// ContentDateRange 时长组件
type ContentDateRange struct {
	// NewBegin 开始时间，unix时间戳
	NewBegin int64 `json:"new_begin"`
	// NewEnd 结束时间，unix时间戳
	NewEnd int64 `json:"new_end"`
	// NewDuration 时长范围，单位秒
	NewDuration int64 `json:"new_duration"`
}

// TemplateDetail 审批模板详情
type TemplateDetail struct {
	// TemplateNames 模板名称，若配置了多语言则会包含中英文的模板名称，默认为zh_CN中文
	TemplateNames []Text `json:"template_names"`
	// TemplateContent 模板控件信息
	TemplateContent TemplateControls `json:"template_content"`
	// Vacation Vacation控件（假勤控件）
	Vacation TemplateControlConfigVacation `json:"vacation_list"`
}

// TemplateControls 模板控件数组。模板详情由多个不同类型的控件组成，控件类型详细说明见附录。
type TemplateControls struct {
	// Controls 模板名称，若配置了多语言则会包含中英文的模板名称，默认为zh_CN中文
	Controls []TemplateControl `json:"controls"`
}

// TemplateControl 模板控件信息
type TemplateControl struct {
	// Property 模板控件属性，包含了模板内控件的各种属性信息
	Property TemplateControlProperty `json:"property"`
	// Config 模板控件配置，包含了部分控件类型的附加类型、属性，详见附录说明。目前有配置信息的控件类型有：Date-日期/日期+时间；Selector-单选/多选；Contact-成员/部门；Table-明细；Attendance-假勤组件（请假、外出、出差、加班）
	Config TemplateControlConfig `json:"config"`
}

// TemplateControlProperty 模板控件属性
type TemplateControlProperty struct {
	// Control 模板控件属性，包含了模板内控件的各种属性信息
	Control Control `json:"control"`
	// ID 模板控件配置，包含了部分控件类型的附加类型、属性，详见附录说明。目前有配置信息的控件类型有：Date-日期/日期+时间；Selector-单选/多选；Contact-成员/部门；Table-明细；Attendance-假勤组件（请假、外出、出差、加班）
	ID string `json:"id"`
	// Title 模板控件配置，包含了部分控件类型的附加类型、属性，详见附录说明。目前有配置信息的控件类型有：Date-日期/日期+时间；Selector-单选/多选；Contact-成员/部门；Table-明细；Attendance-假勤组件（请假、外出、出差、加班）
	Title []Text `json:"title"`
	// Placeholder 模板控件配置，包含了部分控件类型的附加类型、属性，详见附录说明。目前有配置信息的控件类型有：Date-日期/日期+时间；Selector-单选/多选；Contact-成员/部门；Table-明细；Attendance-假勤组件（请假、外出、出差、加班）
	Placeholder []Text `json:"placeholder"`
	// Require 是否必填：1-必填；0-非必填
	Require uint8 `json:"require"`
	// UnPrint 是否参与打印：1-不参与打印；0-参与打印
	UnPrint uint8 `json:"un_print"`
}

// TemplateControlConfig 模板控件配置
type TemplateControlConfig struct {
	// Date Date控件（日期/日期+时间控件）
	Date TemplateControlConfigDate `json:"date"`
	// Selector Selector控件（单选/多选控件）
	Selector TemplateControlConfigSelector `json:"selector"`
	// Contact Contact控件（成员/部门控件）
	Contact TemplateControlConfigContact `json:"contact"`
	// Table Table（明细控件）
	Table TemplateControlConfigTable `json:"table"`
	// Attendance Attendance控件（假勤控件）
	Attendance TemplateControlConfigAttendance `json:"attendance"`
}

// TemplateControlConfigDate 类型标志，日期/日期+时间控件的config中会包含此参数
type TemplateControlConfigDate struct {
	// Type 时间展示类型：day-日期；hour-日期+时间
	Type string `json:"type"`
}

// TemplateControlConfigSelector 类型标志，单选/多选控件的config中会包含此参数
type TemplateControlConfigSelector struct {
	// Type 选择类型：single-单选；multi-多选
	Type string `json:"type"`
	// Options 选项，包含单选/多选控件中的所有选项，可能有多个
	Options []TemplateControlConfigSelectorOption `json:"options"`
}

// TemplateControlConfigSelectorOption 选项，包含单选/多选控件中的所有选项，可能有多个
type TemplateControlConfigSelectorOption struct {
	// Key 选项key，选项的唯一id，可用于发起审批申请，为单选/多选控件赋值
	Key string `json:"key"`
	// Value 选项值，若配置了多语言则会包含中英文的选项值，默认为zh_CN中文
	Value []Text `json:"value"`
}

// TemplateControlConfigContact 类型标志，单选/多选控件的config中会包含此参数
type TemplateControlConfigContact struct {
	// Type 选择类型：single-单选；multi-多选
	Type string `json:"type"`
	// Mode 选择对象：user-成员；department-部门
	Mode string `json:"mode"`
}

// TemplateControlConfigTable 类型标志，明细控件的config中会包含此参数
type TemplateControlConfigTable struct {
	// Children 明细内的子控件，内部结构同controls
	Children []TemplateControl `json:"children"`
}

// TemplateControlConfigAttendance 类型标志，假勤控件的config中会包含此参数
type TemplateControlConfigAttendance struct {
	// DateRange 假期控件属性
	DateRange TemplateControlConfigAttendanceDateRange `json:"date_range"`
	// Type 假勤控件类型：1-请假，3-出差，4-外出，5-加班
	Type uint8 `json:"type"`
}

// TemplateControlConfigAttendanceDateRange 假期控件属性
type TemplateControlConfigAttendanceDateRange struct {
	// Type 时间刻度：hour-精确到分钟, halfday—上午/下午
	Type string `json:"type"`
}

// TemplateControlConfigVacation 类型标志，假勤控件的config中会包含此参数
type TemplateControlConfigVacation struct {
	// Item 单个假期类型属性
	Item []TemplateControlConfigVacationItem `json:"item"`
}

// TemplateControlConfigVacationItem 类型标志，假勤控件的config中会包含此参数
type TemplateControlConfigVacationItem struct {
	// ID 假期类型标识id
	ID int `json:"id"`
	// Name 假期类型名称，默认zh_CN中文名称
	Name []Text `json:"name"`
}

// Control 控件类型
type Control string

// ControlText 文本
const ControlText Control = "Text"

// ControlTextarea 多行文本
const ControlTextarea Control = "Textarea"

// ControlNumber 数字
const ControlNumber Control = "Number"

// ControlMoney 金额
const ControlMoney Control = "Money"

// ControlDate 日期/日期+时间控件
const ControlDate Control = "Date"

// ControlSelector 单选/多选控件
const ControlSelector Control = "Selector"

// ControlContact 成员/部门控件
const ControlContact Control = "Contact"

// ControlTips 说明文字控件
const ControlTips Control = "Tips"

// ControlFile 附件控件
const ControlFile Control = "File"

// ControlTable 明细控件
const ControlTable Control = "Table"

// ControlLocation 位置控件
const ControlLocation Control = "Location"

// ControlRelatedApproval 关联审批单控件
const ControlRelatedApproval Control = "RelatedApproval"

// ControlFormula 公式控件
const ControlFormula Control = "Formula"

// ControlDateRange 时长控件
const ControlDateRange Control = "DateRange"

// ControlVacation 假勤组件-请假组件
const ControlVacation Control = "Vacation"

// ControlAttendance 假勤组件-出差/外出/加班组件
const ControlAttendance Control = "Attendance"

// ApprovalDetail 审批申请详情
type ApprovalDetail struct {
	// SpNo 审批编号
	SpNo string `json:"sp_no"`
	// SpName 审批申请类型名称（审批模板名称）
	SpName string `json:"sp_name"`
	// SpStatus 申请单状态：1-审批中；2-已通过；3-已驳回；4-已撤销；6-通过后撤销；7-已删除；10-已支付
	SpStatus uint8 `json:"sp_status"`
	// TemplateID 审批模板id。可在“获取审批申请详情”、“审批状态变化回调通知”中获得，也可在审批模板的模板编辑页面链接中获得。
	TemplateID string `json:"template_id"`
	// ApplyTime 审批申请提交时间,Unix时间戳
	ApplyTime int64 `json:"apply_time"`
	// Applicant 申请人信息
	Applicant ApprovalDetailApplicant `json:"applyer"`
	// SpRecord 审批流程信息，可能有多个审批节点。
	SpRecord []ApprovalDetailSpRecord `json:"sp_record"`
	// Notifier 抄送信息，可能有多个抄送节点
	Notifier []ApprovalDetailNotifier `json:"notifyer"`
	// ApplyData 审批申请数据
	ApplyData Contents `json:"apply_data"`
	// Comments 审批申请备注信息，可能有多个备注节点
	Comments []ApprovalDetailComment `json:"comments"`
}

// ApprovalDetailApplicant 审批申请详情申请人信息
type ApprovalDetailApplicant struct {
	// UserID 申请人userid
	UserID string `json:"userid"`
	// PartyID 申请人所在部门id
	PartyID string `json:"partyid"`
}

// ApprovalDetailSpRecord 审批流程信息，可能有多个审批节点。
type ApprovalDetailSpRecord struct {
	// SpStatus 审批节点状态：1-审批中；2-已同意；3-已驳回；4-已转审
	SpStatus uint8 `json:"sp_status"`
	// ApproverAttr 节点审批方式：1-或签；2-会签
	ApproverAttr uint8 `json:"approverattr"`
	// Details 审批节点详情,一个审批节点有多个审批人
	Details []ApprovalDetailSpRecordDetail `json:"details"`
}

// ApprovalDetailSpRecordDetail 审批节点详情,一个审批节点有多个审批人
type ApprovalDetailSpRecordDetail struct {
	// Approver 分支审批人
	Approver ApprovalDetailSpRecordDetailApprover `json:"approver"`
	// Speech 审批意见
	Speech string `json:"speech"`
	// SpStatus 分支审批人审批状态：1-审批中；2-已同意；3-已驳回；4-已转审
	SpStatus uint8 `json:"sp_status"`
	// SpTime 节点分支审批人审批操作时间戳，0表示未操作
	SpTime int64 `json:"sptime"`
	// MediaID 节点分支审批人审批意见附件，media_id具体使用请参考：文档-获取临时素材
	MediaID []string `json:"media_id"`
}

// ApprovalDetailSpRecordDetailApprover 分支审批人
type ApprovalDetailSpRecordDetailApprover struct {
	// UserID 分支审批人userid
	UserID string `json:"userid"`
}

// ApprovalDetailNotifier 抄送信息，可能有多个抄送节点
type ApprovalDetailNotifier struct {
	// UserID 节点抄送人userid
	UserID string `json:"userid"`
}

// ApprovalDetailComment 审批申请备注信息，可能有多个备注节点
type ApprovalDetailComment struct {
	// CommentUserInfo 备注人信息
	CommentUserInfo ApprovalDetailCommentUserInfo `json:"commentUserInfo"`
	// CommentTime 备注提交时间戳，Unix时间戳
	CommentTime int64 `json:"commenttime"`
	// CommentTontent 备注文本内容
	CommentTontent string `json:"commentcontent"`
	// CommentID 备注id
	CommentID string `json:"commentid"`
	// MediaID 备注附件id，可能有多个，media_id具体使用请参考：文档-获取临时素材
	MediaID []string `json:"media_id"`
}

// ApprovalDetailCommentUserInfo 备注人信息
type ApprovalDetailCommentUserInfo struct {
	// UserID 备注人userid
	UserID string `json:"userid"`
}

// ApprovalInfoFilter 备注人信息
type ApprovalInfoFilter struct {
	// Key 筛选类型，包括：template_id - 模板类型/模板id；creator - 申请人；department - 审批单提单者所在部门；sp_status - 审批状态。注意:仅“部门”支持同时配置多个筛选条件。不同类型的筛选条件之间为“与”的关系，同类型筛选条件之间为“或”的关系
	Key ApprovalInfoFilterKey `json:"key"`
	// Value 筛选值，对应为：template_id - 模板id；creator - 申请人userid；department - 所在部门id；sp_status - 审批单状态（1-审批中；2-已通过；3-已驳回；4-已撤销；6-通过后撤销；7-已删除；10-已支付）
	Value string `json:"value"`
}

// ApprovalInfoFilterKey 拉取审批筛选类型
type ApprovalInfoFilterKey string

// ApprovalInfoFilterKeyTemplateID 模板类型
const ApprovalInfoFilterKeyTemplateID ApprovalInfoFilterKey = "template_id"

// ApprovalInfoFilterKeyCreator 申请人
const ApprovalInfoFilterKeyCreator ApprovalInfoFilterKey = "creator"

// ApprovalInfoFilterKeyDepartment 审批单提单者所在部门
const ApprovalInfoFilterKeyDepartment ApprovalInfoFilterKey = "department"

// ApprovalInfoFilterKeySpStatus 审批状态
const ApprovalInfoFilterKeySpStatus ApprovalInfoFilterKey = "sp_status"
