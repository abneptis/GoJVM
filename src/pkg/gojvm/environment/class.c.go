package environment

//#cgo LDFLAGS:-ljvm	-L/usr/lib/jvm/java-6-sun/jre/lib/amd64/server/
//#include "helpers.h"
import "C"
import (
	"gojvm/types"
	"log"
	"unsafe"
)

/* represents a class (object) */
type Class struct {
	env    *Environment
	_klass types.ClassName
	class  C.jclass
}

func newClass(env *Environment, name types.ClassName, class C.jclass) *Class {
	return &Class{env, name, class}
}

/*
	returns the (potentially cached) types.ClassName of the class.
*/
func (self *Class) Name() (name types.ClassName, err error) {
	if len(self._klass) == 0 {
		log.Printf("ClassName(miss)")
		var cstr string
		cstr, _, err = self.CallString(false, "getName")
		if err == nil {
			self._klass = types.NewClassName(cstr)
		}
	}
	if err == nil {
		name = self._klass
	}
	return
}

// Calls the named void-method on the class
func (self *Class) CallVoid(static bool, mname string, params ...interface{}) (err error) {
	meth, args, localStack, err := self.env.getClassMethod(self, static, mname, types.Basic(types.VoidKind), params...)
	if err != nil { return }
	defer blowStack(self.env, localStack)
	C.envCallVoidMethodA(self.env.env, self.class, meth.method, args.Ptr())
	if self.env.ExceptionCheck() {
		err = self.env.ExceptionOccurred()
	}
	return

}

// Calls the named int-method on the class
func (self *Class) CallInt(static bool, mname string, params ...interface{}) (i int, err error) {
	meth, args, localStack, err := self.env.getClassMethod(self, static, mname, types.Basic(types.IntKind), params...)
	if err != nil {	return	}
	defer blowStack(self.env, localStack)
	ji := C.envCallIntMethodA(self.env.env, self.class, meth.method, args.Ptr())
	if self.env.ExceptionCheck() {
		err = self.env.ExceptionOccurred()
	}
	if err == nil {
		i = int(ji)
	}
	return
}

// Calls the named obj-method on the class
func (self *Class) CallObj(static bool, mname string, rval types.Typed, params ...interface{}) (vObj *Object, err error) {
	meth, alp, localStack, err := self.env.getClassMethod(self, static, mname, rval, params...)
	if err != nil {	return }
	var oval C.jobject
	defer blowStack(self.env, localStack)
	if static {
		oval = C.envCallStaticObjectMethodA(self.env.env, self.class, meth.method, alp.Ptr())
	} else {
		oval = C.envCallObjectMethodA(self.env.env, self.class, meth.method, alp.Ptr())
	}
	if oval == nil {
		// is this always the case? or should we use exception check?		
		err = self.env.ExceptionOccurred()
	}
	if err == nil {
		vObj = newObject(self.env, nil, oval)
	}
	return
}

/*
	A wrapper around ObjCallObj specific to java/lang/String, that will return the result as a GoString 

	A null string returned with no exception can be identified in the wasNull return type.
*/
func (self *Class) CallString(static bool, mname string, params ...interface{}) (str string, wasNull bool, err error) {
	strobj, err := self.CallObj(static, mname, types.Class{"java/lang/String"}, params...)
	if err != nil {
		return
	}
	if strobj == nil {
		wasNull = true
		if self.env.ExceptionCheck() {
			err = self.env.ExceptionOccurred()
		} // if no exception, they returned a null string
		return
	}
	defer self.env.DeleteLocalRef(strobj)

	bytesObj, err := strobj.CallObj(false, "getBytes", types.ArrayType{types.Basic(types.ByteKind)}, self.env.utf8())
	if err != nil {
		return
	}
	if bytesObj == nil {
		return // they returned an empty string
	}
	defer self.env.DeleteLocalRef(bytesObj)
	//print("getting array length\n")
	alen := C.envGetArrayLength(self.env.env, bytesObj.object)
	//print("got array length\t",alen, "\n")
	_false := C.jboolean(C.JNI_FALSE)
	//print("getting bytes...\n")
	ptr := C.envGetByteArrayElements(self.env.env, bytesObj.object, &_false)
	//print("setting deferral...\n")
	defer C.envReleaseByteArrayElements(self.env.env, bytesObj.object, ptr, 0)
	//print("going ", alen, " bytes...\n")
	str = string(C.GoBytes(unsafe.Pointer(ptr), C.int(alen)))
	//print("str == ", str, "\n")
	return
}
