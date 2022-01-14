package dns_server

import (
	"CloudflareSpeedTest/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// curl -X POST https://dnsapi.cn/Record.Line -d 'login_token=260196,1508ba5160d8647e0b3a660254c3eb2b&format=json&domain=999779.xyz'
type DnsPod struct {
}

func NewDnspod() *DnsPod {
	return &DnsPod{}
}

// 境内线路

func (d *DnsPod) token() string {
	return fmt.Sprintf("login_token=%d,%s&", config.C.Dnspod.ID, config.C.Dnspod.Token)
}
func (d *DnsPod) do(prefix string, body io.Reader) error {
	req, err := http.NewRequest("POST", "https://dnsapi.cn/"+prefix, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	// _, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	return err
	// }
	// v, _ := zhToUnicode(s)
	// b, err := prettyprint(v)
	// fmt.Printf("%s", b)
	defer resp.Body.Close()
	return err
}

func (d *DnsPod) format() url.Values {
	params := url.Values{}
	params.Add("format", "json")
	params.Add("domain", config.C.Dnspod.Domain)
	return params
}

func (d *DnsPod) List() {
	body := strings.NewReader(d.token() + d.format().Encode())
	prefix := "Record.List"
	d.do(prefix, body)
}

func (d *DnsPod) SetRecordModify(ip string) {
	for k, v := range config.C.Dnspod.Record {
		val := d.format()
		val.Set("mx", "0")
		val.Set("record_id", strconv.Itoa(v))
		val.Set("sub_domain", k)
		val.Set("record_type", "A")
		val.Set("record_line_id", recordLineFor())
		val.Set("value", ip)
		body := strings.NewReader(d.token() + val.Encode())
		prefix := "Record.Modify"
		d.do(prefix, body)
	}

}

func recordLineFor() string {
	result := DEFUALT
	if len(config.C.RecordLine) > 0 {
		result = config.C.RecordLine
	}
	return result
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
