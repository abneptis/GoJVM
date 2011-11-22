package gojvm

//#cgo LDFLAGS:-ljvm	-L/usr/lib/jvm/java-6-sun/jre/lib/amd64/server/
//#include</usr/lib/jvm/java-6-sun-1.6.0.26/include/jni.h>
//#include <stdlib.h>
//#include <unistd.h>
//#include "helpers.h"
import "C"
import (
	"unsafe"
)



type callbacks struct {
	calls map[int]interface{}
	next int
}

var cCallbacks = callbacks{map[int]interface{}{}, 0}

func findCCallback(p int) (i interface{}) {
	if v, ok := cCallbacks.calls[p]; ok {
		i = v
	}
	return
}

func addGoCallback(fptr interface{})(id int){
	id = cCallbacks.next
	cCallbacks.next += 1
	cCallbacks.calls[id] = fptr
	return
}

//export callCCallback
func callCCallback(eptr unsafe.Pointer, ptr unsafe.Pointer, id int){
	cb := findCCallback(id)
	if cb != nil {
		print("Callback known")
	}
}
