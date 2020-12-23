package tools

import (
	"testing"
)

func TestPathName(t *testing.T) {
	type testPara struct {
		input1 string
		want   string
	}

	test := map[string]testPara{
		"1": {`file/name1`, "name1"},
		"2": {`file/name2.2`, "name2.2"},
		"3": {`/file\name3.3`, "name3.3"},
	}

	for name, para := range test {
		t.Run(name, func(t *testing.T) {
			got := PathName(para.input1)
			testValue(t, got, para.want)
		})
	}
}

func TestPathStem(t *testing.T) {
	type testPara struct {
		input1 string
		want   string
	}

	test := map[string]testPara{
		"1": {`file/name1`, "name1"},
		"2": {`file/name2.2`, "name2"},
		"3": {`/file\name3.3`, "name3"},
	}

	for name, para := range test {
		t.Run(name, func(t *testing.T) {
			got := PathStem(para.input1)
			testValue(t, got, para.want)
		})
	}
}

func TestIsInFilter(t *testing.T) {
	type testPara struct {
		input1 string
		input2 []string
		want   bool
	}

	test := map[string]testPara{
		"1":  {`11.ini`, []string{`*.ini`}, true},
		"2":  {`11.in`, []string{`*.ini`}, false},
		"3":  {`11.abc`, []string{`11.ini`}, false},
		"4":  {`11.ini`, []string{`11.ini`}, true},
		"5":  {`1.ini`, []string{`11.ini`}, false},
		"6":  {`1.ini`, []string{`11.ini`, `*.ini`}, true},
		"7":  {`1.ini`, []string{`11.ini`, `1.ini`}, true},
		"8":  {`.ini`, []string{`*.ini`}, true},
		"9":  {`ini`, []string{`*.ini`}, false},
		"10": {`1.ini`, []string{`*.INI`}, false},
	}

	for name, para := range test {
		t.Run(name, func(t *testing.T) {
			got := IsInFilter(para.input1, para.input2)
			testValue(t, got, para.want)
		})
	}
}
