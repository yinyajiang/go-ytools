package tools

import "fmt"

//StringToFloat 字符串转浮点数
func StringToFloat(s string) (ret float64) {
	fmt.Sscanf(s, "%f", &ret)
	return
}
