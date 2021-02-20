package tools

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

//VersionSplit 拆分版本信息
func VersionSplit(str string) []string {
	reg := regexp.MustCompile(`\d+`)
	vers := reg.FindAllString(str, -1)
	var results []string
	for _, ver := range vers {
		ver = strings.TrimLeft(ver, "0")
		if len(ver) > 0 {
			results = append(results, ver)
		}
	}
	return results
}

//CmpVersion 比较版本
func CmpVersion(str1 string, str2 string) int {
	results1 := VersionSplit(str1)
	results2 := VersionSplit(str2)
	min := len(results1)
	if len(results1) > len(results2) {
		min = len(results2)
	}
	for i := 0; i < min; i++ {
		if len(results1[i]) != len(results2[i]) {
			if len(results1[i]) > len(results2[i]) {
				return 1
			}
			return -1
		} else if 0 != strings.Compare(results1[i], results2[i]) {
			return strings.Compare(results1[i], results2[i])
		}
	}

	if len(results1[min:]) > len(results2[min:]) {
		return 1
	} else if len(results1[min:]) < len(results2[min:]) {
		return -1
	}
	return 0
}

//IsInArray ...
func IsInArray(arr, val interface{}) bool {
	arrValueOf := reflect.ValueOf(arr)
	for i := 0; i < arrValueOf.Len(); i++ {
		if reflect.DeepEqual(arrValueOf.Index(i).Interface(), val) {
			return true
		}
	}
	return false
}

//ReplaceString 根据map替换文本内容到
func ReplaceString(strdata string, m map[string]string) string {
	for key, val := range m {
		strdata = strings.ReplaceAll(strdata, key, val)
	}
	return strdata
}

//AttachWait ...
func AttachWait() {
	fmt.Println(os.Getpid())
	<-time.After(time.Second * 30)
}

//RandNum 随机数
func RandNum() int64 {
	return rand.Int63()
}

//RandNumN 随机数
func RandNumN(n int64) int64 {
	return rand.Int63n(n)
}

//CmpDate ...
func CmpDate(t1, t2 time.Time) int64 {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()

	tt1 := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	tt2 := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)

	return tt1.Unix() - tt2.Unix()
}

func testValue(t *testing.T, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expect:%v,got:%v", want, got)
	}
}
