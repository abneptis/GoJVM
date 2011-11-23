package gojvm

//#include "helpers.h"
import "C"
import (
	"gojvm/types"
	"unsafe"
	"reflect"
	"sync"
	"errors"
)

// due to the design of the JNI callback method, in order to support
// callbacks (safely), we need to be able to map from C.JVM and C.JNIEnv
// into their local types.
// the AllVMs & AllEnvs structures handles this via a go (RW)mutex.
type vmPtrMap struct {
	jvms 	map[uintptr]*JVM
	maplock	*sync.RWMutex
}

type envPtrMap struct {
	envs 	map[uintptr]*Environment
	maplock	*sync.RWMutex
}

var AllVMs	= &vmPtrMap{ map[uintptr]*JVM{}, &sync.RWMutex{}}
var AllEnvs	= &envPtrMap{ map[uintptr]*Environment{}, &sync.RWMutex{}}

func (self *vmPtrMap)Find(ptr uintptr)(*JVM){
	self.maplock.RLock()
	defer self.maplock.RUnlock()
	return self.jvms[ptr]
}

func (self *vmPtrMap)Add(vm *JVM){
	self.maplock.Lock()
	defer self.maplock.Unlock()
	self.jvms[uintptr(unsafe.Pointer(vm.jvm))] = vm 
}


func (self *envPtrMap)Find(ptr uintptr)(*Environment){
	self.maplock.RLock()
	defer self.maplock.RUnlock()
	return self.envs[ptr]
}

func (self *envPtrMap)Add(env *Environment){
	self.maplock.Lock()
	defer self.maplock.Unlock()
	self.envs[uintptr(unsafe.Pointer(env.env))] = env 
}


//export goCallbackNArgs
func goCallbackNArgs(envp, obj uintptr, fId int)(int){
	env := AllEnvs.Find(envp)
	if env == nil {
		return -1
	}
	if env.jvm == nil { return -1 }
	var cd callbackDescriptor
	var ok bool
	if cd, ok = env.jvm.registered[fId]; !ok { return -1 }
	return len(cd.Signature.Params)
}

// C callbacks actually start in the 'generifiedX' calls (see the .c files)
// Next, goCallbackNArgs is used to determine the number of variadic paramers to expect (
// java doesn't inform us of this, and we can't force any of the callback parameters.
//
// Finally, goCallback looks up the 'fId' - our internal function reference ID (the X in generified),
// un(re) marshalls all the parameters appropriately, calls our function, and returns
// any underlying value back to generified who will return it to the JVM.  
// The initial jbool indicates success, and any failure should check for exceptions.
//
//export goCallback
func goCallback(envp, obj uintptr, fId int, nargs int, argp uintptr)(ok C.jboolean, val interface{}){
	args := C.ArgListPtr(unsafe.Pointer(argp))
	env := AllEnvs.Find(envp)
	if env == nil {
		panic("Got a nil environment")
	}
	if env.jvm == nil {
		panic("Environment pointer has no JVM")
	}
	var cd callbackDescriptor
	var _ok bool
	if cd, _ok = env.jvm.registered[fId]; !_ok {
		print("Unknown callbackId: \t", fId,"\n")
		return
	}
	// TODO: pack argp somehow...
	if nargs != len(cd.Signature.Params) {
		panic("callback/signature length mismatch")
	}
	inCall := []reflect.Value{ reflect.ValueOf(env), reflect.ValueOf(newObject(C.jobject(unsafe.Pointer(obj))))}
	var err error
	for i := 0; i < nargs; i ++ {
		switch cd.Signature.Params[i].Kind() {
			case types.BoolKind: inCall = append(inCall, reflect.ValueOf(bool(0 != C.valBool(C.getArg(args, C.int(i)))) ))
			case types.LongKind: inCall = append(inCall, reflect.ValueOf(int64(C.valLong(C.getArg(args, C.int(i)))) ))
			case types.IntKind: inCall = append(inCall, reflect.ValueOf(int(C.valInt(C.getArg(args, C.int(i)))) ))
			case types.ShortKind: inCall = append(inCall, reflect.ValueOf(int16(C.valShort(C.getArg(args, C.int(i)))) ))
			case types.FloatKind: inCall = append(inCall, reflect.ValueOf(float32(C.valFloat(C.getArg(args, C.int(i)))) ))
			case types.DoubleKind: inCall = append(inCall, reflect.ValueOf(float64(C.valDouble(C.getArg(args, C.int(i)))) ))
			case types.ClassKind: inCall = append(inCall, reflect.ValueOf(newObject(C.valObject(C.getArg(args, C.int(i))) )))
			default:
				err = errors.New("Couldn't reflect kind " + cd.Signature.Params[i].Kind().TypeString())
		}
		if err != nil { break }
	}
	if err != nil { return }
	outCall := reflect.ValueOf(cd.F).Call(inCall)
	switch cd.Signature.Return.Kind() {
		case types.VoidKind:
			return	1,nil
	}
	switch cd.Signature.Return.Kind() {
		case types.BoolKind:
			if outCall[0].Interface().(bool) {
				return 1, C.jboolean(C.JNI_TRUE)
			} else {
				return 1, C.jboolean(C.JNI_FALSE)
			}
		case types.ByteKind: return 1, C.jbyte(outCall[0].Interface().(byte))
		case types.CharKind: return 1, C.jchar(outCall[0].Interface().(int))
		case types.IntKind: return 1, C.jint(outCall[0].Interface().(int))
		case types.ShortKind: return 1, C.jshort(outCall[0].Interface().(int16))
		case types.LongKind: return 1, C.jint(outCall[0].Interface().(int64))
		case types.FloatKind: return 1, C.jfloat(outCall[0].Interface().(float32))
		case types.DoubleKind: return 1, C.jdouble(outCall[0].Interface().(float64))
		case types.ClassKind:
			klass := cd.Signature.Return.(types.Class).Klass
			if klass.Cmp(types.JavaLangString) == 0 {
				var obj *Object
				str := outCall[0].Interface().(string)
				obj, err = env.NewStringObject(str)
				if err == nil {
					return 1, C.jstring(obj.object)
				} // else, exception occurred
				// not needed as callbacks will reap their own refs.
				// env.DeleteLocalRef(obj)
				print("String Error\t", err.Error())
				return 0, nil
			}
			return 1, C.jobject(outCall[0].Interface().(*Object).object)
		default:	panic("array return type not yet supported")
	}
	
	panic("not reached")
}
