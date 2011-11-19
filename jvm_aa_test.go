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
	t.Logf("Testing -- using classpath [.,%s", DefaultJREPath)
	var err error
	_Ctx, err = InitializeJVM(0, []string{".", DefaultJREPath})
	fatalIf(t, err != nil, "Error initializing JVM: %v", err)
	fatalIf(t, _Ctx == nil, "Got a nil context!")
	// if you do(n't) like the java-exception noise in your tests, (un)comment: 
	//_Ctx.env.Mute(true)
	return _Ctx
}

// so the timing of other tests isn't thrown.
func TestJVMFirst(t *testing.T) { setupJVM(t) }
