package wxwork

import (
	"fhyx.online/tencent-api-go/client"
)

// Department 部门
type Department struct {
	Id       int    `json:"id"`       // 部门id，32位整型，指定时必须大于1。若不填该参数，将自动生成id
	Name     string `json:"name"`     // 部门名称。长度限制为1~32个字符
	ParentId int    `json:"parentid"` // 父部门id，32位整型
	Order    int    `json:"order"`    // 在父部门中的次序值。order值大的排序靠前。有效的值范围是[0, 2^32)
}

type Departments []Department

func (z Departments) WithID(id int) *Department {
	for _, dept := range z {
		if dept.Id == id {
			return &dept
		}
	}
	return nil
}

// default sort
func (z Departments) Len() int      { return len(z) }
func (z Departments) Swap(i, j int) { z[i], z[j] = z[j], z[i] }
func (z Departments) Less(i, j int) bool {
	return z[i].ParentId == 0 || z[i].ParentId < z[j].ParentId ||
		z[i].ParentId == z[j].ParentId && z[i].Order > z[j].Order
}

type departmentResponse struct {
	client.Error

	Departments `json:"department"`
}

// FilterDepartment Deprecated with Departments.WithID()
func FilterDepartment(data []Department, id int) (*Department, error) {
	for _, dept := range data {
		if dept.Id == id {
			return &dept, nil
		}
	}
	return nil, ErrNotFound
}
