package gojvm

//#cgo LDFLAGS:-ljvm	-L/usr/lib/jvm/java-6-sun/jre/lib/amd64/server/
//#include "helpers.h"
import "C"
import (
	"gojvm/types"
//	"log"
)

/* represents a class (object) */
type Class struct {
	class  C.jclass
}

func newClass(class C.jclass) *Class {
	return &Class{class}
}

/*
	returns the (potentially cached) types.ClassName of the class.
*/
func (self *Class) GetName(env *Environment) (name types.ClassName, err error) {
		//log.Printf("ClassName(miss)")
	var cstr string
	cstr, _, err = self.CallString(env, false, "getName")
	if err == nil {
		name = types.NewClassName(cstr)
	}
	return
}

// Calls the named void-method on the class
func (self *Class) CallVoid(env *Environment, static bool, mname string, params ...interface{}) (err error) {
	return env.CallClassVoid(self, static, mname, params ...)
}

func (self *Class) CallInt(env *Environment, static bool, mname string, params ...interface{}) (i int, err error) {
	return env.CallClassInt(self, static, mname, params ...)
}

func (self *Class) CallDouble(env *Environment, static bool, mname string, params ...interface{}) (i float64, err error) {
	return env.CallClassDouble(self, static, mname, params ...)
}

func (self *Class) CallFloat(env *Environment, static bool, mname string, params ...interface{}) (i float32, err error) {
	return env.CallClassFloat(self, static, mname, params ...)
}

func (self *Class) CallObject(env *Environment, static bool, mname string, rval types.Typed, params ...interface{}) (o *Object, err error) {
	return env.CallClassObj(self, static, mname, rval, params ...)
}

func (self *Class) CallString(env *Environment, static bool, mname string, params ...interface{}) (str string, isnull bool, err error) {
	return env.CallClassString(self, static, mname, params ...)
}

