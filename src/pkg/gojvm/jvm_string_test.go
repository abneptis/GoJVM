package gojvm

import (
	"testing"
	"math/rand"
	"strconv"
	"strings"	
)

var simpleStringTests = []string{
	"basic",
	"a modestly long string with no special special characters.",
	"", //empty
	"embedded\x00nulls",
	// this won't pass a simpleStringTest...
	// "a modestly long string with invalid \xc0\xaf special characters.",
	"\x00\x00\x00nulledstring",
	"κόσμε",
}

func TestJVMNewString(t *testing.T) {
	env := setupJVM(t)
	for i, str := range simpleStringTests {
		jstr, err := env.NewStringObject(str)
		fatalIf(t, err != nil, "[%d] Couldn't NewString '%q'", i, str)

		// length compare will fail on the UTF8 test, as length is 'runes' in java,
		// and bytes in Go.

		jstr_len, err := jstr.CallInt(env, false, "length")
		fatalIf(t, err != nil, "[%d] Couldn't call length on jstr '%v'", i, err)
		fatalIf(t, i != 5 && jstr_len != len(str), "[%d] Wrong length (Got %d, expected %d)", i, jstr_len, len(str))
		ostr, _, err := env.CallObjectString(jstr, false, "toString")
		fatalIf(t, err != nil, "[%d] Couldn't call toString on jstr '%v'", i, err)
		fatalIf(t, str != ostr, "[%d] Wrong inner string (Got %q, expected %q)", i, ostr, str)

	}
}


var someWords = strings.Split("Mary had a little lamb whose fleece was white as snow and every where that mary went the lamb was sure to go"," ")

// just a reference value for giving relative measures
// not a fair comparison to 'pure' JVM
func BenchmarkGoShortStringsReference(b *testing.B) {
	ll := int64(0)
	words := len(someWords)
	for i := 0; i < b.N; i++ {
		str := someWords[rand.Int() % words] + strconv.Itoa(i)
		ll += int64(len(str))
	}
	b.SetBytes(ll)
}

// Creates a 'random' short string, converts it into a Java String,
// and then immediately dereferences it.
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

