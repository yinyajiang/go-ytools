package tools

import "github.com/yinyajiang/go-ytools/windows"

//GetSysBit 获取当前位数
func GetSysBit() int {
	var si windows.SystemInfo
	windows.GetNativeSystemInfo(&si)
	if si.ProcessorArchitecture == 6 || si.ProcessorArchitecture == 9 {
		return 64
	}
	return 32
}

//GetSysVersion 获取系统版本
func GetSysVersion() string {
	major, minor, _ := windows.RtlGetNtVersionNumbers()
	if major == 6 && minor == 3 {
		return "8.1"
	} else if major == 10 {
		return "10"
	}
	return ""
}

//GetSpecialDir 获取指定目录
func GetSpecialDir(t SpecialDirType) string {
	return windows.SHGetSpecialFolderPath(int32(t))
}
