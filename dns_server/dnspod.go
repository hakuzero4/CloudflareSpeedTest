package dns_server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// curl -X POST https://dnsapi.cn/Record.Line -d 'login_token=260196,1508ba5160d8647e0b3a660254c3eb2b&format=json&domain=999779.xyz'
type DnsPod struct {
	config *viper.Viper
}

func NewDnspod() *DnsPod {

	return &DnsPod{}
}

// 境内线路

func (d *DnsPod) token() string {
	return fmt.Sprintf("login_token=%s&", viper.GetString("dnspod.id")+","+viper.GetString("dnspod.token"))
}
func (d *DnsPod) do(prefix string, body io.Reader) {
	req, err := http.NewRequest("POST", "https://dnsapi.cn/"+prefix, body)
	if err != nil {
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
	}
	s, err := ioutil.ReadAll(resp.Body)
	v, err := zhToUnicode(s)
	b, err := prettyprint(v)
	fmt.Printf("%s", b)
	defer resp.Body.Close()
}

func (d *DnsPod) format() url2.Values {
	params := url2.Values{}
	params.Add("format", "json")
	params.Add("domain", viper.GetString("dnspod.domain"))
	return params
}

func (d *DnsPod) List() {
	body := strings.NewReader(d.token() + d.format().Encode())
	prefix := "Record.List"
	d.do(prefix, body)
}

func (d *DnsPod) SetRecordModify(ip string) {
	record_line := viper.GetString("record_line")
	if record_line == "" {
		record_line = "7=0"
	}
	for k, v := range viper.GetStringMapString("dnspod.record") {
		val := d.format()
		val.Set("mx", "0")
		val.Set("record_id", v)
		val.Set("sub_domain", k)
		val.Set("record_type", "A")
		val.Set("record_line_id", record_line)
		val.Set("value", ip)
		body := strings.NewReader(d.token() + val.Encode())
		prefix := "Record.Modify"
		d.do(prefix, body)
	}

}

// unicode 转中文
func zhToUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

// json 格式化
func prettyprint(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
