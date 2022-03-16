package webhook

import (
	"fhyx.online/tencent-api-go/log"
)

func logger() log.Logger {
	return log.GetLogger()
}
