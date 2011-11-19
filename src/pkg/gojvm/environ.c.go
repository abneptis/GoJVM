package gojvm

//#cgo LDFLAGS:-ljvm	-L/usr/lib/jvm/java-6-sun/jre/lib/amd64/server/
//#include</usr/lib/jvm/java-6-sun-1.6.0.26/include/jni.h>
//#include <stdlib.h>
//#include <libio.h>
//#include <unistd.h>
//#include "helpers.h"
import "C"
import (
	"errors"
	"unsafe"
)

const (
	JAVAClass = iota
	JAVAObject
)

/* 

	An environment consists of a pointer to a JNI environment
	and a caching map of class names to (globally referenced) class objects.
	
	TODO: Handle references on other items (nominally) correctly.

*/
type Environment struct {
	env             *C.JNIEnv
	classes         map[string]*Class
	quietExceptions bool
	// various 'consts'
	_UTF8 C.jstring // "UTF8" parameter
}

func (self *Environment) getObjectMethod(obj *Object, static bool, mname string, rType JavaType, params ...interface{}) (meth *Method, args ArgList, err error) {
	meth, err = self._objMethod(obj, mname, rType, params...)
	if err != nil {
		return
	}
	args, err = newArgList(self, params...)
	return
}

func (self *Environment) defMute()(func()){
	muted := self.Muted()
	self.Mute(true)
	return func(){
		self.Mute(muted)
	}
}

func (self *Environment) getClassMethod(c *Class, static bool, mname string, rType JavaType, params ...interface{}) (meth *Method, args ArgList, err error) {
	if !static {
		meth, err = self._classMethod(c, mname, rType, params...)
	} else {
		meth, err = self._classStaticMethod(c, mname, rType, params...)
	}
	if err != nil {
		return
	}
	args, err = newArgList(self, params...)
	return
}

// (Un)Suppress the java console barf of exceptions
// (execeptions are still caught, cleared and returned)
func (self *Environment) Mute(mute bool) { self.quietExceptions = mute }

// Returns the current state of the environmental exception mute.
func (self *Environment) Muted()( bool) { return self.quietExceptions }

// Refcounting is probably needed here, TODO: figure that out...
func (self *Environment) utf8() C.jstring {
	if self._UTF8 == nil {
		cs := C.CString("UTF8")
		defer C.free(unsafe.Pointer(cs))
		self._UTF8 = C.envNewStringUTF(self.env, cs)
	}
	return self._UTF8
}

func NewEnvironment() *Environment {
	return &Environment{
		env:     new(C.JNIEnv),
		classes: map[string]*Class{},
	}
}

func (self Class) JavaType() int { return JAVAClass }

/* represents JNI method call;  without subject, style & parameters,
it is useless.  It (appears) to be an error to ref/unref methods.
*/ 
type Method struct {
	method C.jmethodID
}

func (self *Environment) findCachedClass(klass ClassName) (c *Class, err error) {
	if class, ok := self.classes[klass.AsPath()]; ok {
		c = class
	} else {
		err = ErrUnknownClass // not technically an exception, but shouldn't bubble either.
	}
	return
}

/* 
	returns a new *Object of class 'java/lang/String', containing the (UTF16 reinterpreted)
	representation of 's'.  Mostly a helper for passing strings into Java.
*/
func (self *Environment) NewStringObject(s string) (obj *Object, err error) {
	if err == nil {
		obj, err = self.NewInstanceStr("java/lang/String", []byte(s), self.utf8())
	}
	/*
		The naieve approach doesn't work w/ complex or bad UTF8
		    obj = &Object{
		    	object: C.jobject(C.envNewStringUTF(self.env,cs)),
		    }
	*/
	return
}


func (self *Environment) setObjectArrayElement(arr *Object, pos int, item *Object)(err error) {
	C.envSetObjectArrayElement(self.env, arr.object, C.jint(pos), item.object)
	return
}

