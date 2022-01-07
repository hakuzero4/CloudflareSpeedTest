package notify

import (
	"CloudflareSpeedTest/config"
	"fmt"
	"strings"
)

func AddServer(n iNotify) {
	notifyServers = append(notifyServers, n)
}

func Setup() {
	if config.C.Wechat.Webhook != "" {
		AddServer(NewWechatMsg(config.C.Wechat.Webhook))
	}
	if config.C.Dnspod.Token != "" {
		AddServer(NewDingDingNotify(config.C.Dingding.Token, config.C.Dingding.Secret))
	}
}

func SendTxtMsg(s string) error {
	return send(&TextStruct{
		Msgtype: "text",
		Text: struct {
			Content string "json:\"content\""
		}{
			Content: s,
		},
	})
}

func SendMarkDownMsg(s string) error {
	return send(&MarkDownStruct{
		Msgtype: "markdown",
		Markdown: struct {
			Content string "json:\"content\""
			Title   string "json:\"title\""
			Text    string "json:\"text\""
		}{Title: "消息推送", Content: s, Text: s},
	})
}

func send(s interface{}) error {
	if err := status(); err != nil {
		return err
	}
	var errstrings []string
	for _, sv := range notifyServers {
		if err := sv.SendMsg(s); err != nil {
			errstrings = append(errstrings, err.Error())
		}
	}
	if len(errstrings) > 0 {
		return fmt.Errorf(strings.Join(errstrings, "\n"))
	}
	return nil
}

func status() (err error) {
	if !config.C.Notify {
		return ErrNotifyStatusIsFalse
	}
	return err
}
