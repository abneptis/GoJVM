package gojvm

//#cgo LDFLAGS:-ljvm	-L/usr/lib/jvm/java-6-sun/jre/lib/amd64/server/
//#include "helpers.h"
import "C"
import (
	"os"
	"log"
	"unsafe"
)
type Class	struct {
	env		*Environment
	_klass	ClassName
	class	C.jclass
}

func newClass(env *Environment, name ClassName, class C.jclass)(*Class){
	return &Class{env,name,class}
}


func (self *Class)Name()(name ClassName, err os.Error){
	if len(self._klass) == 0 {
		log.Printf("ClassName(miss)")
		var cstr string
		cstr, _, err = self.CallString(false, "getName")
		if err == nil {
			self._klass = NewClassName(cstr)
		}
	}
	if err == nil {
		name = self._klass
	}
	return
}

func (self *Class)CallVoid(static bool, mname string, params ...interface{})(err os.Error){
	meth, args, err := self.env.getClassMethod(self, static,   mname,  BasicType(JavaVoidKind), params...)
	if err != nil { return }
	C.envCallVoidMethodA(self.env.env, self.class, meth.method, args.Ptr())
	if self.env.exceptionCheck() {
		err = self.env.exceptionOccurred()
	}
	return
	
}


func (self *Class)CallInt(static bool, mname string, params ...interface{})(i int, err os.Error){
	meth, args, err := self.env.getClassMethod(self, static,   mname,  BasicType(JavaIntKind), params...)
	if err != nil { return }
	ji := C.envCallIntMethodA(self.env.env, self.class, meth.method, args.Ptr())
	if self.env.exceptionCheck() {
		err = self.env.exceptionOccurred()
	}
	if err == nil { i = int(ji) }
	return
}


func (self *Class)CallObj(static bool, mname string, rval JavaType, params ...interface{})(vObj *Object, err os.Error){
	meth, alp, err := self.env.getClassMethod(self, static, mname, rval, params...)
	if err != nil { return }
	var oval C.jobject
	if static {
		oval = C.envCallStaticObjectMethodA(self.env.env, self.class, meth.method, alp.Ptr())
	} else {
		oval = C.envCallObjectMethodA(self.env.env, self.class, meth.method, alp.Ptr())
	}
	if oval == nil {
		// is this always the case? or should we use exception check?		
		err = self.env.exceptionOccurred()
	}
	if err == nil {
		vObj = newObject(self.env, nil, oval)
	}
	return
}

/* A wrapper around ObjCallObj specific to java/lang/String, that will return the result as a GoString 

	A null string returned with no exception can be identified in the wasNull return type.
*/
func (self *Class)CallString(static bool, mname string,params ...interface{})(str string, wasNull bool, err os.Error){
	strobj, err := self.CallObj(static, mname, ClassType{"java/lang/String"}, params...)
	if err != nil { return }
	if strobj == nil {
		wasNull = true
		if self.env.exceptionCheck() {
			err = self.env.exceptionOccurred()
		}// if no exception, they returned a null string
		return
	}
	defer self.env.LocalUnref(strobj)
	
	bytesObj, err := strobj.CallObj( false, "getBytes", ArrayType{BasicType(JavaByteKind)}, self.env.utf8())
	if err != nil { return }
	if bytesObj == nil {
		return	// they returned an empty string
	}
	defer self.env.LocalUnref(bytesObj)
	//print("getting array length\n")
	alen := C.envGetArrayLength(self.env.env,	bytesObj.object)
	//print("got array length\t",alen, "\n")
	_false := C.jboolean(C.JNI_FALSE)
	//print("getting bytes...\n")
	ptr :=  C.envGetByteArrayElements(self.env.env, bytesObj.object, &_false)
	//print("setting deferral...\n")
	defer C.envReleaseByteArrayElements(self.env.env, bytesObj.object, ptr, 0)
	//print("going ", alen, " bytes...\n")
	str = string(C.GoBytes(unsafe.Pointer(ptr), C.int(alen)))
	//print("str == ", str, "\n")
	return
}
	

