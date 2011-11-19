package gojvm

import (
	"os"
//	"log"
)


/* Returns the signature string for the configured environment, return type and parameters
	(exported for convenience) */
func FormFor(ctx *Environment, ret JavaType, params ...interface{})(s string, err os.Error){
	return formFor(ctx,ret,params...)
}	
func formFor(ctx *Environment, ret JavaType, params ...interface{})(s string, err os.Error){
	s, err = ParameterString(ctx, params...)
	if err == nil {
		s = s + ret.String()
	}
	return
}

func MethodSignature(in []JavaType, out JavaType)(string){
	return JavaMethodSignature{in,out}.String()
}

func ParameterTypes(ctx *Environment, params ...interface{})(jtypes []JavaType, err os.Error){
	for _, param := range(params){
		var jt JavaType
		jt, err = TypeOf(ctx, param)
		//log.Printf("param[%d] '%+v' %v", i, param, reflect.TypeOf(param).String())
		if err != nil { break }
		jtypes = append(jtypes, jt)
	}
	return
}

func ParameterString(ctx *Environment, params ...interface{})(s string, err os.Error){
	jms := JavaMethodSignature{}
	jms.Params, err = ParameterTypes(ctx, params...)
	if err == nil {
		s = jms.ParameterString()
	}
	return
}


type JavaMethodSignature struct {
	Params	[]JavaType
	Return	JavaType
}

func (self JavaMethodSignature)ParameterString()(string){
	plist := ""
	for _, parm := range(self.Params){
		plist += parm.String()
	}
	return "(" + plist + ")"
}

func (self JavaMethodSignature)String()(string){
	return self.ParameterString() + self.Return.String()
}

type JavaType interface {
	Kind()(ArgKind)
	String()(string)
}

type BasicType ArgKind

func (self BasicType)String()(string) { return string(self) }
func (self BasicType)Kind()(ArgKind) { return ArgKind(self) }

type ClassType struct {
	Klass	string
}

func (self ClassType)String()(string) { return "L" + self.Klass + ";" }
func (self ClassType)Kind()(ArgKind) { return JavaClassKind }

type ArrayType struct {
	Underlying	JavaType
}

func (self ArrayType)String()(string) { return "[" + self.Underlying.String() }
func (self ArrayType)Kind()(ArgKind) { return JavaArrayKind }



type ArgKind int
const (
	JavaBoolKind	ArgKind =	'Z'
	JavaByteKind	ArgKind =	'B'
	JavaCharKind	ArgKind =	'C'
	JavaShortKind	ArgKind =	'S'
	JavaIntKind		ArgKind =	'I'
	JavaLongKind	ArgKind =	'J'
	JavaFloatKind	ArgKind =	'F'
	JavaDoubleKind	ArgKind =	'D'
	JavaClassKind	ArgKind =	'L'
	JavaArrayKind	ArgKind =	'['
	JavaVoidKind	ArgKind =	'V'	// return only
)
