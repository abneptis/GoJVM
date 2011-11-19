package types


// A 'Typed' contains information for method signatures
type Typed interface {
	Kind() Kind
	String() string
}


// returns the resultant String() of a JavaMethodSignature compiled from a set of
// JavaType's
func MethodSignatureString(in []Typed, out Typed) string {
	return MethodSignature{in, out}.String()
}


type MethodSignature struct {
	Params []Typed
	Return Typed
}

/* Implements the parameter-side only of a JavaMethodSignature.
the resultant string is () quoted.
*/
func (self MethodSignature) ParameterString() string {
	plist := ""
	for _, parm := range self.Params {
		plist += parm.String()
	}
	return "(" + plist + ")"
}

// Returns the Java method-signature as a string.
func (self MethodSignature) String() string {
	return self.ParameterString() + self.Return.String()
}


// a 'Basic' is a single 'primative' type
type Basic Kind

func (self Basic) String() string { return string(self) }
func (self Basic) Kind() Kind  { return Kind(self) }

// a 'Class' is a Java Class name type  a la "Lclass/path/name;"
type Class struct {
	Klass string
}

func (self Class) String() string { return "L" + self.Klass + ";" }
func (self Class) Kind() Kind  { return ClassKind }

// Java arrays consist of a single type (though the type itself could be
// the generic 'Object' class
type ArrayType struct {
	Underlying Typed
}

func (self ArrayType) String() string { return "[" + self.Underlying.String() }
func (self ArrayType) Kind() Kind  { return ArrayKind }

type Kind int

const (
	BoolKind   Kind = 'Z'
	ByteKind   Kind = 'B'
	CharKind   Kind = 'C'
	ShortKind  Kind = 'S'
	IntKind    Kind = 'I'
	LongKind   Kind = 'J'
	FloatKind  Kind = 'F'
	DoubleKind Kind = 'D'
	ClassKind  Kind = 'L'
	ArrayKind  Kind = '['
	VoidKind   Kind = 'V' // return only
)
