package environment

import (
	"gojvm/types"
	"reflect"
)

type callbackDescriptor struct {
	Signature	types.MethodSignature
	PTypes		[]reflect.Type
	// name?
	F	interface{}	// func
}

// Reflects a given function (w/ the specifid environment), and returns a callbackDescriptor
// including the reflected signature, reflected parameter types and the function pointer
// itself.  This descriptor is used in the global tables when using 'registerNatives'
func CallbackDescriptor(env *Environment, f interface{})(cd callbackDescriptor, err error){
	csig, err := ReflectedSignature(env, f)
	if err !=  nil { return }
	cd.Signature = csig
	cd.F = f
	// reflected sig has already done basic verification
	// of params & return
	//print("Reflected signature is ", cd.Signature.String(), "\n")
	rfv := reflect.TypeOf(f)
	for i := 2; i < rfv.NumIn(); i++ {
		cd.PTypes = append(cd.PTypes, rfv.In(i))
	}
	return
}