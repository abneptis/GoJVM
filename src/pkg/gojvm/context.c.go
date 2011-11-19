package gojvm

//#cgo LDFLAGS:-ljvm	-L/usr/lib/jvm/java-6-sun/jre/lib/amd64/server/
//#include</usr/lib/jvm/java-6-sun-1.6.0.26/include/jni.h>
//#include <stdlib.h>
//#include <libio.h>
//#include <unistd.h>
//#include "helpers.h"
import "C"
import
//	"log"

"strings"

type Context struct {
	jvm     *C.JavaVM
	Env     *Environment
	classes map[string]C.jclass
}

func newContext() (ctx *Context) {
	ctx = &Context{}
	ctx.Env = NewEnvironment()
	ctx.jvm = new(C.JavaVM)
	ctx.classes = map[string]C.jclass{}
	return
}


func initializeContext(args *C.JavaVMInitArgs) (ctx *Context, err error) {
	ctx = newContext()
	//eptr := unsafe.Pointer(&ctx.env)
	err = JVMError(C.newJVMContext(&ctx.jvm, &ctx.Env.env, args))
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
