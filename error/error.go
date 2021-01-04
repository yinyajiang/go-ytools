package yerror

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	js "github.com/bitly/go-simplejson"
)

var (
	_ Error        = (*_Error)(nil)
	_ fmt.Stringer = (*_Error)(nil)
)

//Error ...
type Error interface {
	error
	fmt.Stringer
	Code() int
	Caller() []CallerInfo
	CallerInfoStr() string
	CallerJSONInfo() *js.Json
	private()
}

type _Error struct {
	err     error
	code    int
	callers []CallerInfo
	wraped  error
}

var _codeTranslate = map[int]string{}

//AddCodeTranslateMap 添加错误码映射关系
func AddCodeTranslateMap(codeMsg map[int]string) {
	for code, msg := range codeMsg {
		_codeTranslate[code] = msg
	}
}

//GetCodeTranslate 获取错误码对应的翻译
func GetCodeTranslate(code int) string {
	msg, ok := _codeTranslate[code]
	if !ok {
		if code == 0 {
			return "Successed"
		}
		return "Unknow"
	}
	return msg
}

//New ...
func New(e interface{}, v ...interface{}) Error {
	if e == nil {
		return nil
	}
	code := -1
	var endmsg string
	if c, ok := e.(int); ok {
		code = c
		if code == 0 {
			return nil
		}
	} else if yerr, ok := e.(Error); ok {
		code = yerr.Code()
		endmsg = yerr.Error()
	} else if err, ok := e.(error); ok {
		code = -1
		endmsg = err.Error()
	} else {
		code = -1
		endmsg = fmt.Sprint(e)
	}

	var msg string
	if len(v) > 0 {
		if fmat, ok := v[0].(string); ok && strings.HasPrefix(fmat, "format!") {
			fmat := fmat[7:]
			v = v[1:]
			msg = fmt.Sprintf(fmat, v...)
		} else {
			msg = fmt.Sprint(v...)
		}
	}
	if len(msg) > 0 {
		msg += "," + endmsg
	} else {
		msg = endmsg
	}
	return &_Error{
		code:    code,
		err:     errors.New(msg),
		callers: CallerList(1),
	}
}

func (p *_Error) Code() int {
	return p.code
}

func (p *_Error) Caller() []CallerInfo {
	return p.callers
}

func (p *_Error) CallerInfoStr() (ret string) {
	first := true
	for _, info := range p.callers {
		if first {
			first = false
		} else {
			ret += " => "
		}
		ret += "file:" + info.File + ",func:" + info.Fun + ",line:" + strconv.Itoa(info.Line)
	}
	return
}

func (p *_Error) CallerJSONInfo() *js.Json {
	b, err := json.Marshal(p.callers)
	if err != nil {
		return nil
	}
	j, err := js.NewJson(b)
	if err != nil {
		return nil
	}
	return j
}

func (p *_Error) Error() string {
	if nil == p.err {
		return GetCodeTranslate(p.code)
	}
	return p.err.Error()
}

func (p *_Error) String() (ret string) {
	ret = p.Error()
	return
}

func (p *_Error) private() {
	panic("unreached")
}
