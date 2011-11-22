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
		return types.Class{types.JavaLangString}, nil
	case types.Typed:
		return vt, nil
	case *Object:
		return types.Class{types.JavaLangObject}, nil
	}
	k, err = reflectedType(env, v)
	return
}

func reflectedType(env *Environment, v interface{})(k types.Typed, err error){
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
		k = types.Class{[]string{"golang",vkind.String()}}
	case reflect.String:
		k = types.Class{types.JavaLangString}
	case reflect.Slice, reflect.Array:
		sltype := vtype.Elem()
		switch sltype.Kind() {
		case reflect.Uint8:
			k = types.ArrayType{types.Basic(types.ByteKind)}
		case reflect.String:
			k = types.ArrayType{types.Class{types.JavaLangString}}
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

type errInterface interface{
	Error()(string)
}

func ReflectedSignature(ctx *Environment, f interface{})(sig types.MethodSignature, err error){
	sig.Return = types.Basic(types.UnspecKind)
	
	ftype :=  reflect.TypeOf(f)
	sig.Params = make([]types.Typed, 0)
	if ftype.Kind() != reflect.Func {
		err = errors.New("ReflectedSignature: f is not a function")
	}
	if err == nil && ftype.NumIn() < 2 {
		err = errors.New("ReflectedSignature: f is not a callback (insufficient args)")
	}
	if err == nil &&  ftype.NumOut() > 1 {
		err = errors.New("ReflectedSignature: f is not a callback (too many returns)")
	}
	if err == nil && ftype.In(0) != reflect.TypeOf(&Environment{}) {
		err = errors.New("bad first-arg Type: must be *Environment")
	}
	if err == nil && ftype.In(1) != reflect.TypeOf(&Object{}) {
		err = errors.New("bad second-arg Type: must be *Object")
	}
	if err != nil { return }
	for i := 2; i < ftype.NumIn(); i++ {
		var k types.Typed
		if ftype.In(i).Kind() == reflect.Ptr {
			if ftype.In(i) == reflect.TypeOf(&Object{}){
				k = types.Class{types.JavaLangObject}
			} else {
				itype := ftype.In(i).Elem()
				pobj := reflect.New(itype).Interface()
				k, err = TypeOf(ctx, pobj)
			}
		} else {
			k, err = TypeOf(ctx, reflect.New(ftype.In(i)).Interface())
		}
		if err != nil { break }
		sig.Params = append(sig.Params, k)
	}
	if ftype.NumOut() == 1 {
		var k types.Typed
		k, err = reflectedType(ctx, reflect.New(ftype.Out(0)).Interface())
		if err == nil {
			sig.Return = k
		}
	} else {
		sig.Return = types.Basic(types.VoidKind)
	}
	return	
}

