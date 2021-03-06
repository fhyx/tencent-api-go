package wxwork

// ListReq ...
type ListReq struct {
	DeptID   string `json:"deptID,omitempty"`
	IncChild bool   `json:"incChild,omitempty"`
	IsSimple bool   `json:"isSimple,omitempty"`
}

// ListResult ...
type ListResult interface {
	Users() Users
}

// IClient ... interface of API client
type IClient interface {
	CountActivity(date string) (int, error)
	ListDepartment(id string) (data Departments, err error)
	ListUser(r ListReq) (res ListResult, err error)
	// SyncDepartment(data []DepartmentUp) error
	// SyncUser(user UserUp) error
}
