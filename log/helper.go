package ylog

import (
	js "github.com/bitly/go-simplejson"
)

const (
	ostd = 1 << iota
	odbg
	otrace
)

func otoPrefix(f int) string {
	var s string
	switch f {
	case odbg:
		s = "[debug] "
	case otrace:
		s = "[trace] "
	default:
		return ""
	}
	return s
}

func createJSONLog(code int, message string, data *js.Json) *js.Json {
	j := js.New()
	j.Set("code", code)
	j.Set("message", message)
	if data != nil {
		j.Set("data", data)
	} else {
		j.Set("data", "")
	}
	return j
}
