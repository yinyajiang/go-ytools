package ylog

import (
	"fmt"
	"io"
	"log"
	"os"

	js "github.com/bitly/go-simplejson"
	yerror "github.com/yinyajiang/go-ytools/error"
	tools "github.com/yinyajiang/go-ytools/utils"
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
		f, err := tools.OpenApptendFile(file)
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
func (p *Log) CodePrint(e interface{}, v ...interface{}) {
	var (
		code     int
		msg      string
		data     *js.Json
		callinfo string
	)
	if e == nil {
		code = 0
	} else if c, ok := e.(int); ok {
		code = c
		msg = yerror.GetCodeTranslate(c)
	} else if yerr, ok := e.(yerror.Error); ok {
		code = yerr.Code()
		msg = yerr.Error()
		callinfo = yerr.CallerInfoStr()
	} else if err, ok := e.(error); ok {
		code = -1
		msg = err.Error()
	} else {
		code = -1
		msg = fmt.Sprint(e)
	}

	vt := make([]interface{}, 0, len(v))
	for _, i := range v {
		if j, ok := i.(*js.Json); ok {
			data = j
		} else {
			vt = append(vt, i)
		}
	}
	if len(vt) > 0 {
		if len(msg) > 0 {
			msg += " | "
		}
		msg += fmt.Sprint(vt...)
	}
	p.codePrint(code, msg, data, true, callinfo)

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

	p.codePrint(1, "", data, false, "")

	if progress == 0.0 {
		p.print(otrace, "", "Progress begin")
	} else if progress == 100.0 {
		p.print(otrace, "", "Progress end,Success")
	}
}

func (p *Log) codePrint(code int, msg string, data *js.Json, record bool, stack string) {
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
