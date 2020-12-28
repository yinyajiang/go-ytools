package ylog

import (
	"fmt"
	"io"
	"log"
	"os"

	js "github.com/bitly/go-simplejson"
	yerror "github.com/yinyajiang/go-ytools/error"
)

//Log ...
type Log struct {
	log.Logger
}

//New ...
func New() *Log {
	return newLog("")
}

//NewWithFile ...
func NewWithFile(file string) *Log {
	return newLog(file)
}

func newLog(file string) *Log {
	l := new(Log)
	var w io.Writer
	if len(file) > 0 {
		f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Printf("Log create fail:%v\n", err)
			w = os.Stdout
		} else {
			w = io.MultiWriter(os.Stdout, f)
		}
	} else {
		w = os.Stdout
	}
	l.SetOutput(w)
	l.SetFlags(log.LstdFlags)
	return l
}

//DbgPrintf ...
func (p *Log) DbgPrintf(format string, v ...interface{}) {
	p.print(odbg, format, v...)
}

//DbgPrint ...
func (p *Log) DbgPrint(v ...interface{}) {
	p.print(odbg, "", v...)
}

//TracePrintf ...
func (p *Log) TracePrintf(format string, v ...interface{}) {
	p.print(otrace, format, v...)
}

//TracePrint ...
func (p *Log) TracePrint(v ...interface{}) {
	p.print(otrace, "", v...)
}

//StdPrintf ...
func (p *Log) StdPrintf(format string, v ...interface{}) {
	p.print(ostd, format, v...)
}

//StdPrint ...
func (p *Log) StdPrint(v ...interface{}) {
	p.print(ostd, "", v...)
}

//CodePrint ...
func (p *Log) CodePrint(code int, msg string, data *js.Json) {
	p.uiPrint(code, msg, data, true, "")
}

//ErrPrintWithData ...
func (p *Log) ErrPrintWithData(err error, data *js.Json) {
	p.uiErrPrint(err, data)
}

//ErrPrint ...
func (p *Log) ErrPrint(err error) {
	p.uiErrPrint(err, nil)
}

//ProgressPrint 标准输出进度相关信息
func (p *Log) ProgressPrint(progress float64, speed, size, transffred int, phase string) {
	data := js.New()
	if 100.0 == progress || 0.0 == progress {
		data.Set("progress", fmt.Sprintf("%d", int(progress)))
	} else {
		data.Set("progress", fmt.Sprintf("%.2f", progress))
	}
	data.Set("speed", speed)
	data.Set("size", size)
	data.Set("transffred", transffred)
	data.Set("phase", phase)

	p.uiPrint(1, "", data, false, "")

	if progress == 0.0 {
		p.print(otrace, "", "Progress begin")
	} else if progress == 100.0 {
		p.print(otrace, "", "Progress end,Success")
	}
}

//PrintSuccess 打印成功信息
func (p *Log) PrintSuccess() {
	p.uiErrPrint(nil, nil)
}

func (p *Log) uiPrint(code int, msg string, data *js.Json, record bool, stack string) {
	//for ui
	j := createJSONLog(code, msg, data)
	b, _ := j.MarshalJSON()
	fmt.Println(string(b))

	//for other
	if record {
		if len(stack) > 0 {
			if len(msg) > 0 {
				p.print(otrace, "%s,code(%d),stack(%s)", msg, code, stack)
			} else {
				p.print(otrace, "code(%d),stack(%s)", code, stack)
			}
		} else {
			if len(msg) > 0 {
				p.print(otrace, "%s,code(%d)", msg, code)
			} else {
				p.print(otrace, "code(%d)", code)
			}

		}

	}
}

func (p *Log) uiErrPrint(err error, data *js.Json) {

	if err == nil {
		p.uiPrint(0, yerror.GetCodeTranslate(0), data, true, "")
		return
	}
	stack := true
	e, ok := err.(yerror.Error)
	if !ok {
		e = yerror.NewWithError(-1, err)
		stack = false
	}

	if stack {
		p.uiPrint(e.Code(), e.String(), data, true, e.CallerInfoStr())
	} else {
		p.uiPrint(e.Code(), e.String(), data, true, "")
	}

}

//print ...
func (p *Log) print(lev int, format string, v ...interface{}) {
	p.SetPrefix(otoPrefix(lev))
	s := ""
	if ostd == lev {
		fmt.Println(v...)
		return
	} else if len(format) > 0 {
		s = fmt.Sprintf(format, v...)
	} else {
		s = fmt.Sprint(v...)
	}
	p.Println(":", s)
}
