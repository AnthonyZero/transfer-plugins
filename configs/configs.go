package configs

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

//config 项目配置文件信息
type config struct {
	//Base
	Base struct {
		AppName string `toml:"appName"`
		LogPath string `toml:"logPath"`
	} `toml:"base"`
	//Kafka配置
	Kafka struct {
		Addr        string `toml:"addr"`
		Group       string `toml:"group"`
		Assignor    string `toml:"assignor"`
		Oldest      bool   `toml:"oldest"`
		TopicAction string `toml:"topicAction"`
	} `toml:"kafka"`
	//influxdb配置
	InfluxDB struct {
		Addr          string `toml:"addr"`
		Token         string `toml:"token"`
		Bucket        string `toml:"bucket"`
		Org           string `toml:"org"`
		BatchSize     uint   `toml:"batchSize"`
		FlushInterval uint   `toml:"flushInterval"`
	} `toml:"kafka"`
}

//private
var configs = new(config)

//Get 获取配置文件 public
func Get() config {
	return *configs
}

//Init 配置初始化 active -> Env(dev fat uat pro)
func Init(active string) {
	viper.SetConfigName(active + "_configs")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		// viper解析文件错误
		panic(err)
	}

	if err := viper.Unmarshal(configs); err != nil {
		panic(err)
	}

	// 监控配置文件变化
	watchConfig()
}

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("config file changed : %s", in.Name)
	})
}
