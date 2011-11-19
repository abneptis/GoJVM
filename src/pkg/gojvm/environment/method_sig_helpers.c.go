package environment

//#cgo CFLAGS:-I../include/
//#cgo LDFLAGS:-ljvm	-L/usr/lib/jvm/java-6-sun/jre/lib/amd64/server/
//#include "helpers.h"
import "C"
import (
	"errors"
	"fmt"
	"reflect"
	"gojvm/types"
)

func TypeOf(env *Environment, v interface{}) (k types.Typed, err error) {
	if kind, ok := v.(types.Typed); ok {
		return kind, nil
	}
	switch vt := v.(type) {
	case C.jstring:
		return types.Class{"java/lang/String"}, nil
	case types.Typed:
		return vt, nil
	case *Object:
		var klass *Class
		var name types.ClassName
		klass, err = vt.ObjectClass()
		if err == nil {
			name, err = klass.Name()
		}
		if err == nil {
			k = types.Class{name.AsPath()}
		}
		return
	}
	vtype := reflect.TypeOf(v)
	vkind := vtype.Kind()
	switch vkind {
	case reflect.Ptr:
		k, err = TypeOf(env, reflect.Indirect(reflect.ValueOf(v)).Interface())
	case reflect.Bool:
		k = types.Basic(types.BoolKind)
	case reflect.Uint8, reflect.Int8:
		k = types.Basic(types.ByteKind)
	case reflect.Int16, reflect.Uint16:
		k = types.Basic(types.ShortKind)
	case reflect.Int32, reflect.Uint32:
		k = types.Basic(types.IntKind)
	case reflect.Uint64, reflect.Int64:
		k = types.Basic(types.LongKind)
	case reflect.Int, reflect.Uint:
		k = types.Basic(types.IntKind)
	case reflect.Float32:
		k = types.Basic(types.FloatKind)
	case reflect.Float64:
		k = types.Basic(types.DoubleKind)
	case reflect.Struct:
		k = types.Class{"golang/" + vkind.String()}
	case reflect.String:
		k = types.Class{"java/lang/String"}
	case reflect.Slice, reflect.Array:
		sltype := vtype.Elem()
		switch sltype.Kind() {
		case reflect.Uint8:
			k = types.ArrayType{types.Basic(types.ByteKind)}
		case reflect.String:
			k = types.ArrayType{types.Class{"java/lang/String"}}
		default:
			err = errors.New("Unhandled slice type " + sltype.String())
		}
	default:
		switch T := v.(type) {
		default:
			err = errors.New(fmt.Sprintf("Unsure how to TypeOf '%s'/%T", vkind.String(), T))
		}
	}
	return
}
