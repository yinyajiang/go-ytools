package yerror

const (
	//ECSuccess ...
	ECSuccess = iota
	_
	_
	_
	//ECNotAllowedInLocale ...
	ECNotAllowedInLocale
	_
	//ECAsynUnkonw ...
	ECAsynUnkonw
	//ECTimeOut ...
	ECTimeOut
	//ECHasBeenCancle ...
	ECHasBeenCancle
	//ECArgError ...
	ECArgError
)

var (
	codeTranslate = map[int]string{
		ECSuccess:            "Successed",
		ECNotAllowedInLocale: "Not allowed in locale",
		ECAsynUnkonw:         "Asyn unkown",
		ECTimeOut:            "Time out",
		ECHasBeenCancle:      "Has been cancle",
		ECArgError:           "Arg error",
	}
)
