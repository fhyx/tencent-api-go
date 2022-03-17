package models

// Department 部门
// 参数	说明
// errcode	返回码
// errmsg	对返回码的文本描述内容
// department	部门列表数据。
// id	创建的部门id
// name	部门名称，此字段从2019年12月30日起，对新创建第三方应用不再返回，2020年6月30日起，对所有历史第三方应用不再返回，后续第三方仅通讯录应用可获取，第三方页面需要通过通讯录展示组件来展示部门名称
// name_en	英文名称
// parentid	父部门id。根部门为1
// order	在父部门中的次序值。order值大的排序靠前。值范围是[0, 2^32)
type Department struct {
	ID       int    `json:"id"`                 // 部门id，32位整型，指定时必须大于1。若不填该参数，将自动生成id
	Name     string `json:"name"`               // 部门名称。长度限制为1~32个字符
	NameEN   string `json:"name_en,omitempty"`  // 英文名称
	ParentID int    `json:"parentid,omitempty"` // 父部门id，32位整型
	Order    int    `json:"order,omitempty"`    // 在父部门中的次序值。order值大的排序靠前。有效的值范围是[0, 2^32)
}

type Departments []Department

func (z Departments) WithID(id int) *Department {
	for _, dept := range z {
		if dept.ID == id {
			return &dept
		}
	}
	return nil
}

// default sort
func (z Departments) Len() int      { return len(z) }
func (z Departments) Swap(i, j int) { z[i], z[j] = z[j], z[i] }
func (z Departments) Less(i, j int) bool {
	return z[i].ParentID == 0 || z[i].ParentID < z[j].ParentID ||
		z[i].ParentID == z[j].ParentID && z[i].Order > z[j].Order
}
