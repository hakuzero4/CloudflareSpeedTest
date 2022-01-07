package notify

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrNotifyStatusIsFalse = errors.New("已关闭消息推送")
)

type iNotify interface {
	SendMsg(interface{}) error
}
type TextStruct struct {
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

type MarkDownStruct struct {
	Msgtype  string `json:"msgtype"`
	Markdown struct {
		Content string `json:"content"`
		Title   string `json:"title"`
		Text    string `json:"text"`
	} `json:"markdown"`
}

type responError struct {
	ErrorCode int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
}

func (r *responError) check(res []byte) (err error) {
	err = json.Unmarshal(res, &r)
	if err != nil {
		return
	}
	if r.ErrorCode != 0 {
		err = fmt.Errorf(r.ErrMsg)
	}
	return
}

var notifyServers []iNotify
