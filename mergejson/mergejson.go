package mergejson

import (
	"io"
	"reflect"

	js "github.com/bitly/go-simplejson"
)

func isContain(s []interface{}, f interface{}) bool {
	for _, val := range s {
		if reflect.DeepEqual(val, f) {
			return true
		}
	}
	return false
}

func mergeJSONArray(toarray []interface{}, fromarray []interface{}) []interface{} {
	for _, fval := range fromarray {
		if !isContain(toarray, fval) {
			toarray = append(toarray, fval)
		}
	}
	return toarray
}

func mergeJSONMap(tomap map[string]interface{}, frommap map[string]interface{}) {
	for fkey, fval := range frommap {
		tval, ok := tomap[fkey]
		if !ok {
			tomap[fkey] = fval
			continue
		}

		typeof := reflect.TypeOf(fval)
		switch typeof.Kind() {
		case reflect.Slice:
			tomap[fkey] = mergeJSONArray(tval.([]interface{}), fval.([]interface{}))
		case reflect.Map:
			mergeJSONMap(tval.(map[string]interface{}), fval.(map[string]interface{}))
		default:
			tomap[fkey] = fval
		}
	}
}

//MergeJSON from合并json到to,以from为准
func MergeJSON(to *js.Json, from *js.Json) {
	totype := reflect.TypeOf(to.Interface())
	fromtype := reflect.TypeOf(from.Interface())
	if totype.Kind() != fromtype.Kind() {
		return
	}

	switch totype.Kind() {
	case reflect.Slice:
		mergeJSONArray(to.MustArray(), from.MustArray())
	case reflect.Map:
		mergeJSONMap(to.MustMap(), from.MustMap())
	}

}

//MergeFileJSON from合并json到to,以from为准
func MergeFileJSON(to io.Reader, from io.Reader, result io.Writer) error {
	fromjson, err := js.NewFromReader(from)
	if err != nil {
		return err
	}
	tojson, err := js.NewFromReader(to)
	if err != nil {
		return err
	}

	MergeJSON(tojson, fromjson)
	data, err := tojson.MarshalJSON()
	if err != nil {
		return err
	}
	_, err = result.Write(data)
	if err != nil {
		return err
	}
	return nil
}
