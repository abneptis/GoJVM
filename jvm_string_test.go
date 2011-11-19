package gojvm

import (
	"testing"
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
	ctx := setupJVM(t)
	for i, str := range simpleStringTests {
		jstr, err := ctx.Env.NewStringObject(str)
		fatalIf(t, err != nil, "[%d] Couldn't NewString '%q'", i, str)

		// length compare will fail on the UTF8 test, as length is 'runes' in java,
		// and bytes in Go.

		jstr_len, err := jstr.CallInt(false, "length")
		fatalIf(t, err != nil, "[%d] Couldn't call length on jstr '%v'", i, err)
		fatalIf(t, i != 5 && jstr_len != len(str), "[%d] Wrong length (Got %d, expected %d)", i, jstr_len, len(str))
		ostr, _, err := jstr.CallString(false, "toString")
		fatalIf(t, err != nil, "[%d] Couldn't call toString on jstr '%v'", i, err)
		fatalIf(t, str != ostr, "[%d] Wrong inner string (Got %q, expected %q)", i, ostr, str)

	}
}
