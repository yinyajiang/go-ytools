package tools

import (
	"syscall"
	"unsafe"
)

//StringToWCharPtr windows
func StringToWCharPtr(str string) uintptr {
	strp, _ := syscall.UTF16PtrFromString(str)
	return uintptr(unsafe.Pointer(strp))
}

//WCharByteToString windows
func WCharByteToString(buff []byte) string {
	u16 := (*[1 << 29]uint16)(unsafe.Pointer(&buff[0]))[0 : len(buff)/2 : len(buff)/2]
	return syscall.UTF16ToString(u16)
}