func (self *Environment) newObjectArray(sz int, klass *Class, init C.jobject) (o *Object, err error) {
	ja := C.envNewObjectArray(self.env, C.jint(sz), klass.class, init)
	if ja == nil {
		err = self.exceptionOccurred()
	}
	if err == nil {
		o = newObject(self, klass, C.jobject(ja))
	}
	return
}

func (self *Environment) newByteObject(bts []byte) (o *Object, err error) {
	ja := C.envNewByteArray(self.env, C.jint(len(bts)))
	if ja == nil {
		err = errors.New("Error allocating byte array")
	}
	if err == nil && len(bts) > 0 {
		bptr := make([]byte, len(bts))
		copy(bptr, bts)
		//log.Printf("bptr: %s %p %p", bptr,bptr, &bptr[0] )
		C.envSetByteArrayRegion(self.env, ja, 0, C.jint(len(bptr)), unsafe.Pointer(&bptr[0]))
	}
	if err == nil {
		o = newObject(self, nil, C.jobject(ja))
	}
	return
}

/* 
	returns a new *Object of the class named by 'klass' (Wrapper around NewInstance(NewClassName(...)))
*/
func (self *Environment) NewInstanceStr(klass string, params ...interface{}) (obj *Object, err error) {
	class, err := self.GetClass(NewClassName(klass))
	if err != nil {
		return
	}
	return self.NewInstance(class, params...)
}

/*
	returns a new *Object of type *Class, using the constructor identified by []params
*/
func (self *Environment) NewInstance(c *Class, params ...interface{}) (o *Object, err error) {
	meth, alp, err := self.getClassMethod(c, false, "<init>", BasicType(JavaVoidKind), params...)
	//	meth, alp, err := self.getObjectMethod(newObject(self, c, C.jobject( c.class)), "<init>", BasicType(JavaVoidKind), params...)
	if err != nil {
		return
	}
	obj := C.envNewObjectA(self.env, c.class, meth.method, alp.Ptr())
	if obj != nil {
		obj = C.envNewGlobalRef(self.env, obj)
		o = newObject(self, c, obj)
	} else {
		err = self.exceptionOccurred()
	}
	return
}


// returns a Class object;  the object will first be looked up in cache,
// and if not found there, resolved via Java and stored in the cache path.
// classes returned via /THIS/ channel, need not be unrefed, as they all
// hold a global ref.
//
// TODO: in truth, they should probably ALL be local-refs of the cached one...
func (self *Environment) GetClass(klass ClassName) (c *Class, err error) {
	c, err = self.findCachedClass(klass)
	if err == nil {
		return
	}
	s := C.CString(klass.AsPath())
	defer C.free(unsafe.Pointer(s))
	// print("envFindClass ", klass,"\n")
	kl := C.envFindClass(self.env, s)
	if kl == nil {
		err = self.exceptionOccurred()
	} else {
		err = nil // clear the cache error
		// print("found ", klass,"\n")
		kl = C.jclass(C.envNewGlobalRef(self.env, kl))
		c = newClass(self, klass, kl)
		self.classes[klass.AsPath()] = c
	}
	return
}
func (self *Environment) GetClassStr(klass string) (c *Class, err error) {
	class := NewClassName(klass)
	if err == nil {
		c, err = self.GetClass(class)
	}
	return
}

func (self *Environment) getObjectClass(o *Object) (c *Class, err error) {
	kl := C.envGetObjectClass(self.env, o.object)
	if kl == nil {
		err = self.exceptionOccurred()
	} else {
		c = self.LocalRefClass(newClass(self, nil, /*TODO: nil's probably not the best..*/ kl))
	}
	return
}

