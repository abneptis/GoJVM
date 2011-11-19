package gojvm

//	"log"

/* 
	Returns the signature string for the configured environment, return type and parameters
	(exported for convenience) ;
	
	Reflection is handled by ParameterString
*/
func FormFor(ctx *Environment, ret JavaType, params ...interface{}) (s string, err error) {
	return formFor(ctx, ret, params...)
}

// documented as FormFor
func formFor(ctx *Environment, ret JavaType, params ...interface{}) (s string, err error) {
	s, err = ParameterString(ctx, params...)
	if err == nil {
		s = s + ret.String()
	}
	return
}


// returns the resultant String() of a JavaMethodSignature compiled from a set of
// JavaType's
func MethodSignature(in []JavaType, out JavaType) string {
	return JavaMethodSignature{in, out}.String()
}

/*
	Returns an array of JavaTypes matching the params list (which may be empty, resulting
	in a valid empty (nil) list of JavaTypes).
	
	Reflection is done by TypeOf.

	* Environment is required in order to interpolate abstract objects into
	  class names.
*/
func ParameterTypes(ctx *Environment, params ...interface{}) (jtypes []JavaType, err error) {
	for _, param := range params {
		var jt JavaType
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
	jms := JavaMethodSignature{}
	jms.Params, err = ParameterTypes(ctx, params...)
	if err == nil {
		s = jms.ParameterString()
	}
	return
}

type JavaMethodSignature struct {
	Params []JavaType
	Return JavaType
}

/* Implements the parameter-side only of a JavaMethodSignature.
the resultant string is () quoted.
*/
func (self JavaMethodSignature) ParameterString() string {
	plist := ""
	for _, parm := range self.Params {
		plist += parm.String()
	}
	return "(" + plist + ")"
}

// Returns the Java method-signature as a string.
func (self JavaMethodSignature) String() string {
	return self.ParameterString() + self.Return.String()
}

type JavaType interface {
	Kind() ArgKind
	String() string
}

type BasicType ArgKind

func (self BasicType) String() string { return string(self) }
func (self BasicType) Kind() ArgKind  { return ArgKind(self) }

type ClassType struct {
	Klass string
}

func (self ClassType) String() string { return "L" + self.Klass + ";" }
func (self ClassType) Kind() ArgKind  { return JavaClassKind }

type ArrayType struct {
	Underlying JavaType
}

func (self ArrayType) String() string { return "[" + self.Underlying.String() }
func (self ArrayType) Kind() ArgKind  { return JavaArrayKind }

type ArgKind int

const (
	JavaBoolKind   ArgKind = 'Z'
	JavaByteKind   ArgKind = 'B'
	JavaCharKind   ArgKind = 'C'
	JavaShortKind  ArgKind = 'S'
	JavaIntKind    ArgKind = 'I'
	JavaLongKind   ArgKind = 'J'
	JavaFloatKind  ArgKind = 'F'
	JavaDoubleKind ArgKind = 'D'
	JavaClassKind  ArgKind = 'L'
	JavaArrayKind  ArgKind = '['
	JavaVoidKind   ArgKind = 'V' // return only
)
