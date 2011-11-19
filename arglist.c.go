package gojvm
//#cgo LDFLAGS:-ljvm	-L/usr/lib/jvm/java-6-sun/jre/lib/amd64/server/
//#include "helpers.h"
import "C"
import (
	"errors"
	"fmt"

	"strconv"
	"unsafe"
)

// Implements the *C.jvalue pointer required
// for passing arguments _into_ java.
type ArgList []C.jvalue

func (self *ArgList) Ptr() unsafe.Pointer {
	if self == nil || len(*self) == 0 {
		//fmt.Printf("Returning a nil ptr!\n")
		return nil
	}
	return unsafe.Pointer(&((*self)[0]))
}

func newArgList(ctx *Environment, params ...interface{}) (alp ArgList, err error) {
	//fmt.Printf("newArgList: %+v\n", params...)
	alp = make(ArgList, 0)
	for i, param := range params {
		var ok C.int
		switch v := param.(type) {
		case int:
			alp = append(alp, C.intValue(C.jint(v)))
		case int64:
			alp = append(alp, C.longValue(C.jlong(v)))
		case C.jstring:
			alp = append(alp, C.objValue(v))
		case C.jboolean:
			alp = append(alp, C.boolValue(v))
		case C.jint:
			alp = append(alp, C.intValue(v))
		case C.jobject:
			alp = append(alp, C.objValue(v))
		case *Object:
			alp = append(alp, C.objValue(v.object))
		case *Class:
			alp = append(alp, C.objValue(v.class))
		case C.jvalue:
			alp = append(alp, v)
		case string:
			var str *Object
			str, err = ctx.NewStringObject(v)
			if err == nil {
				alp = append(alp, C.objValue(str.object))
			}
		case []byte:
			var obj *Object
			obj, err = ctx.newByteObject(v)
			if err == nil {
				alp = append(alp, C.objValue(obj.object))
			}
		case bool:
			val := C.jboolean(C.JNI_FALSE)
			if v {
				val = C.JNI_TRUE
			}
			alp = append(alp, C.boolValue(val))
		default:
			err = errors.New(fmt.Sprintf("Unknown type: %T", v))
		}
		if ok != 0 {
			err = errors.New("Couldn't parse arg #" + strconv.Itoa(i+1))
		}
		if err != nil {
			break
		}
	}
	return
}

type Value struct {
	val C.jvalue
}
/*
func ValueOf(i interface{})(v Value, err os.Error){
	switch i.(type) {
		case bool:
			val := C.JNI_FALSE
			if i == true { val = C.JNI_TRUE }
			v.val.z = C.jboolean(val)
	}
}*/