func (self *Environment) _objMethod(obj *Object, name string, jt JavaType, params ...interface{}) (meth *Method, err error) {
	class, err := self.getObjectClass(obj)
	defer self.LocalUnrefClass(class)
	if err != nil {
		return
	}
	form, err := formFor(self, jt, params...)
	if err != nil {
		return
	}

	cmethod := C.CString(name)
	defer C.free(unsafe.Pointer(cmethod))
	cform := C.CString(form)
	defer C.free(unsafe.Pointer(cform))

	m := C.envGetMethodID(self.env, class.class, cmethod, cform)
	if m == nil {
		err = self.exceptionOccurred()
	} else {
		meth = &Method{m}
	}
	return

}

func (self *Environment) _classMethod(class *Class, name string, jt JavaType, params ...interface{}) (meth *Method, err error) {
	form, err := formFor(self, jt, params...)
	if err != nil {
		return
	}
	cmethod := C.CString(name)
	defer C.free(unsafe.Pointer(cmethod))
	cform := C.CString(form)
	defer C.free(unsafe.Pointer(cform))
	//cname, err := class.Name()
	//if err != nil { return }
	//print("Looking for ", name, "\t", form, "\t in ", cname.AsPath(), "\n")
	m := C.envGetMethodID(self.env, class.class, cmethod, cform)
	if m == nil {
		err = self.exceptionOccurred()
	} else {
		meth = &Method{m}
	}
	return
}

func (self *Environment) _classStaticMethod(class *Class, name string, jt JavaType, params ...interface{}) (meth *Method, err error) {
	form, err := formFor(self, jt, params...)
	if err != nil {
		return
	}
	cmethod := C.CString(name)
	defer C.free(unsafe.Pointer(cmethod))
	cform := C.CString(form)
	defer C.free(unsafe.Pointer(cform))
	//cname, err := class.Name()
	//if err != nil { return }
	//print("Looking for (static)", name, "\t", form, "\t in ", cname.AsPath(), "\n")
	m := C.envGetStaticMethodID(self.env, class.class, cmethod, cform)
	if m == nil {
		err = self.exceptionOccurred()
	} else {
		meth = &Method{m}
	}
	return
}

type Exception struct {
	ex C.jthrowable
}

func (self *Exception) Error() string {
	return "{JavaException:<TODO>}"
}

func (self *Environment) exceptionOccurred() (ex *Exception) {
	throwable := C.envExceptionOccurred(self.env)
	if throwable != nil {
		ex = &Exception{throwable}
		if !self.quietExceptions {
			C.envExceptionDescribe(self.env)
		}
		C.envExceptionClear(self.env)
	}
	return
}

func (self *Environment) exceptionCheck() bool {
	return (C.envExceptionCheck(self.env) != C.JNI_FALSE)
}

// Syntactic sugar around &Class{C.jclass(LocalRef(&Object{C.jobject(class.class)}))}
func (self *Environment) LocalRefClass(c *Class) *Class {
	return newClass(self, c._klass, C.jclass(C.envNewLocalRef(self.env, c.class)))
}

// Syntactic sugar around LocalUnref(&Object{C.jobject(class.class)})
func (self *Environment) LocalUnrefClass(c *Class) {
	C.envDeleteLocalRef(self.env, c.class)
}

// Adds a 'local' ref to the JVM for Object, and returns an object that is contains reference
func (self *Environment) LocalRef(o *Object) *Object {
	return newObject(self, o._klass, C.envNewLocalRef(self.env, o.object))
}

// Release a local reference (returned from LocalRef) back to the JVM
func (self *Environment) LocalUnref(o *Object) {
	C.envDeleteLocalRef(self.env, o.object)
}

func (self *Environment) GlobalRef(o *Object) *Object {
	return newObject(self, o._klass, C.envNewGlobalRef(self.env, o.object))
}

/*
func (self *Environment)NewInstance(klass string, params ...interface{})(obj *Object, err os.Error){
	obj, err = env.NewString(klass, params...)
	return
}*/

/*func (self *Context)NewString(s string)(obj *Object, err os.Error){
	return self.env.NewStringObject(s)
}*/
