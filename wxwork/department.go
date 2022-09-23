package wxwork

import (
	"daxv.cn/gopak/tencent-api-go/client"
	"daxv.cn/gopak/tencent-api-go/models"
)

type Department = models.Department
type Departments = models.Departments

type ListIDsReq struct {
	Cursor string `json:"cursor"`
	Limit  uint32 `json:"limit"`
}

type DeptUser struct {
	UserID string `json:"userid"`     // 用户userid，当用户在多个部门下时会有多条记录
	DeptID uint32 `json:"department"` // 用户所属部门
}

type DeptUsers []DeptUser

type listIDsResponse struct {
	models.Error

	// 分页游标，下次请求时填写以获取之后分页的记录。如果该字段返回空则表示已没有更多数据
	NextCursor string `json:"next_cursor"`

	// 用户-部门关系列表
	DeptUsers DeptUsers `json:"dept_user"`
}

type departmentResponse struct {
	client.Error

	Departments `json:"department_id"`
}

// FilterDepartment Deprecated with Departments.WithID()
func FilterDepartment(data []Department, id int) (*Department, error) {
	for _, dept := range data {
		if dept.ID == uint32(id) {
			return &dept, nil
		}
	}
	return nil, ErrNotFound
}

type DepartmentUp = Department
