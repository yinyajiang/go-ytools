package tools

import (
	"os"
	"strings"
)

//GetEnv 获取环境变量
func GetEnv(key string) string {
	for _, item := range os.Environ() {
		keyvals := strings.Split(item, "=")
		if keyvals[0] == key {
			return keyvals[1]
		}
	}
	return ""
}
