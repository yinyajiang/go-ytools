package tools

import (
	"runtime"
)

//GetSysBit 获取当前位数
func GetSysBit() int {
	if "amd64" == runtime.GOARCH {
		return 64
	}
	return 32
}

//GetSpecialDir 获取指定目录
func GetSpecialDir(t SpecialDirType) string {
	if t == LocalAppdata {
		return GetEnv("HOME")
	}
	return ""
}

//GetSysVersion 获取系统版本
func GetSysVersion() string {
	return ""
}
