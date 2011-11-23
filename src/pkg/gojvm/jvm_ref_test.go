package gojvm

import (
	"testing"
	"time"
	"gojvm/types"
	"math/rand"
)

var CleanerClass = types.ClassName{"org","golang","ext","gojvm","testing","Cleaner"}
var CleanableClass = types.ClassName{"org","golang","ext","gojvm","testing","Cleaner$Cleanable"}


// tests basic 'created & destroyed' semantics of our interaction w. the JVM.
// a better test would put some objects through their paces, release, and 'somehow'
// get the number of refs unreaped...
func TestJVMBasicRefCounting(t *testing.T) {
	env := setupJVM(t)
	system := systemClass(env, t)
	cleaner, err := env.NewInstanceStr(CleanerClass.AsPath())
	fatalIf(t, err != nil, "Got an exception instantiating %s", CleanerClass.String())
	dead, err := cleaner.CallInt(env, false, "getDeadKids")
	fatalIf(t, err != nil, "Got an exception calling getDeadKids %v", err)
	fatalIf(t, dead != 0, "Wrong number of dead kids: (Got: %d, exp: %d)", dead, 0)
	for i := 0; i < 100; i++ {
		obj, err := cleaner.CallObj(env, false, "NewChild", types.Class{CleanableClass})
		fatalIf(t, err != nil, "Got an exception calling NewChild %v", err)
		env.DeleteLocalRef(obj)
	}
	err = system.CallVoid(env, true, "gc")
	fatalIf(t, err != nil, "Got an exception calling gc() :%v", err)
	// gc is not a blocking call... nor is it required it would actually finalize all of our 
	// objects, but it seems to work on most JVMs...
	time.Sleep(500000000)
	dead, err = cleaner.CallInt(env, false, "getDeadKids")
	fatalIf(t, err != nil, "Got an exception calling getDeadKids %v", err)
	fatalIf(t, dead != 100, "Wrong number of dead kids: (Got: %d, exp: %d)", dead, 100)
}



// simply for comparison, a totally unfair test
// since java has to break into 'native' for us.
func BenchmarkGoRandInt(b *testing.B) {
	ll := int64(0)
	for i := 0; i < b.N; i++ {
		_ = rand.Int()
		ll += 4	// assuming 32 byte ints...
	}
	b.SetBytes(ll)
}


func BenchmarkJavaRandInt(b *testing.B) {
	env := setupJVM(nil)
	ll := int64(0)
	robj, err := env.NewInstanceStr("java/util/Random")
	if err != nil {
		panic("BenchmarkJavaRand failed to setup class")
	}
	for i := 0; i < b.N; i++ {
		_, err = robj.CallInt(env, false, "nextInt")
		if err == nil {
			ll += 4	// assuming 32 bit ints...
		}
	}
	env.DeleteLocalRef(robj)
	b.SetBytes(ll)
}


func BenchmarkJavaRandLong(b *testing.B) {
	env := setupJVM(nil)
	ll := int64(0)
	robj, err := env.NewInstanceStr("java/util/Random")
	if err != nil {
		panic("BenchmarkJavaRand failed to setup class")
	}
	for i := 0; i < b.N; i++ {
		_, err = robj.CallLong(env, false, "nextLong")
		if err == nil {
			ll += 8	// assuming 64 bit longs...
		}
	}
	env.DeleteLocalRef(robj)
	b.SetBytes(ll)
}
