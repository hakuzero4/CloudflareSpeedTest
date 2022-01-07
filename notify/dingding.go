package notify

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var _ iNotify = (*Dingding)(nil)

type Dingding struct {
	token, secret string
	resErr        *responError
}

func sign(t int64, secret string) string {
	strToHash := fmt.Sprintf("%d\n%s", t, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(strToHash))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}

func (d *Dingding) SendMsg(content interface{}) error {
	m, _ := json.Marshal(content)
	value := url.Values{}
	value.Set("access_token", d.token)
	if d.secret != "" {
		t := time.Now().UnixNano() / 1e6
		value.Set("timestamp", fmt.Sprintf("%d", t))
		value.Set("sign", sign(t, d.secret))
	}

	req, err := http.NewRequest(http.MethodPost, "https://oapi.dingtalk.com/robot/send", strings.NewReader(string(m)))
	if err != nil {
		return err
	}
	req.URL.RawQuery = value.Encode()
	req.Header.Add("Content-Type", "application/json")
	res, _ := (&http.Client{}).Do(req)
	if err != nil {
		fmt.Printf("钉钉机器人发送消息失败%e", err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err := d.resErr.check(b); err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		fmt.Printf("消息发送失败%s", b)
	}
	return nil
}

func NewDingDingNotify(token, secret string) *Dingding {
	return &Dingding{
		token:  token,
		secret: secret,
		resErr: &responError{},
	}
}
