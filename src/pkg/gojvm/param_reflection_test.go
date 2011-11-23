package gojvm

import (
	"gojvm/types"
	"testing"
)

type formForTest struct {
	rtype  types.Typed
	params []interface{}
	eform  string
	eerror error
}

var formOfTests = []formForTest{
	formForTest{types.Basic(types.VoidKind), []interface{}{}, "()V", nil},
	formForTest{types.Basic(types.VoidKind), []interface{}{int(5)}, "(I)V", nil},
	formForTest{types.Basic(types.VoidKind), []interface{}{int16(5)}, "(S)V", nil},
	formForTest{types.Basic(types.VoidKind), []interface{}{int64(5)}, "(J)V", nil},
	formForTest{types.Basic(types.VoidKind), []interface{}{"5"}, "(Ljava/lang/String;)V", nil},
	formForTest{types.Basic(types.VoidKind), []interface{}{"5", &Object{}}, "(Ljava/lang/String;Ljava/lang/Object;)V", nil},
	formForTest{types.Basic(types.ShortKind), []interface{}{int64(5)}, "(J)S", nil},
	formForTest{types.Basic(types.BoolKind), []interface{}{int(0)}, "(I)Z", nil},
	formForTest{types.Class{types.JavaLangObject}, []interface{}{int(0)}, "(I)Ljava/lang/Object;", nil},
	formForTest{types.ArrayType{types.Basic(types.ByteKind)}, []interface{}{[]byte{1, 2, 3}}, "([B)[B", nil},
}

func TestTrivialFormFor(t *testing.T) {
	env := setupJVM(t)
	for i, test := range formOfTests {
		form, err := FormFor(env, test.rtype, test.params...)
		fatalIf(t, err != test.eerror, "[%d] Unexpected error %v", i, err)
		fatalIf(t, form != test.eform, "[%d] Unexpected form  (got %s, wanted %s)", i, form, test.eform)
	}

}
