package wxwork

// ListReq ...
type ListReq struct {
	DeptID   int  `json:"deptID,omitempty"`
	IncChild bool `json:"incChild,omitempty"`
	IsSimple bool `json:"isSimple,omitempty"`
}

// ListResult ...
type ListResult interface {
	Users() Users
}

// IClient ... interface of API client
type IClient interface {
	ListDepartment(id int) (data Departments, err error)
	ListUser(r ListReq) (res Users, err error)
	// SyncDepartment(data []DepartmentUp) error
	// SyncUser(user UserUp) error
}
