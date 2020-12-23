package tools

import (
	"fmt"
	"path/filepath"

	"github.com/yinyajiang/go-ytools/windows"
)

//DiskUsage 获取路径的磁盘信息
func DiskUsage(path string) (disk DiskStatus, err error) {
	vol := filepath.VolumeName(AbsPath(path))
	freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes, errt := windows.GetDiskFreeSpaceEx(vol)

	return DiskStatus{
		All:  totalNumberOfBytes,
		Used: totalNumberOfBytes - totalNumberOfFreeBytes,
		Free: freeBytesAvailable,
	}, errt
}

//GetFileVersion 获取文件版本信息
func GetFileVersion(file string) string {
	data := windows.GetFileVersionInfo(file)
	if len(data) == 0 {
		return ""
	}
	fileInfo := windows.VerQueryValue(data)

	ver := fmt.Sprintf("%d.%d.%d.%d", fileInfo.FileVersionMS>>16,
		fileInfo.FileVersionMS&0xffff,
		fileInfo.FileVersionLS>>16,
		fileInfo.FileVersionLS&0xffff)
	return ver
}
