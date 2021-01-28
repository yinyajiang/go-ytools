package tools

import (
	"unsafe"
)

//StringToWCharPtr mac
func StringToWCharPtr(str string) uintptr {
	utf32 := []rune(str)
	return uintptr(unsafe.Pointer(&utf32[0]))
}

//WCharSplitToString mac
func WCharSplitToString(buff []byte) string {
	u32 := (*[1 << 29]rune)(unsafe.Pointer(&buff[0]))[0 : len(buff)/4 : len(buff)/4]
	return string(u32)
}
