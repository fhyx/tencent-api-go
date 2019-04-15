package exwechat

import (
	"github.com/wealthworks/go-tencent-api/client"
)

// 部门
type Department struct {
	Id       int    `json:"id"`       // 部门id，32位整型，指定时必须大于1。若不填该参数，将自动生成id
	Name     string `json:"name"`     // 部门名称。长度限制为1~32个字符
	ParentId int    `json:"parentid"` // 父部门id，32位整型
	Order    int    `json:"order"`    // 在父部门中的次序值。order值大的排序靠前。有效的值范围是[0, 2^32)
}

type departmentResponse struct {
	client.Error

	Department []Department `json:"department"`
}

func FilterDepartment(data []Department, id int) (*Department, error) {
	for _, dept := range data {
		if dept.Id == id {
			return &dept, nil
		}
	}
	return nil, ErrNotFound
}
