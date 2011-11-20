package gojvm

/* Tests various external classes pre-disposed to have certain.. 'issues'
 */
import (
	"testing"
	"gojvm/environment"
	"gojvm/types"
)

var TrivialClass = "org/golang/ext/gojvm/testing/Trivial"
/* Provides a variety of constructors, and one niladic getConstructorUsed;

Verifies:
	Method (parameter) reflection
	Ctx.NewClass,
	Marshalling (int,int64,string,{})
*/

var PathosClass = "org/golang/ext/gojvm/testing/Pathos"
/* Throws an exception on construction;
Verifies, exception on obj.New()
*/

/* Doesn't exist

Verifies:
	Exception for missing class	
*/
var MissingClass = "org/golang/ext/MissingClass"

type trivialClassTest struct {
	ConstArgs []interface{}
}

var trivialClassTests = []trivialClassTest{
	trivialClassTest{[]interface{}{}},
	trivialClassTest{[]interface{}{3}},
	trivialClassTest{[]interface{}{"aString"}},
	trivialClassTest{[]interface{}{int64(32)}},
}

func TestJVMTrivialClass(t *testing.T) {
	env := setupJVM(t)
	for i, test := range trivialClassTests {
		form, err := environment.FormFor(env, types.Basic(types.VoidKind), test.ConstArgs...)
		fatalIf(t, err != nil, "[%d] Error generating formFor: %v", i, err)
		fatalIf(t, form == "", "Got nil form")
		klass, err := env.NewInstanceStr(TrivialClass, test.ConstArgs...)
		fatalIf(t, err != nil, "[%d] Error generating formFor: %v", i, err)
		kused, _, err := klass.CallString(false, "getConstructorUsed")
		fatalIf(t, err != nil, "[%d] Error getting constructor used: %v", i, err)
		fatalIf(t, kused != form, "[%d] Constructor called was wrong (Exp: %s, got: %s)", form, kused)

		cn, err := klass.ClassName()
		fatalIf(t, err != nil, "Pathos name threw an error: %v", err)
		fatalIf(t, cn.AsPath() != TrivialClass, "Returned wrong name: (exp %q, got %q)", TrivialClass, cn)

	}
}

func TestJVMPathosClass(t *testing.T) {
	env := setupJVM(t)
	// We mute expected exceptions because otherwise the test looks sloppy (and FAILS are hard to see)
	defer defMute(env)()
	klass, err := env.NewInstanceStr(PathosClass)
	fatalIf(t, klass != nil, "Pathos should throw an exception (be nil), but got %v", klass)
	fatalIf(t, err == nil, "Pathos didn't throw an exception")
}

func TestJVMMissingClass(t *testing.T) {
	env := setupJVM(t)
	defer defMute(env)()
	// We mute expected exceptions because otherwise the test looks sloppy (and FAILS are hard to see)
	klass, err := env.NewInstanceStr(MissingClass)
	fatalIf(t, klass != nil, "Missing should throw an exception (be nil), but got %v", klass)
	fatalIf(t, err == nil, "Missing didn't throw an exception")
}
