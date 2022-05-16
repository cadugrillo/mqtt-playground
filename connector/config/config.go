package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ClientSub struct {
		ClientId          string `yaml:"clientId"`
		ServerAddress     string `yaml:"serverAddress"`
		Qos               int    `yaml:"qos"`
		ConnectionTimeout int    `yaml:"connectionTimeout"`
		WriteTimeout      int    `yaml:"writeTimeout"`
		KeepAlive         int    `yaml:"keepAlive"`
		PingTimeout       int    `yaml:"pingTimeout"`
		ConnectRetry      bool   `yaml:"connectRetry"`
		AutoConnect       bool   `yaml:"autoConnect"`
		OrderMaters       bool   `yaml:"orderMaters"`
	} `yaml:"clientSub"`
	ClientPub struct {
		ClientId          string `yaml:"clientId"`
		ServerAddress     string `yaml:"serverAddress"`
		Qos               int    `yaml:"qos"`
		ConnectionTimeout int    `yaml:"connectionTimeout"`
		WriteTimeout      int    `yaml:"writeTimeout"`
		KeepAlive         int    `yaml:"keepAlive"`
		PingTimeout       int    `yaml:"pingTimeout"`
		ConnectRetry      bool   `yaml:"connectRetry"`
		AutoConnect       bool   `yaml:"autoConnect"`
		OrderMaters       bool   `yaml:"orderMaters"`
	} `yaml:"clientPub"`
	Logs struct {
		SubPayload bool `yaml:"subPayload"`
		Debug      bool `yaml:"debug"`
		Warning    bool `yaml:"warning"`
		Error      bool `yaml:"error"`
		Critical   bool `yaml:"critical"`
	} `yaml:"logs"`
	Topics struct {
		Topic []string
	} `yaml:"topics"`
}

func ReadConfig() Config {
	f, err := os.Open("./config/config.yml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
