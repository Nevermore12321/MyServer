package config

import (
	"github.com/spf13/viper"
	"sync"
)

//  用于 只执行一次的任务
var Once sync.Once

// 用于全局读取config 文件内容的 句柄
var ConfigViper *viper.Viper

func init() {

	Once.Do(func() {

		//  将 全局 的 viper 读取配置文件的句柄 初始化
		ConfigViper = viper.New()

		//  viper 读取配置的默认配置
		ConfigViper.SetDefault("test", "value")
		ConfigViper.SetConfigName("configuration")
		ConfigViper.SetConfigType("yaml")
		ConfigViper.AddConfigPath(".")
		ConfigViper.AddConfigPath("$HOME/,MyServer")

		//  读取配置文件
		err := ConfigViper.ReadInConfig()
		if err != nil {

		}
	})

}
