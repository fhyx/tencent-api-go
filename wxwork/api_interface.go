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
	GetUser(userId string) (*User, error)
	AddUser(user *User) (err error)
	DeleteUser(userId string) (err error)
	CountActivity(date string) (int, error)
	GetDepartment(id string) (dept *Department, err error)
	ListDepartment(id ...string) (data Departments, err error)
	ListDepartmentID(id ...string) (data Departments, err error)
	ListUser(r ListReq) (res ListResult, err error)
	// SyncDepartment(data []DepartmentUp) error
	// SyncUser(user UserUp) error
}
