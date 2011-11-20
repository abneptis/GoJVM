package gojvm

//#cgo LDFLAGS:-ljvm	-L/usr/lib/jvm/java-6-sun/jre/lib/amd64/server/
//#include "helpers.h"
import "C"
import (
	"strings"
	"gojvm/environment"
)
type Context struct {
	jvm     *C.JavaVM
	Env     *environment.Environment
}

func (self *Context)AttachCurrentThread()(env *environment.Environment, err error){
	env = environment.NewEnvironment()
	//print ("Allocated environment for thread\t", env.Ptr(),"\n")
	err = JVMError(C.vmAttachCurrentThread(self.jvm, env.Ptr(), nil))
	return
}

func newContext() (ctx *Context) {
	ctx = &Context{}
	ctx.Env = environment.NewEnvironment()
	ctx.jvm = new(C.JavaVM)
	return
}


func initializeContext(args *C.JavaVMInitArgs) (ctx *Context, err error) {
	ctx = newContext()
	envPtr := ctx.Env.Ptr()
	err = JVMError(C.newJVMContext(&ctx.jvm, envPtr, args))
	return
}

func InitializeJVM(ver int, cpath []string) (ctx *Context, err error) {
	args, err := DefaultJVMArgs(ver)
	if err != nil {
		return
	}
	pathStr := strings.Join(cpath, ":")
	//print("Adding class path\n")
	err = AddStringArg(args, "-Djava.class.path="+pathStr)
	if err == nil {
		//print("Initializing JVM Context\n")
		ctx, err = initializeContext(args)
	}
	return
}
