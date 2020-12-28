package yerror

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

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
	Wraped() error
	WrapedList() []error
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
		msg, ok = codeTranslate[code]
		if !ok {
			if code == 0 {
				return "Successed"
			}
			return "Unknow"
		}
	}
	return msg
}

//New ...
func New(code int) Error {
	if 0 == code {
		return nil
	}
	return &_Error{
		code:    code,
		callers: CallerList(1),
	}
}

//NewWithMsg ...
func NewWithMsg(code int, msg string) Error {
	if 0 == code {
		return nil
	}
	return &_Error{
		code:    code,
		err:     errors.New(msg),
		callers: CallerList(1),
	}
}

//NewWithError ...
func NewWithError(code int, err error) Error {
	if 0 == code {
		return nil
	}
	if nil == err {
		return New(code)
	}
	return &_Error{
		code:    code,
		err:     err,
		callers: CallerList(1),
	}

}

//NewWithErrorMsg ...
func NewWithErrorMsg(code int, err error, msg string) Error {
	if 0 == code {
		return nil
	}
	if nil == err {
		return New(code)
	}
	return &_Error{
		code:    code,
		err:     errors.New(msg + " | " + err.Error()),
		callers: CallerList(1),
	}

}

//NewF ...
func NewF(code int, format string, v ...interface{}) Error {
	if 0 == code {
		return nil
	}
	return &_Error{
		code:    code,
		err:     fmt.Errorf(format, v...),
		callers: CallerList(1),
	}
}

//Wrap ...
func Wrap(werr error, code int) Error {
	if 0 == code {
		return nil
	}
	return &_Error{
		code:    code,
		callers: CallerList(1),
		wraped:  werr,
	}
}

//WrapWithMsg ...
func WrapWithMsg(werr error, code int, msg string) Error {
	if 0 == code {
		return nil
	}
	return &_Error{
		code:    code,
		err:     errors.New(msg),
		callers: CallerList(1),
		wraped:  werr,
	}
}

//WrapWithError ...
func WrapWithError(werr error, code int, err error) Error {
	if 0 == code {
		return nil
	}
	if nil == err {
		return Wrap(werr, code)
	}
	return &_Error{
		code:    code,
		err:     err,
		callers: CallerList(1),
		wraped:  werr,
	}
}

//WrapF ...
func WrapF(werr error, code int, format string, v ...interface{}) Error {
	if 0 == code {
		return nil
	}
	return &_Error{
		code:    code,
		err:     fmt.Errorf(format, v...),
		callers: CallerList(1),
		wraped:  werr,
	}
}

func (p *_Error) Code() int {
	return p.code
}

func (p *_Error) Wraped() error {
	return p.wraped
}

func (p *_Error) WrapedList() (ret []error) {
	last := p
	for last != nil && nil != last.wraped {
		ret = append(ret, last.wraped)
		ok := false
		last, ok = last.wraped.(*_Error)
		if !ok {
			break
		}
	}
	return
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
	for _, w := range p.WrapedList() {
		ret += " <- " + w.Error()
	}
	return
}

func (p *_Error) private() {
	panic("unreached")
}
