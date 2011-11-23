package gojvm
/*
	This file exists (and is named) primarily for common test functionality;
	We also include the first SetupJVM call in here, so that it will not skew the timing
	of other tests regardless of other jvm_XX_test files
*/
import (
	"testing"
	"flag"
	"sync"
)

func fatalIf(t *testing.T, tv bool, msg string, args ...interface{}) {
	if tv {
		t.Fatalf(msg, args...)
	}
}

func fatalEquals(t *testing.T, val interface{}, val2 interface{}, msg string, args ...interface{}) {
	args = append(args, []interface{}{val, val2}...)
	fatalIf(t, val == val2, msg + " (%v == %v)", args...)
}

func fatalInEq(t *testing.T, val interface{}, val2 interface{}, msg string, args ...interface{}) {
	args = append(args, []interface{}{val, val2}...)
	fatalIf(t, val != val2, msg + " (Expected: %v;\t Got: %v)", args...)
}

var SystemClass = "java/lang/System"

func systemClass(env *Environment, t *testing.T) (c *Class) {
	c, err := env.GetClassStr(SystemClass)
	fatalIf(t, err != nil, "Error loading system class: %v", err)
	return
}


var _jvm *JVM
var squelchExceptions bool /* = false */

// used in testing;  a 'squelch' helper
// such that:
//  func X(){
//    defer env.defMute()() /*note the double parens!!!*/
//    doSomeJavaCall
//  }
//
// would not output an exception to the console during processing
// regardless othe explicit 'mutedness'.
// there is a race condition here, but you're not supposed
// to be using *Environment in multiple threads anyhow :P
func defMute(env *Environment)(func()){
  muted := env.Muted()
  env.Mute(true)
  return func(){
    env.Mute(muted)
  }
}




var startLock = &sync.Mutex{}
func setupJVM(t *testing.T) (env *Environment){
	startLock.Lock()
	defer startLock.Unlock()
	if _jvm != nil {
		var err error
		env, err = _jvm.AttachCurrentThread()
		if err != nil {
			t.Fatalf("Couldn't attach thread: %v", err)
		}
		return
	}
	t.Logf("Testing -- using classpath [../../../java/,%s", DefaultJREPath)
	var err error
	_jvm, env, err = NewJVM(0, JvmConfig{[]string{"../../../java/", DefaultJREPath}})
	fatalIf(t, err != nil, "Error initializing JVM: %v", err)
	fatalIf(t, _jvm == nil, "Got a nil context!")
	// expected exceptions are pre-muted/unmuted, but if you're testing something
	// that causes them to throw, and want readable tests, this is the line
	// to uncomment.
	//_Ctx.env.Mute(true)
	return
}

// so the timing of other tests/bench's isn't thrown.
func TestJVMFirst(t *testing.T) { setupJVM(t) }


func init(){
	// this only works if you run it directly (not in the gotest framework, as it does not inherit options)
	flag.BoolVar(&squelchExceptions, "squelch-ex", squelchExceptions, "Squelch unexpected exceptions from printing")
}
