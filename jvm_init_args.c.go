package gojvm

//#cgo LDFLAGS:-ljvm	-L/usr/lib/jvm/java-6-sun/jre/lib/amd64/server/
//#include</usr/lib/jvm/java-6-sun-1.6.0.26/include/jni.h>
//#include <stdlib.h>
//#include <unistd.h>
//#include "helpers.h"
import "C"
import (
	"os"
	"unsafe"
)

func DefaultJVMArgs(ver int)(args *C.JavaVMInitArgs, err os.Error){
	if ver == 0 { ver = DEFAULT_JVM_VERSION }
	args = new(C.JavaVMInitArgs)
	//print("Default args\t", ver,"\n")
	args.version = C.jint(ver)
	ok := C.JNI_GetDefaultJavaVMInitArgs(unsafe.Pointer(args))
	err = JVMError(ok)
	return
}

func AddStringArg(args *C.JavaVMInitArgs, s string)(err os.Error){
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	ok := C.addStringArgument(args, cstr)
	if ok != 0 {
		err = os.NewError("addStringArg failed")
	}
	return
}
