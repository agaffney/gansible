package python

import (
	"fmt"
	"github.com/sbinet/go-python"
	"reflect"
)

func String(val string) *python.PyObject {
	return python.PyString_FromString(val)
}

func Int(val int) *python.PyObject {
	return python.PyInt_FromLong(val)
}

func Bool(val bool) *python.PyObject {
	v := 0
	if val {
		v = 1
	}
	return python.PyBool_FromLong(v)
}

func Tuple(items []interface{}) *python.PyObject {
	itemsLen := len(items)
	if itemsLen == 0 {
		return python.PyTuple_New(0)
	} else {
		tmpItems := make([]*python.PyObject, 0)
		for _, val := range items {
			tmpItems = append(tmpItems, ToPythonType(val))
		}
		return python.PyTuple_Pack(len(items), tmpItems...)
	}
}

func List(items []interface{}) *python.PyObject {
	pyList := python.PyList_New(0)
	for _, val := range items {
		python.PyList_Append(pyList, ToPythonType(val))
	}
	return pyList
}

func Dict(items map[interface{}]interface{}) *python.PyObject {
	pyDict := python.PyDict_New()
	for key, val := range items {
		python.PyDict_SetItem(pyDict, ToPythonType(key), ToPythonType(val))
	}
	return pyDict
}

func getPythonType(val *python.PyObject) string {
	pyType := val.Type().String()
	// extract from: <type 'foo'>
	return pyType[7 : len(pyType)-2]
}

func ToPythonType(val interface{}) *python.PyObject {
	switch v := reflect.ValueOf(val); v.Kind() {
	case reflect.Bool:
		return Bool(v.Bool())
	// XXX: break some of these out into PyLong?
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Int(int(v.Int()))
	case reflect.String:
		return String(v.String())
	case reflect.Map:
		return Dict(val.(map[interface{}]interface{}))
	case reflect.Array, reflect.Slice:
		return List(val.([]interface{}))
	default:
		return nil
	}
}

func ToGoType(val *python.PyObject) interface{} {
	pyType := getPythonType(val)
	switch pyType {
	case "list":
		ret := make([]interface{}, 0)
		for idx := 0; idx < python.PyList_Size(val); idx++ {
			ret = append(ret, ToGoType(python.PyList_GetItem(val, idx)))
		}
		return ret
	case "str":
		return python.PyString_AsString(val)
	default:
		fmt.Printf("unsupported type: %s\n", pyType)
		return nil
	}
}
