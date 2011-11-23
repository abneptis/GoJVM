package gojvm
//#cgo CFLAGS:-I../include/
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
type argList []C.jvalue


/// Returns an unsafe.Pointer to the arglist, suitable 
/// for passing to the JNI xxxA() variadic methods.
/// The pointer is only valid at the time it is returned,
/// and changes to the underlying arglist could invalidate
/// this list.
///
/// (If you're constructing w/o newArgList, be sure to use make(),
/// in order to ensure the value references are aligned)
func (self *argList) Ptr() unsafe.Pointer {
	if self == nil || len(*self) == 0 {
		//fmt.Printf("Returning a nil ptr!\n")
		return nil
	}
	return unsafe.Pointer(&((*self)[0]))
}

/* dereferences objects in a list, useful for deferrals */
func blowStack(env *Environment, objs []*Object){
	for _, obj := range(objs){
		env.DeleteLocalRef(obj)
	}
}


//	TODO(refcounting): any constructed objects will be leaked on call return, 
//	as nothing cleans up proxy objects.  I'm also torn on how to differentiate
//	the objects made here  and those coming in from other references.
//
//	Refcounting attempt 1, objects _we_ construct will be returned in the objStack,
//	otherwise refcounts of 'pass-through' java natives are untouched by the call to newArgList.
//
//	On error, the stack has already been blown (and will be empty).
func newArgList(ctx *Environment, params ...interface{}) (alp argList, objStack []*Object, err error) {
	alp = make(argList, 0)
	defer func(){
		if err != nil {
			blowStack(ctx, objStack)
			objStack = []*Object{}
		}
	}()
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
				objStack = append(objStack, str)
				alp = append(alp, C.objValue(str.object))
			}
		case []string:
			var klass *Class
			var	obj	  *Object
		 	klass, err = ctx.GetClassStr("java/lang/String")
		 	// classes via this channel are cached and globally referenced by gojvm, not stacked.
			if err == nil {
				obj, err = ctx.newObjectArray(len(v), klass, nil)
			}
			if err == nil {
				objStack = append(objStack, obj)
				for i, s := range(v){
					var str *Object
					str, err = ctx.NewStringObject(s)
					if err == nil {
		  				// I'm assuming stuffing the array adds a reference inside the JVM.
		  				defer ctx.DeleteLocalRef(str)
	  					ctx.setObjectArrayElement(obj, i, str)
	  					if ctx.ExceptionCheck(){
	  						err = ctx.ExceptionOccurred()
	  					}
	  					
	  				}
	  				if err != nil {
	  					break
	  				}
				}
			}
			if err == nil {
				alp = append(alp, C.objValue(obj.object))
			}
		case []byte:
			var obj *Object
			obj, err = ctx.newByteObject(v)
			if err == nil {
				alp = append(alp, C.objValue(obj.object))
			}
			objStack = append(objStack, obj)
		case bool:
			val := C.jboolean(C.JNI_FALSE)
			if v {
				val = C.JNI_TRUE
			}
			alp = append(alp, C.boolValue(val))
		default:
			err = errors.New(fmt.Sprintf("Unknown type: %T/%s", v, v))
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

// Essentially, the java generic.  Type information is NOT carried
// with the value, however is required for proper use.  (don't use
// unless you know the distinction between jobject, jclass and jvalue).
type Value struct {
	val C.jvalue
}

