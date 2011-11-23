package gojvm

import (
	"gojvm/types"
//	"reflect"
//	"log"
)

/* 
	Returns the signature string for the configured environment, return type and parameters
	(exported for convenience) ;
	
	Reflection is handled by ParameterString
*/
func FormFor(ctx *Environment, ret types.Typed, params ...interface{}) (s string, err error) {
	return formFor(ctx, ret, params...)
}

// documented as FormFor
func formFor(ctx *Environment, ret types.Typed, params ...interface{}) (s string, err error) {
	s, err = ParameterString(ctx, params...)
	if err == nil {
		s = s + ret.TypeString()
	}
	return
}

/*
	Returns an array of types.Typeds matching the params list (which may be empty, resulting
	in a valid empty (nil) list of types.Typeds).
	
	Reflection is done by TypeOf.

	* Environment is required in order to interpolate abstract objects into
	  class names.
*/
func ParameterTypes(ctx *Environment, params ...interface{}) (jtypes []types.Typed, err error) {
	for _, param := range params {
		var jt types.Typed
		jt, err = TypeOf(ctx, param)
		//log.Printf("param[%d] '%+v' %v", i, param, reflect.TypeOf(param).String())
		if err != nil {
			break
		}
		jtypes = append(jtypes, jt)
	}
	return
}

/* Helper function for the parameter-side only of a JavaMethodSignature.
the resultant string is () quoted by jms.ParameterString()
*/
func ParameterString(ctx *Environment, params ...interface{}) (s string, err error) {
	jms := types.MethodSignature{}
	jms.Params, err = ParameterTypes(ctx, params...)
	if err == nil {
		s = jms.ParameterString()
	}
	return
}

