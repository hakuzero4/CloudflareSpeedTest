package config

import "github.com/spf13/viper"

type Config struct {
	Dnspod struct {
		ID     int            `yaml:"id"`
		Token  string         `yaml:"token"`
		Domain string         `yaml:"domain"`
		Record map[string]int `yaml:"record"`
	} `yaml:"dnspod"`
	RecordLine string `yaml:"record_line"`
	Dingding   struct {
		Token  string `yaml:"token"`
		Secret string `yaml:"secret"`
	} `yaml:"dingding"`
	Wechat struct {
		Webhook string `yaml:"webhook"`
	} `yaml:"wechat"`
	Notify bool `yaml:"notify"`
}

var C *Config

type config struct {
	path string
}

func Setup(p string) error {
	c := &config{path: p}
	if err := c.readConfig(); err != nil {
		return err
	}
	return nil
}

func (c *config) readConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(c.path)
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return err
	}
	if err := viper.Unmarshal(C); err != nil {
		panic(err)
	}
	return nil
}
