package gojvm

//#cgo LDFLAGS:-ljvm	-L/usr/lib/jvm/java-6-sun/jre/lib/amd64/server/
//#include</usr/lib/jvm/java-6-sun-1.6.0.26/include/jni.h>
//#include <stdlib.h>
//#include <unistd.h>
//#include "helpers.h"
import "C"

const (
	JNI_VERSION_1_2		= C.JNI_VERSION_1_2
	JNI_VERSION_1_4		= C.JNI_VERSION_1_4
	JNI_VERSION_1_6		= C.JNI_VERSION_1_6
)

const	DEFAULT_JVM_VERSION	= JNI_VERSION_1_6

const SystemDefaultJREPath = "/usr/lib/jvm/default-java/jre/lib"
const SunJREPath = "/usr/lib/jvm/java-6-sun/jre/lib"

var DefaultJREPath = SystemDefaultJREPath