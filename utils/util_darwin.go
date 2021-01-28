package tools

import (
	"unsafe"
)

//StringToWCharPtr mac
func StringToWCharPtr(str string) uintptr {
	utf32 := []rune(str)
	return uintptr(unsafe.Pointer(&utf32[0]))
}
