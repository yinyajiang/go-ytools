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

	js "github.com/bitly/go-simplejson"
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

//OpenJSON ...
func OpenJSON(path string) (*js.Json, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return js.NewFromReader(file)
}

//CreateJSON 创建json节点
func CreateJSON() *js.Json {
	return js.New()
}

//MarshalToJSON ...
func MarshalToJSON(path string, j *js.Json) error {
	data, err := j.MarshalJSON()
	if err != nil {
		return err
	}
	return WriteFileString(path, string(data))
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

//CmpDate ...
func CmpDate(t1, t2 time.Time) int64 {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()

	tt1 := time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC)
	tt2 := time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC)

	return tt1.Unix() - tt2.Unix()
}

//AddBodyPara ...
func AddBodyPara(body *string, k string, v interface{}) {
	if len(*body) > 0 {
		*body += "&"
	}
	*body += k + "=" + fmt.Sprint(v)
}

func testValue(t *testing.T, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expect:%v,got:%v", want, got)
	}
}
