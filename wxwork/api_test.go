package wxwork

import (
	"os"
	"testing"
)

func TestAPI(t *testing.T) {
	wa := NewAPI(os.Getenv("WXWORK_CORP_ID"), os.Getenv("WXWORK_CONTACTS_SECRET"))

	deptID := "1"
	depts, err := wa.ListDepartment(deptID)
	if err != nil {
		t.Errorf("list dept fail: %s", err)
		return
	}

	t.Logf("depts: %+v", depts)

	deptUsers, err := wa.ListIDs("", 200)
	if err != nil {
		t.Errorf("list ids fail: %s", err)
		return
	}

	t.Logf("deptUsers: %+v", deptUsers)
}
