package gojvm

//#cgo LDFLAGS:-ljvm	-L/usr/lib/jvm/java-6-sun/jre/lib/amd64/server/
//#include</usr/lib/jvm/java-6-sun-1.6.0.26/include/jni.h>
//#include <stdlib.h>
//#include <unistd.h>
import "C"
import (
	"os"
	"fmt"
)

type jvmError int
func (self jvmError)String()(string){
	return	fmt.Sprintf("<jvmError: %d>",int(self)) 
}
func JVMError(i C.jint)(err os.Error){
	if i != 0 {
		err = jvmError(int(i))
	}
	return
}

