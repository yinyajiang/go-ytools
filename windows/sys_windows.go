package windows

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	//loaddll
	modkernel32  = syscall.NewLazyDLL("kernel32.dll")
	modshell32   = syscall.NewLazyDLL("shell32.dll")
	modmswincore = syscall.NewLazyDLL("Api-ms-win-core-version-l1-1-0.dll")
	modntdll     = syscall.NewLazyDLL("ntdll.dll")

	//kernel32
	getNativeSystemInfo = modkernel32.NewProc("GetNativeSystemInfo")
	getDiskFreeSpaceExW = modkernel32.NewProc("GetDiskFreeSpaceExW")

	//mswincore
	verQueryValueA      = modmswincore.NewProc("VerQueryValueA")
	getFileVersionInfoW = modmswincore.NewProc("GetFileVersionInfoW")

	//ntdll
	rtlGetNtVersionNumbers = modntdll.NewProc("RtlGetNtVersionNumbers")

	//shell32
	sHGetSpecialFolderPathW = modshell32.NewProc("SHGetSpecialFolderPathW")
)

func stringToUTF16Ptr(s string) uintptr {
	strp, _ := syscall.UTF16PtrFromString(s)
	return uintptr(unsafe.Pointer(strp))
}

func spliceToPtr(data []byte) uintptr {
	return uintptr(unsafe.Pointer(&data[0]))
}

func utf16SpliceToPtr(data []uint16) uintptr {
	return uintptr(unsafe.Pointer(&data[0]))
}

//VsFIXEDFILEINFO windows结构
type VsFIXEDFILEINFO struct {
	Signature        uint32
	StrucVersion     uint32
	FileVersionMS    uint32
	FileVersionLS    uint32
	ProductVersionMS uint32
	ProductVersionLS uint32
	FileFlagsMask    uint32
	FileFlags        uint32
	FileOS           uint32
	FileType         uint32
	FileSubtype      uint32
	FileDateMS       uint32
	FileDateLS       uint32
}

//VerQueryValue windowsAPI
func VerQueryValue(data []byte) *VsFIXEDFILEINFO {
	var info *VsFIXEDFILEINFO
	var bytes uint32
	r, _, err := syscall.Syscall6(verQueryValueA.Addr(), 4, spliceToPtr(data), stringToUTF16Ptr("\\"), uintptr(unsafe.Pointer(&info)), uintptr(unsafe.Pointer(&bytes)), 0, 0)
	if r != 1 {
		fmt.Println(err)
		return nil
	}
	return info
}

//GetFileVersionInfo windowsapi
func GetFileVersionInfo(file string) []byte {
	filep, _ := syscall.UTF16PtrFromString(file)
	buff := make([]byte, 2048)
	r, _, err := syscall.Syscall6(getFileVersionInfoW.Addr(), 4, uintptr(unsafe.Pointer(filep)), 0, 2048, uintptr(unsafe.Pointer(&buff[0])), 0, 0)
	if r != 1 {
		fmt.Println(err)
		return []byte{}
	}
	return buff
}

//GetSystemVersion 系统版本
func GetSystemVersion() string {
	version, err := syscall.GetVersion()
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%d.%d (%d)", byte(version), uint8(version>>8), version>>16)
}

//GetDiskFreeSpaceEx windows的GetDiskFreeSpaceEx
func GetDiskFreeSpaceEx(disk string) (freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64, err error) {
	r, _, _ := syscall.Syscall6(getDiskFreeSpaceExW.Addr(), 4,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(disk))),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)), 0, 0)
	err = nil
	if r != 1 {
		err = fmt.Errorf("windows GetDiskFreeSpaceEx error")
	}
	return
}

//SystemInfo windows SYSTEM_INFO结构
type SystemInfo struct {
	ProcessorArchitecture     uint16
	Reserved                  uint16
	PageSize                  uint32
	MinimumApplicationAddress uintptr
	MaximumApplicationAddress uintptr
	ActiveProcessorMask       uint32
	NumberOfProcessors        uint32
	ProcessorType             uint32
	AllocationGranularity     uint32
	ProcessorLevel            uint16
	ProcessorRevision         uint16
}

//GetNativeSystemInfo windows接口
func GetNativeSystemInfo(lpSystemInfo *SystemInfo) {
	syscall.Syscall(getNativeSystemInfo.Addr(), 1, uintptr(unsafe.Pointer(lpSystemInfo)), 0, 0)
}

//RtlGetNtVersionNumbers windowsAPI
func RtlGetNtVersionNumbers() (Major, Minor, BuildNumber uint32) {
	syscall.Syscall(rtlGetNtVersionNumbers.Addr(), 3, uintptr(unsafe.Pointer(&Major)), uintptr(unsafe.Pointer(&Minor)), uintptr(unsafe.Pointer(&BuildNumber)))
	return
}

//SHGetSpecialFolderPath 获取系统目录
func SHGetSpecialFolderPath(ty int32) string {
	buff := make([]uint16, 260)
	r, _, _ := syscall.Syscall6(sHGetSpecialFolderPathW.Addr(), 4, 0, uintptr(unsafe.Pointer(&buff[0])), uintptr(ty), 0, 0, 0)
	if 1 != r {
		return ""
	}
	return syscall.UTF16ToString(buff)
}
