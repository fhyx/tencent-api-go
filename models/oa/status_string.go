// Code generated by "stringer -type=SpStatus,ApprovalStatuChangeEvent -linecomment -output status_string.go"; DO NOT EDIT.

package oa

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SpNone-0]
	_ = x[SpInProgress-1]
	_ = x[SpApproved-2]
	_ = x[SpRejected-3]
	_ = x[SpWithdrawn-4]
	_ = x[SpWithdrawAfterAccept-6]
	_ = x[SpDeleted-7]
	_ = x[SpPaid-10]
}

const (
	_SpStatus_name_0 = "nonein progressapprovedrejectedwithdrawn"
	_SpStatus_name_1 = "withdrawn after acceptdeleted"
	_SpStatus_name_2 = "paid"
)

var (
	_SpStatus_index_0 = [...]uint8{0, 4, 15, 23, 31, 40}
	_SpStatus_index_1 = [...]uint8{0, 22, 29}
)

func (i SpStatus) String() string {
	switch {
	case 0 <= i && i <= 4:
		return _SpStatus_name_0[_SpStatus_index_0[i]:_SpStatus_index_0[i+1]]
	case 6 <= i && i <= 7:
		i -= 6
		return _SpStatus_name_1[_SpStatus_index_1[i]:_SpStatus_index_1[i+1]]
	case i == 10:
		return _SpStatus_name_2
	default:
		return "SpStatus(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[AsceNone-0]
	_ = x[AsceSubmit-1]
	_ = x[AsceAccept-2]
	_ = x[AsceReject-3]
	_ = x[AsceForward-4]
	_ = x[AsceRemind-5]
	_ = x[AsceWithdraw-6]
	_ = x[AsceWithdrawAfterAccept-8]
	_ = x[AsceRemark-10]
}

const (
	_ApprovalStatuChangeEvent_name_0 = "AsceNonein progressacceptrejectforwardremindwithdraw"
	_ApprovalStatuChangeEvent_name_1 = "withdrawn after accept"
	_ApprovalStatuChangeEvent_name_2 = "remark"
)

var (
	_ApprovalStatuChangeEvent_index_0 = [...]uint8{0, 8, 19, 25, 31, 38, 44, 52}
)

func (i ApprovalStatuChangeEvent) String() string {
	switch {
	case 0 <= i && i <= 6:
		return _ApprovalStatuChangeEvent_name_0[_ApprovalStatuChangeEvent_index_0[i]:_ApprovalStatuChangeEvent_index_0[i+1]]
	case i == 8:
		return _ApprovalStatuChangeEvent_name_1
	case i == 10:
		return _ApprovalStatuChangeEvent_name_2
	default:
		return "ApprovalStatuChangeEvent(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
