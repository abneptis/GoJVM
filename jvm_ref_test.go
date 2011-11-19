package gojvm

import (
	"testing"
	"time"
)

var CleanerClass = "cc/qwe/gojvm/Cleaner"

var SystemClass = "java/lang/System"

func systemClass(ctx *Context, t *testing.T) (c *Class) {
	obj, err := ctx.Env.NewInstanceStr(SystemClass)
	fatalIf(t, err != nil, "Error loading system class: %v", err)
	fatalIf(t, obj._klass == nil, "Object has no class!: %+v", obj)
	return obj._klass
}

func TestJVMRefCounting(t *testing.T) {
	ctx := setupJVM(t)
	system := systemClass(ctx, t)
	cleaner, err := ctx.Env.NewInstanceStr(CleanerClass)
	fatalIf(t, err != nil, "Got an exception instantiating %s", CleanerClass)
	dead, err := cleaner.CallInt(false, "getDeadKids")
	fatalIf(t, err != nil, "Got an exception calling getDeadKids %v", err)
	fatalIf(t, dead != 0, "Wrong number of dead kids: (Got: %d, exp: %d)", dead, 0)
	for i := 0; i < 100; i++ {
		obj, err := cleaner.CallObj(false, "NewChild", ClassType{CleanerClass + "$Cleanable"})
		fatalIf(t, err != nil, "Got an exception calling NewChild %v", err)
		ctx.Env.LocalUnref(obj)
	}
	err = system.CallVoid(true, "gc")
	fatalIf(t, err != nil, "Got an exception calling gc() :%v", err)
	// gc is not a blocking call... nor is it required it would actually finalize all of our 
	// objects, but it seems to work on most JVMs...
	time.Sleep(500000000)
	dead, err = cleaner.CallInt(false, "getDeadKids")
	fatalIf(t, err != nil, "Got an exception calling getDeadKids %v", err)
	fatalIf(t, dead != 100, "Wrong number of dead kids: (Got: %d, exp: %d)", dead, 100)
}
