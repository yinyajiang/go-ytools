package tools

import (
	"github.com/yinyajiang/go-w32"
)

//GetSysBit 获取当前位数
func GetSysBit() int {
	return w32.GetSysBit()
}

//GetSpecialDir 获取指定目录
func GetSpecialDir(t SpecialDirType) string {
	return w32.SHGetSpecialFolderPath(int32(t))
}
