package gojvm

import (
	"testing"
	"time"
	"gojvm/types"
	"strconv"
	"strings"
	"math/rand"
)

var CleanerClass = "org/golang/ext/gojvm/testing/Cleaner"


// tests basic 'created & destroyed' semantics of our interaction w. the JVM.
// a better test would put some objects through their paces, release, and 'somehow'
// get the number of refs unreaped...
func TestJVMBasicRefCounting(t *testing.T) {
	env := setupJVM(t)
	system := systemClass(env, t)
	cleaner, err := env.NewInstanceStr(CleanerClass)
	fatalIf(t, err != nil, "Got an exception instantiating %s", CleanerClass)
	dead, err := cleaner.CallInt(false, "getDeadKids")
	fatalIf(t, err != nil, "Got an exception calling getDeadKids %v", err)
	fatalIf(t, dead != 0, "Wrong number of dead kids: (Got: %d, exp: %d)", dead, 0)
	for i := 0; i < 100; i++ {
		obj, err := cleaner.CallObj(false, "NewChild", types.Class{CleanerClass + "$Cleanable"})
		fatalIf(t, err != nil, "Got an exception calling NewChild %v", err)
		env.DeleteLocalRef(obj)
	}
	err = system.CallVoid(true, "gc")
	fatalIf(t, err != nil, "Got an exception calling gc() :%v", err)
	// gc is not a blocking call... nor is it required it would actually finalize all of our 
	// objects, but it seems to work on most JVMs...
	time.Sleep(500000000)
	dead, err = cleaner.CallInt(false, "getDeadKids")
	fatalIf(t, err != nil, "Got an exception calling getDeadKids %v", err)
	fatalIf(t, dead != 100, "Wrong number of dead kids: (Got: %d, exp: %d)", dead, 100)
}


var someWords = strings.Split("Mary had a little lamb whose fleece was white as snow and every where that mary went the lamb was sure to go"," ")

func BenchmarkGoShortStringsReference(b *testing.B) {
	ll := int64(0)
	words := len(someWords)
	for i := 0; i < b.N; i++ {
		str := someWords[rand.Int() % words] + strconv.Itoa(i)
		ll += int64(len(str))
	}
	b.SetBytes(ll)
}

func BenchmarkShortStrings(b *testing.B) {
	env := setupJVM(nil)
	ll := int64(0)
	words := len(someWords)
	for i := 0; i < b.N; i++ {
		str := someWords[rand.Int() % words] + strconv.Itoa(i)
		obj, err := env.NewStringObject(str)
		if err == nil {
			env.DeleteLocalRef(obj)
			ll += int64(len(str))
		} // else why didn't the tests fail...
	}
	
	b.SetBytes(ll)
}


// Benchmarks a long string conversion (but under a page)
func BenchmarkLongStrings(b *testing.B) {
	// ~2048 bytes (1/2 page) @ ~4ch/word => 512
	wordsPer := 512
	env := setupJVM(nil)
	ll := int64(0)
	words := len(someWords)
	for i := 0; i < b.N; i++ {
		str := ""
		for j := 0; j < wordsPer; j++ {
			str += someWords[rand.Int() % words]
		}
		obj, err := env.NewStringObject(str)
		if err == nil {
			env.DeleteLocalRef(obj)
			ll += int64(len(str))
		} // else why didn't the tests fail...
	}
	
	b.SetBytes(ll)
}


// Benchmarks a very long string conversion (~2 pages, possibly 3)
func BenchmarkVeryLongStrings(b *testing.B) {
	// ~8192 bytes (2 pages) @ ~4ch/word => 2048
	wordsPer := 2048
	env := setupJVM(nil)
	ll := int64(0)
	words := len(someWords)
	for i := 0; i < b.N; i++ {
		str := ""
		for j := 0; j < wordsPer; j++ {
			str += someWords[rand.Int() % words]
		}
		obj, err := env.NewStringObject(str)
		if err == nil {
			env.DeleteLocalRef(obj)
			ll += int64(len(str))
		} // else why didn't the tests fail...
	}
	
	b.SetBytes(ll)
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
		_, err = robj.CallInt(false, "nextInt")
		if err == nil {
			ll += 4	// assuming 32 byte ints...
		}
	}
	env.DeleteLocalRef(robj)
	b.SetBytes(ll)
}
