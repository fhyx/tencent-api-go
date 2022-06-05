package oa

//go:generate stringer -type=SpStatus,ApprovalStatuChangeEvent -linecomment -output status_string.go

// SpStatus 申请单状态：1-审批中；2-已通过；3-已驳回；4-已撤销；6-通过后撤销；7-已删除；10-已支付
type SpStatus int32

const (
	SpNone       SpStatus = iota // none
	SpInProgress                 // in progress
	SpApproved                   // approved
	SpRejected                   // rejected
	SpWithdrawn                  // withdrawn
	_
	SpWithdrawAfterAccept // withdrawn after accept
	SpDeleted             // deleted
	_
	_
	SpPaid // paid
)

// ApprovalStatuChangeEvent StatuChangeEvent
// 审批申请状态变化类型：1-提单；2-同意；3-驳回；4-转审；5-催办；6-撤销；8-通过后撤销；10-添加备注
type ApprovalStatuChangeEvent int32

const (
	AsceNone     ApprovalStatuChangeEvent = iota
	AsceSubmit                            // in progress
	AsceAccept                            // accept
	AsceReject                            // reject
	AsceForward                           // forward
	AsceRemind                            // remind
	AsceWithdraw                          // withdraw
	_
	AsceWithdrawAfterAccept // withdrawn after accept
	_
	AsceRemark // remark
)
