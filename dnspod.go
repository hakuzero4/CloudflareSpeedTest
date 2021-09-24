package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"strconv"
	"strings"
)

// curl -X POST https://dnsapi.cn/Record.Line -d 'login_token=260196,1508ba5160d8647e0b3a660254c3eb2b&format=json&domain=999779.xyz'
type DnsPod struct {
	config *viper.Viper
}

var dnspod = &DnsPod{}

// 境内线路
var record_line = "7=0"

func (d *DnsPod) token() string {
	return fmt.Sprintf("login_token=%s&", d.config.GetString("dnspod.id")+","+d.config.GetString("dnspod.token"))
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
	params.Add("domain", d.config.GetString("dnspod.domain"))
	return params
}

func (d *DnsPod) List() {
	body := strings.NewReader(d.token() + d.format().Encode())
	prefix := "Record.List"
	d.do(prefix, body)
}

func (d *DnsPod) setRecordModify(ip string) {
	for k, v := range d.config.GetStringMapString("dnspod.record") {
		val := d.format()
		val.Set("mx", "0")
		val.Set("record_id", v)
		val.Set("sub_domain", k)
		val.Set("record_type", "A")
		val.Set("record_line", "境内")
		val.Set("record_line_id", record_line)
		val.Set("value", ip)
		body := strings.NewReader(d.token() + val.Encode())
		prefix := "Record.Modify"
		d.do(prefix, body)
	}

}

func init() {
	dnspod.config = viper.New()
	dnspod.config.SetConfigName("config")
	dnspod.config.SetConfigType("yaml")
	dnspod.config.AddConfigPath(".")
	if err := dnspod.config.ReadInConfig(); err != nil {
		panic(err)
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
