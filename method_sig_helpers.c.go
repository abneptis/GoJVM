package gojvm

//#cgo LDFLAGS:-ljvm	-L/usr/lib/jvm/java-6-sun/jre/lib/amd64/server/
//#include</usr/lib/jvm/java-6-sun-1.6.0.26/include/jni.h>
//#include "helpers.h"
import "C"
import (
	"errors"
	"fmt"
	"reflect"
)

func TypeOf(env *Environment, v interface{}) (k JavaType, err error) {
	if kind, ok := v.(JavaType); ok {
		return kind, nil
	}
	switch vt := v.(type) {
	case C.jstring:
		return ClassType{"java/lang/String"}, nil
	case JavaType:
		return vt, nil
	case *Object:
		var name ClassName
		klass, err := vt.ObjectClass()
		if err == nil {
			name, err = klass.Name()
		}
		if err == nil {
			k = ClassType{name.AsPath()}
		}
		return
	}
	vtype := reflect.TypeOf(v)
	vkind := vtype.Kind()
	switch vkind {
	case reflect.Ptr:
		k, err = TypeOf(env, reflect.Indirect(reflect.ValueOf(v)).Interface())
	case reflect.Bool:
		k = BasicType(JavaBoolKind)
	case reflect.Uint8, reflect.Int8:
		k = BasicType(JavaByteKind)
	case reflect.Int16, reflect.Uint16:
		k = BasicType(JavaShortKind)
	case reflect.Int32, reflect.Uint32:
		k = BasicType(JavaIntKind)
	case reflect.Uint64, reflect.Int64:
		k = BasicType(JavaLongKind)
	case reflect.Int, reflect.Uint:
		k = BasicType(JavaIntKind)
	case reflect.Float32:
		k = BasicType(JavaFloatKind)
	case reflect.Float64:
		k = BasicType(JavaDoubleKind)
	case reflect.Struct:
		k = ClassType{"golang/" + vkind.String()}
	case reflect.String:
		k = ClassType{"java/lang/String"}
	case reflect.Slice, reflect.Array:
		sltype := vtype.Elem()
		switch sltype.Kind() {
		case reflect.Uint8:
			k = ArrayType{BasicType(JavaByteKind)}
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
