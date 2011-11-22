package gojvm

//#cgo LDFLAGS:-ljvm	-L/usr/lib/jvm/java-6-sun/jre/lib/amd64/server/
//#include</usr/lib/jvm/java-6-sun-1.6.0.26/include/jni.h>
//#include <stdlib.h>
//#include <unistd.h>
import "C"
import (
	"fmt"
	"unsafe"
)

type jvmError int

func (self jvmError) Error() string {
	return fmt.Sprintf("<jvmError: %d>", int(self))
}
func JVMError(i C.jint) (err error) {
	if i != 0 {
		err = jvmError(int(i))
	}
	return
}

//export TestingCallback
func TestingCallback(unsafe.Pointer, unsafe.Pointer){
	print("Testing callback...")
}

