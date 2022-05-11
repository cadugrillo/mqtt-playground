package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Client struct {
		ClientId          string `yaml:"clientId"`
		ServerAddress     string `yaml:"serverAddress"`
		Qos               int    `yaml:"qos"`
		ConnectionTimeout int    `yaml:"connectionTimeout"`
		WriteTimeout      int    `yaml:"writeTimeout"`
		KeepAlive         int    `yaml:"keepAlive"`
		PingTimeout       int    `yaml:"pingTimeout"`
		ConnectRetry      bool   `yaml:"connectRetry"`
		AutoConnect       bool   `yaml:"autoConnect"`
	} `yaml:"client"`
	Logs struct {
		WriteToLog  bool   `yaml:"writeToLog"`
		WriteToDisk bool   `yaml:"writeToDisk"`
		OutputFile  string `yaml:"outPutFile"`
	} `yaml:"logs"`
	Topics struct {
		Topic []string
	} `yaml:"topics"`
}

var (
	cfg Config
)

func ReadConfig() Config {
	f, err := os.Open("config.yml")
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

func main() {

	cfg = ReadConfig()

	fmt.Println(cfg.Client.ServerAddress)
	fmt.Println(cfg.Topics.Topic)

}
