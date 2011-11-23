package gojvm

//#include "helpers.h"
import "C"

import (
	"errors"
	"gojvm/types"
	"strings"
	"unsafe"
)

const (
	JNI_VERSION_1_2 = C.JNI_VERSION_1_2
	JNI_VERSION_1_4 = C.JNI_VERSION_1_4
	JNI_VERSION_1_6 = C.JNI_VERSION_1_6
)

const DEFAULT_JVM_VERSION = JNI_VERSION_1_6

const SystemDefaultJREPath = "/usr/lib/jvm/default-java/jre/lib"
const SunJREPath = "/usr/lib/jvm/java-6-sun/jre/lib"

var DefaultJREPath = SystemDefaultJREPath

type JVM struct {
	jvm        *C.JavaVM
	registered map[int]callbackDescriptor
	regId      int
}

func (self *JVM) addNative(env *Environment, f interface{}) (id int, csig types.MethodSignature, err error) {
	//func CallbackSignature(f interface{})(sig MethodSignature, err error){
	cbd, err := CallbackDescriptor(env, f)
	if err != nil {
		return
	}
	/// todo, this is not thread safe...
	id = self.regId
	self.registered[id] = cbd
	self.regId++
	csig = cbd.Signature
	//log.Printf("Calculated signature: %v", csig)
	/// TODO
	return
}

/* 
	Returns a new environment pointer that is appropriate for the /currently executing/
	thread to use.  It is safe to call multiply (idempotent), and can be returned 
	via DetachCurrentThread (not idempotent!)
*/
func (self *JVM) AttachCurrentThread() (env *Environment, err error) {
	env = NewEnvironment(self)
	//print ("Allocated environment for thread\t", env.Ptr(),"\n")
	if 0 != C.vmAttachCurrentThread(self.jvm, env.Ptr(), nil) {
		err = errors.New("Couldn't attach thread (and thus cannot gather exception)")
	} else {
		AllEnvs.Add(env)
	}
	return
}

// notifies the JVM of your threads done-ness w/ it, and deallocates the associated
// environment pointer.  Depending on the exact JDK version, there are differing semantics
// on whether the 'original' thread can call this (else JVM Shutdown), but most modern
// stacks (>=1.2) should allow this from the 'main' thread.
func (self *JVM) DetachCurrentThread() (err error) {
	if 0 != C.vmDetachCurrentThread(self.jvm) {
		err = errors.New("Couldn't attach thread (and thus cannot gather exception)")
	} else {
		print("TODO: DetachCurrentThread - don't know which thread I am :-( (cant be unmapped)\n")
	}
	return
}

type JvmConfig struct {
	ClassPath []string
}

func NewJVM(ver int, conf JvmConfig) (jvm *JVM, env *Environment, err error) {
	args, err := defaultJVMArgs(ver)
	if err != nil {
		return
	}
	pathStr := strings.Join(conf.ClassPath, ":")
	//print("Adding class path\n")
	err = addStringArg(args, "-Djava.class.path="+pathStr)
	if err == nil {
		//print("Initializing JVM Context\n")
		jvm = &JVM{
			registered: map[int]callbackDescriptor{},
		}
		env = NewEnvironment(jvm)
		if 0 != C.newJVMContext(&jvm.jvm, env.Ptr(), args) {
			err = errors.New("Couldn't instantiate JVM")
		} else {
			AllVMs.Add(jvm)
		}
	}
	return
}

func defaultJVMArgs(ver int) (args *C.JavaVMInitArgs, err error) {
	if ver == 0 {
		ver = DEFAULT_JVM_VERSION
	}
	args = new(C.JavaVMInitArgs)
	//print("Default args\t", ver,"\n")
	args.version = C.jint(ver)
	if 0 != C.JNI_GetDefaultJavaVMInitArgs(unsafe.Pointer(args)) {
		err = errors.New("Couldn't contruct default JVM args")
	}
	return
}

func addStringArg(args *C.JavaVMInitArgs, s string) (err error) {
	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))
	ok := C.addStringArgument(args, cstr)
	if ok != 0 {
		err = errors.New("addStringArg failed")
	}
	return
}
