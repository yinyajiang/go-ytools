package tools

import (
	js "github.com/bitly/go-simplejson"
	"gopkg.in/ini.v1"
)

var (
	//ProductName 产品名
	_ProductName string
)

//ConfigValueFile ...
func ConfigValueFile(file string, keypath ...string) *js.Json {
	j, err := OpenJSON(file)
	if err != nil {
		return nil
	}
	return j.GetPath(keypath...)
}

//ConfigValue 从config.json中获取
func ConfigValue(keypath ...string) *js.Json {
	return ConfigValueFile(LocalPath("config.json"), keypath...)
}

//FrameConfigValue 从frameConfig.json中获取
func FrameConfigValue(keypath ...string) *js.Json {
	return ConfigValueFile(LocalPath("frameConfig.json"), keypath...)
}

//GetProductName ...
func GetProductName() (product string) {
	if len(_ProductName) == 0 {
		product = "unknowProduct"
		if j := FrameConfigValue("productName"); j != nil {
			product = j.MustString()
		}
		_ProductName = product
	} else {
		product = _ProductName
	}
	return
}

//ParseIniString ....
func ParseIniString(key, def, session, file string) string {
	cfg, err := ini.Load(file)
	if err != nil {
		return def
	}
	return cfg.Section(session).Key(key).String()
}
