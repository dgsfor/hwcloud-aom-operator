package utils

import (
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

func Init() {
	SetupConfig()
}

var (
	Config *ini.File
)

func SetupConfig() {
	var err error
	Config, err = ini.Load("config/const/prod.ini")
	if err != nil {
		return
	}
	// 加入环境变量
	Config.ValueMapper = os.ExpandEnv

}
func GetConfig(key string) *ini.Key {
	parts := strings.Split(key, "::")
	section := parts[0]
	keyStr := parts[1]
	return Config.Section(section).Key(keyStr)
}
