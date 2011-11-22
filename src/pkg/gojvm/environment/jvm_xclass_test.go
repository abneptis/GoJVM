package environment

/* Tests various external classes pre-disposed to have certain.. 'issues'
 */
import (
	"testing"
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

var NativeClass = "org/golang/ext/gojvm/testing/Native"
/* 
	Has attachable native methods
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
		form, err := FormFor(env, types.Basic(types.VoidKind), test.ConstArgs...)
		fatalIf(t, err != nil, "[%d] Error generating formFor: %v", i, err)
		fatalIf(t, form == "", "Got nil form")
		klass, err := env.NewInstanceStr(TrivialClass, test.ConstArgs...)
		fatalIf(t, err != nil, "[%d] Error generating formFor: %v", i, err)
		kused, _, err := klass.CallString(env, false, "getConstructorUsed")
		fatalIf(t, err != nil, "[%d] Error getting constructor used: %v", i, err)
		fatalIf(t, kused != form, "[%d] Constructor called was wrong (Exp: %s, got: %s)", form, kused)

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


func TestJVMNativeVoidClass(t *testing.T) {
	env := setupJVM(t)
	//defer defMute(env)()
	nativePings := 0
	klass, err := env.GetClassStr(NativeClass)
	fatalIf(t, err != nil, "Native threw an exception", err)
	fatalIf(t, klass == nil, "Native klass is nil!")
	err = env.RegisterNative(klass, "NativePing", func(E *Environment, O *Object){
		nativePings += 1
	})
	fatalIf(t, err != nil, "RegisterNative threw an exception", err)
	obj, err := env.NewInstanceStr(NativeClass)
	fatalIf(t, err != nil, "Couldn't instantiate NativeClass: %v", err)
	fatalIf(t, obj == nil, "Instantiated NativeClass is nil")
	err = obj.CallVoid(env, false, "NativePing")
	fatalIf(t, err != nil, "Couldn't call NativeClass.NativePing(): %v", err)
	fatalIf(t, nativePings != 1, "Wrong native ping count: %d\n", nativePings)
}

func TestJVMNativeIntClass(t *testing.T) {
	env := setupJVM(t)
	//defer defMute(env)()
	klass, err := env.GetClassStr(NativeClass)
	fatalIf(t, err != nil, "Native threw an exception", err)
	fatalIf(t, klass == nil, "Native klass is nil!")
	err = env.RegisterNative(klass, "NativeInt", func(E *Environment, O *Object)(i int){
		return 15
	})
	fatalIf(t, err != nil, "RegisterNative threw an exception", err)
	obj, err := env.NewInstanceStr(NativeClass)
	fatalIf(t, err != nil, "Couldn't instantiate NativeClass: %v", err)
	fatalIf(t, obj == nil, "Instantiated NativeClass is nil")
	ival, err := obj.CallInt(env, false, "NativeInt")
	fatalIf(t, err != nil, "Couldn't call NativeClass.NativeInt(): %v", err)
	fatalIf(t, ival != 15, "wrong returned value from native: %d", ival)
}

func TestJVMNativeComplexClass(t *testing.T) {
	env := setupJVM(t)
	//defer defMute(env)()
	klass, err := env.GetClassStr(NativeClass)
	fatalIf(t, err != nil, "Native threw an exception", err)
	fatalIf(t, klass == nil, "Native klass is nil!")
	obj1, err := env.NewInstanceStr("java/lang/Object")
	fatalIf(t, err != nil, "new(Object) threw an exception", err)
	obj2, err := env.NewInstanceStr("java/lang/Object")
	fatalIf(t, err != nil, "new(Object2) threw an exception", err)
	hit := false
	err = env.RegisterNative(klass, "NativeComplex", func(E *Environment, O *Object, o1 *Object, o2 *Object, i1 int)(){
		hit = true
	})
	fatalIf(t, err != nil, "RegisterNative threw an exception", err)
	obj, err := env.NewInstanceStr(NativeClass)
	fatalIf(t, err != nil, "Couldn't instantiate NativeClass: %v", err)
	fatalIf(t, obj == nil, "Instantiated NativeClass is nil")
	err = env.CallObjectVoid(obj, false, "NativeComplex", obj1, obj2, 13)
	fatalIf(t, err != nil, "Couldn't call NativeClass.NativeComplex(): %v", err)
	fatalIf(t, !hit, "Native complex never got called", err)
}


func TestJVMNativeBoolClass(t *testing.T) {
	env := setupJVM(t)
	//defer defMute(env)()
	klass, err := env.GetClassStr(NativeClass)
	fatalIf(t, err != nil, "Native threw an exception", err)
	fatalIf(t, klass == nil, "Native klass is nil!")
	hit := false
	err = env.RegisterNative(klass, "NativeBool", func(E *Environment, O *Object)(bool){
		return !hit
	})
	fatalIf(t, err != nil, "RegisterNative threw an exception", err)
	obj, err := env.NewInstanceStr(NativeClass)
	fatalIf(t, err != nil, "Couldn't instantiate NativeClass: %v", err)
	fatalIf(t, obj == nil, "Instantiated NativeClass is nil")
	ok, err := obj.CallBool(env, false, "NativeBool")
	fatalIf(t, err != nil, "Couldn't call NativeClass.NativeBool(): %v", err)
	fatalIf(t, ok == hit, "Native complex never got called", err)
	hit = !hit
	ok, err = obj.CallBool(env, false, "NativeBool")
	fatalIf(t, err != nil, "Couldn't call NativeClass.NativeBool(): %v", err)
	fatalIf(t, ok == hit, "Native complex never got called", err)
}

func TestJVMNativeLongClass(t *testing.T) {
	env := setupJVM(t)
	//defer defMute(env)()
	klass, err := env.GetClassStr(NativeClass)
	fatalIf(t, err != nil, "Native threw an exception", err)
	fatalIf(t, klass == nil, "Native klass is nil!")
	hit := int64(0)
	err = env.RegisterNative(klass, "NativeLong", func(E *Environment, O *Object)(int64){
		return hit
	})
	fatalIf(t, err != nil, "RegisterNative threw an exception", err)
	obj, err := env.NewInstanceStr(NativeClass)
	fatalIf(t, err != nil, "Couldn't instantiate NativeClass: %v", err)
	fatalIf(t, obj == nil, "Instantiated NativeClass is nil")
	ok, err := obj.CallLong(env, false, "NativeLong")
	fatalIf(t, err != nil, "Couldn't call NativeClass.NativeBool(): %v", err)
	fatalIf(t, ok != hit, "NativeLong got wrong value: %d", ok)
	hit = -5128
	ok, err = obj.CallLong(env, false, "NativeLong")
	fatalIf(t, err != nil, "Couldn't call NativeClass.NativeBool(): %v", err)
	fatalIf(t, ok != hit, "NativeLong got wrong value: %d", ok)
}

func TestJVMNativeFloatClass(t *testing.T) {
	env := setupJVM(t)
	//defer defMute(env)()
	klass, err := env.GetClassStr(NativeClass)
	fatalIf(t, err != nil, "Native threw an exception", err)
	fatalIf(t, klass == nil, "Native klass is nil!")
	hit := float32(.1234)
	err = env.RegisterNative(klass, "NativeFloat", func(E *Environment, O *Object)(float32){
		return hit
	})
	fatalIf(t, err != nil, "RegisterNative threw an exception", err)
	obj, err := env.NewInstanceStr(NativeClass)
	fatalIf(t, err != nil, "Couldn't instantiate NativeClass: %v", err)
	fatalIf(t, obj == nil, "Instantiated NativeClass is nil")
	ok, err := obj.CallFloat(env, false, "NativeFloat")
	fatalIf(t, err != nil, "Couldn't call NativeClass.NativeFloat(): %v", err)
	fatalIf(t, ok != hit, "NativeLong got wrong value: %d", ok)
	hit = 1/20
	ok, err = obj.CallFloat(env, false, "NativeFloat")
	fatalIf(t, err != nil, "Couldn't call NativeClass.NativeFloat(): %v", err)
	fatalIf(t, ok != hit, "NativeFloat got wrong value: %d", ok)
}

func TestJVMNativeShortClass(t *testing.T) {
	env := setupJVM(t)
	//defer defMute(env)()
	klass, err := env.GetClassStr(NativeClass)
	fatalIf(t, err != nil, "Native threw an exception", err)
	fatalIf(t, klass == nil, "Native klass is nil!")
	hit := int16(0)
	err = env.RegisterNative(klass, "NativeShort", func(E *Environment, O *Object)(int16){
		return hit
	})
	fatalIf(t, err != nil, "RegisterNative threw an exception", err)
	obj, err := env.NewInstanceStr(NativeClass)
	fatalIf(t, err != nil, "Couldn't instantiate NativeClass: %v", err)
	fatalIf(t, obj == nil, "Instantiated NativeClass is nil")
	ok, err := obj.CallShort(env, false, "NativeShort")
	fatalIf(t, err != nil, "Couldn't call NativeClass.NativeBool(): %v", err)
	fatalIf(t, ok != hit, "NativeShort got wrong value: %d", ok)
	hit = -5128
	ok, err = obj.CallShort(env, false, "NativeShort")
	fatalIf(t, err != nil, "Couldn't call NativeClass.NativeBool(): %v", err)
	fatalIf(t, ok != hit, "NativeShort got wrong value: %d", ok)
}

func TestJVMNativeDoubleClass(t *testing.T) {
	env := setupJVM(t)
	//defer defMute(env)()
	klass, err := env.GetClassStr(NativeClass)
	fatalIf(t, err != nil, "Native threw an exception", err)
	fatalIf(t, klass == nil, "Native klass is nil!")
	hit := float64(1234)
	err = env.RegisterNative(klass, "NativeDouble", func(E *Environment, O *Object)(float64){
		return hit
	})
	fatalIf(t, err != nil, "RegisterNative threw an exception", err)
	obj, err := env.NewInstanceStr(NativeClass)
	fatalIf(t, err != nil, "Couldn't instantiate NativeClass: %v", err)
	defer env.DeleteLocalRef(obj)
	fatalIf(t, obj == nil, "Instantiated NativeClass is nil")
	ok, err := obj.CallDouble(env, false, "NativeDouble")
	fatalIf(t, err != nil, "Couldn't call NativeClass.NativeDouble(): %v", err)
	fatalIf(t, ok != hit, "NativeDouble got wrong value: %d", ok)
	hit = float64(-125/7)
	ok, err = obj.CallDouble(env, false, "NativeDouble")
	fatalIf(t, err != nil, "Couldn't call NativeClass.NativeDouble(): %v", err)
	fatalIf(t, ok != hit, "NativeDouble got wrong value: %d", ok)
}

func TestJVMNativeStringClass(t *testing.T) {
	env := setupJVM(t)
	//defer defMute(env)()
	klass, err := env.GetClassStr(NativeClass)
	fatalIf(t, err != nil, "Native threw an exception", err)
	fatalIf(t, klass == nil, "Native klass is nil!")
	s := "test-string"
	err = env.RegisterNative(klass, "NativeString", func(E *Environment, O *Object)(string){
		return s
	})
	fatalIf(t, err != nil, "RegisterNative threw an exception", err)
	obj, err := env.NewInstanceStr(NativeClass)
	fatalIf(t, err != nil, "Couldn't instantiate NativeClass: %v", err)
	fatalIf(t, obj == nil, "Instantiated NativeClass is nil")
	ok,_, err := obj.CallString(env, false, "NativeString")
	fatalIf(t, err != nil, "Couldn't call NativeClass.NativeString(): %v", err)
	fatalIf(t, ok != s, "NativeString got wrong value: %d", ok)
	s = "testStr2"
	ok,_, err = obj.CallString(env, false, "NativeString")
	fatalIf(t, err != nil, "Couldn't call NativeClass.NativeString(): %v", err)
	fatalIf(t, ok != s, "NativeString got wrong value: %d", ok)
}

func BenchmarkJVMNativePing(b *testing.B) {
	env := setupJVM(nil)
	//defer defMute(env)()
	// We mute expected exceptions because otherwise the test looks sloppy (and FAILS are hard to see)
	nativePings := 0
	klass, err := env.GetClassStr(NativeClass)
	if err != nil {
		print("benchmark failed: ", err.Error(), "\n")
		return
	}
	err = env.UnregisterNatives(klass)
	if err != nil {
		print("benchmark failed: ", err.Error(), "\n")
		return
	}
	err = env.RegisterNative(klass, "NativePing", func(E *Environment, O *Object){
		nativePings += 1
	})
	if err != nil {
		print("benchmark failed: ", err.Error(), "\n")
		return
	}
	obj, err := env.NewInstanceStr(NativeClass)
	if err != nil {
		print("benchmark failed: ", err.Error(), "\n")
		return
	}
	ll := int64(0)
	for i := 0; i < b.N; i++{
		obj.CallVoid(env, false, "NativePing")
		ll += 1
	}
	env.DeleteLocalRef(obj)
	
	b.SetBytes(ll)
}

