package gojvm
/*
	This file exists (and is named) primarily for common test functionality;
	We also include the first SetupJVM call in here, so that it will not skew the timing
	of other tests regardless of other jvm_XX_test files
*/
import "testing"

func fatalIf(t *testing.T, tv bool, msg string, args ...interface{}) {
	if tv {
		t.Fatalf(msg, args...)
	}
}

var _Ctx *Context

func setupJVM(t *testing.T) *Context {
	if _Ctx != nil {
		return _Ctx
	}
	t.Logf("Testing -- using classpath [../../../java/,%s", DefaultJREPath)
	var err error
	_Ctx, err = InitializeJVM(0, []string{"../../../java/", DefaultJREPath})
	fatalIf(t, err != nil, "Error initializing JVM: %v", err)
	fatalIf(t, _Ctx == nil, "Got a nil context!")
	// expected exceptions are pre-muted/unmuted, but if you're testing something
	// that causes them to throw, and want readable tests, this is the line
	// to uncomment.
	//_Ctx.env.Mute(true)
	return _Ctx
}

// so the timing of other tests/bench's isn't thrown.
func TestJVMFirst(t *testing.T) { setupJVM(t) }
