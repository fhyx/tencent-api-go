package main

import (
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"

	talog "daxv.cn/gopak/tencent-api-go/log"
	"daxv.cn/gopak/tencent-api-go/wxwork/webhook"
)

const (
	defaultWebHookUrlTemplate = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s"
)

var (
	botKey string
	text   string
)

func init() {
	flag.StringVar(&botKey, "key", os.Getenv("WXWORK_BOT_KEY"), "key of chat bot")
	flag.StringVar(&text, "text", "", "message text")
}

func main() {
	lgr, _ := zap.NewDevelopment()
	defer func() {
		_ = lgr.Sync() // flushes buffer, if any
	}()
	sugar := lgr.Sugar()
	talog.SetLogger(sugar)

	flag.Parse()

	if botKey == "" || text == "" {
		flag.Usage()
		return
	}

	msg := webhook.NewTextMessage(text)

	notifier := webhook.NewClient(fmt.Sprintf(defaultWebHookUrlTemplate, botKey))
	_ = notifier.Notify(msg)
}
