package notify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type WechatMsg struct {
	token string
}

type wechatMsgBody struct {
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
}

func (w *WechatMsg) SendMsg(txt string) error {
	msg := &wechatMsgBody{Msgtype: "text", Text: struct {
		Content string "json:\"content\""
	}{txt}}
	m, _ := json.Marshal(msg)
	res, err := http.Post("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="+w.token, "application/json", strings.NewReader(string(m)))
	if err != nil {
		fmt.Printf("企业微信机器人发送消息失败%e", err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		fmt.Printf("消息发送失败%s", b)
	}
	return err
}

func NewWechatMsg(h string) *WechatMsg {
	return &WechatMsg{token: h}
}
