package config

import "github.com/spf13/viper"

type config struct {
	path string
}

func Init(p string) error {
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
	return nil
}
