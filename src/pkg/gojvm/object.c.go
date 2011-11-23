package gojvm

//#cgo CFLAGS:-I../include/
//#cgo LDFLAGS:-ljvm	-L/usr/lib/jvm/java-6-sun/jre/lib/amd64/server/
//#include "helpers.h"
import "C"
import (
	"gojvm/types"
	"log"
)

type Object struct {
	object C.jobject
}

// returns a new object value with specified parameters
// NB: refs are NOT adjusted directly by this call! Use it as a casting/construction-helper,
// not a Clone()
func newObject(obj C.jobject) *Object {
	return &Object{obj}
}

/* 
	Returns the Class() associated with the object
*/
func (self *Object) ObjectClass(env *Environment) (c *Class, err error) {
	return env.GetObjectClass(self)
}

/*
	Returns the (potentially cached) name of the ObjectClass of the
	named object.
*/
func (self *Object) Name(env *Environment) (name types.Name, err error) {
	var c *Class
	c, err = self.ObjectClass(env)
	if err == nil {
		defer env.DeleteLocalClassRef(c)
		name, err = c.GetName(env)
	} else {
		log.Printf("Couldn't get object class!")
	}
	return
}

func (self *Object) CallVoid(env *Environment, static bool, mname string, params ...interface{}) (err error) {
	return env.CallObjectVoid(self, static, mname, params...)
}

func (self *Object) CallInt(env *Environment, static bool, mname string, params ...interface{}) (i int, err error) {
	return env.CallObjectInt(self, static, mname, params...)
}

func (self *Object) CallLong(env *Environment, static bool, mname string, params ...interface{}) (i int64, err error) {
	return env.CallObjectLong(self, static, mname, params...)
}

func (self *Object) CallBool(env *Environment, static bool, mname string, params ...interface{}) (i bool, err error) {
	return env.CallObjectBool(self, static, mname, params...)
}

func (self *Object) CallFloat(env *Environment, static bool, mname string, params ...interface{}) (i float32, err error) {
	return env.CallObjectFloat(self, static, mname, params...)
}

func (self *Object) CallShort(env *Environment, static bool, mname string, params ...interface{}) (i int16, err error) {
	return env.CallObjectShort(self, static, mname, params...)
}

func (self *Object) CallDouble(env *Environment, static bool, mname string, params ...interface{}) (i float64, err error) {
	return env.CallObjectDouble(self, static, mname, params...)
}

// Calls the named Object-method on the object instance
func (self *Object) CallObj(env *Environment, static bool, mname string, rval types.Typed, params ...interface{}) (vObj *Object, err error) {
	return env.CallObjectObj(self, static, mname, rval, params...)
}

/* 
	A wrapper around ObjCallObj specific to java/lang/String, that will return the result as a GoString 

	A null string returned with no Exception can be differentiated via the wasNull return value.
*/
func (self *Object) CallString(env *Environment, static bool, mname string, params ...interface{}) (str string, wasNull bool, err error) {
	return env.CallObjectString(self, static, mname, params...)
}
