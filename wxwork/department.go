package wxwork

import (
	"daxv.cn/gopak/tencent-api-go/client"
	"daxv.cn/gopak/tencent-api-go/models"
)

type Department = models.Department
type Departments = models.Departments

type departmentResponse struct {
	client.Error

	Departments `json:"department"`
}

// FilterDepartment Deprecated with Departments.WithID()
func FilterDepartment(data []Department, id int) (*Department, error) {
	for _, dept := range data {
		if dept.ID == id {
			return &dept, nil
		}
	}
	return nil, ErrNotFound
}

type DepartmentUp = Department
