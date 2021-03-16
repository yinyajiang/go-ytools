package tools

import "fmt"

//StringToFloat 字符串转浮点数
func StringToFloat(s string) (ret float64) {
	fmt.Sscanf(s, "%f", &ret)
	return
}

func StringToUnicode(s string) string {
	us := ""
	for _, a := range s {
		us += fmt.Sprintf(`\u%04x`, a)
	}
	return us
}

func UnicodeToString(us string) string {
	s := ""
	l := len(us) / 6
	for i := 0; i < l; i++ {
		var c rune
		fmt.Sscanf(us[i*6:i*6+6], `\u%04x`, &c)
		s += string(c)
	}
	return s
}
