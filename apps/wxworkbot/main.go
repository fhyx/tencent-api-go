package main

import (
	"flag"
	"log/slog"
	"os"

	"daxv.cn/gopak/tencent-api-go/wxwork/webhook"
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
	currentLogLevel := slog.SetLogLoggerLevel(slog.LevelDebug)
	defer slog.SetLogLoggerLevel(currentLogLevel)

	flag.Parse()

	if botKey == "" || text == "" {
		flag.Usage()
		return
	}

	msg := webhook.NewTextMessage(text)

	notifier := webhook.NewWithKey(botKey)
	_ = notifier.Notify(msg)
}
