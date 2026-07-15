package wxwork

import (
	"log/slog"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	currentLogLevel := slog.SetLogLoggerLevel(slog.LevelDebug)
	defer slog.SetLogLoggerLevel(currentLogLevel)

	ret := m.Run()

	os.Exit(ret)
}

func TestAPI(t *testing.T) {
	wa := NewAPI(os.Getenv("WXWORK_CORP_ID"), os.Getenv("WXWORK_APP_SECRET"))

	depts, err := wa.ListDepartmentID()
	if err != nil {
		t.Errorf("list dept fail: %s", err)
		return
	}
	t.Logf("list simple ok: %d", len(depts))
	depts, err = wa.ListDepartment()
	if err != nil {
		t.Errorf("list dept fail: %s", err)
		return
	}
	t.Logf("list ok: %d", len(depts))

	t.Logf("depts: %+v", depts)

	deptUsers, err := wa.ListIDs("", 200)
	if err != nil {
		t.Errorf("list ids fail: %s", err)
		return
	}

	t.Logf("deptUsers: %+v", deptUsers)
}
