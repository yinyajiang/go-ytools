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
