package gojvm

//#include "helpers.h"
import "C"
import "unsafe"

type Object struct {
	env    *Environment
	_klass *Class
	object C.jobject
}



// returns a new object value with specified parameters
// NB: refs are NOT adjusted directly by this call! Use it as a casting/construction-helper,
// not a Clone()
func newObject(env *Environment, klass *Class, obj C.jobject) *Object {
	return &Object{env, klass, obj}
}

/* 
	Returns the Class() associated with the object;
	If this was known at call, that value will be used,
	else it will be resolved through the environment into a 
	class type.
*/
func (self *Object) ObjectClass() (c *Class, err error) {
	if self._klass == nil {
		self._klass, err = self.env.getObjectClass(self)
	}
	if err == nil {
		c = self._klass
	}
	return
}

/*
	Returns the (potentially cached) name of the ObjectClass of the
	named object.
*/
func (self *Object) ClassName() (name ClassName, err error) {
	c, err := self.ObjectClass()
	if err == nil {
		name, err = c.Name()
	}
	return
}

func (self Object) JavaType() int { return JAVAObject }

// Calls the named void-method on the object instance
func (self *Object) CallVoid(static bool, mname string, params ...interface{}) (err error) {
	meth, args, err := self.env.getObjectMethod(self, static, mname, BasicType(JavaVoidKind), params...)
	if err != nil {
		return
	}
	if static {
		C.envCallStaticVoidMethodA(self.env.env, self.object, meth.method, args.Ptr())
	} else {
		C.envCallVoidMethodA(self.env.env, self.object, meth.method, args.Ptr())
	}
	if self.env.exceptionCheck() {
		err = self.env.exceptionOccurred()
	}
	return
}

// Calls the named int-method on the object instance
func (self *Object) CallInt(static bool, mname string, params ...interface{}) (i int, err error) {
	meth, args, err := self.env.getObjectMethod(self, static, mname, BasicType(JavaIntKind), params...)
	if err != nil {
		return
	}
	var ji C.jint
	if static {
		ji = C.envCallStaticIntMethodA(self.env.env, self.object, meth.method, args.Ptr())
	} else {
		ji = C.envCallIntMethodA(self.env.env, self.object, meth.method, args.Ptr())
	}
	if self.env.exceptionCheck() {
		err = self.env.exceptionOccurred()
	}
	if err == nil {
		i = int(ji)
	}
	return
}

// Calls the named Object-method on the object instance
func (self *Object) CallObj(static bool, mname string, rval JavaType, params ...interface{}) (vObj *Object, err error) {
	meth, alp, err := self.env.getObjectMethod(self, static, mname, rval, params...)
	if err != nil {
		return
	}
	var oval C.jobject
	if static {
		oval = C.envCallStaticObjectMethodA(self.env.env, self.object, meth.method, alp.Ptr())
	} else {
		oval = C.envCallObjectMethodA(self.env.env, self.object, meth.method, alp.Ptr())
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

/* 
	A wrapper around ObjCallObj specific to java/lang/String, that will return the result as a GoString 

	A null string returned with no exception can be differentiated via the wasNull return value.
*/
func (self *Object) CallString(static bool, mname string, params ...interface{}) (str string, wasNull bool, err error) {
	strobj, err := self.CallObj(static, mname, ClassType{"java/lang/String"}, params...)
	if err != nil {
		return
	}
	if strobj == nil {
		wasNull = true
		if self.env.exceptionCheck() {
			err = self.env.exceptionOccurred()
		} // if no exception, they returned a null string
		return
	}
	defer self.env.LocalUnref(strobj)

	bytesObj, err := strobj.CallObj(false, "getBytes", ArrayType{BasicType(JavaByteKind)}, self.env.utf8())
	if err != nil {
		return
	}
	if bytesObj == nil {
		return // they returned an empty string
	}
	defer self.env.LocalUnref(bytesObj)
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
