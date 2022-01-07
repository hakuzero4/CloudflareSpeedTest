package notify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var _ iNotify = (*WechatMsg)(nil)

type WechatMsg struct {
	token  string
	resErr *responError
}

func (w *WechatMsg) SendMsg(content interface{}) error {
	m, _ := json.Marshal(content)
	res, err := http.Post("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="+w.token, "application/json", strings.NewReader(string(m)))
	if err != nil {
		return fmt.Errorf("企业微信机器人发送消息失败%e", err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("消息发送失败%s", b)
	}
	if err := w.resErr.check(b); err != nil {
		return err
	}

	return err
}

func NewWechatMsg(h string) *WechatMsg {
	return &WechatMsg{token: h, resErr: &responError{}}
}
