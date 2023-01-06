package wxwork

import (
	"os"
	"testing"

	"go.uber.org/zap"

	"daxv.cn/gopak/tencent-api-go/log"
)

func TestMain(m *testing.M) {
	lgr, _ := zap.NewDevelopment()
	defer func() {
		_ = lgr.Sync() // flushes buffer, if any
	}()
	sugar := lgr.Sugar()
	log.SetLogger(sugar)

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
