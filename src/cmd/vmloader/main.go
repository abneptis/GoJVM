package main

import (
	"gojvm"
	"log"
	"flag"
)

var jrePath = gojvm.DefaultJREPath
var cpBase  = "."
func init(){
	flag.StringVar(&jrePath, "jre", jrePath, "Path to JRE (classes)")
	flag.StringVar(&cpBase, "cp", cpBase, "(single) ClassPath")
}




func main(){
	flag.Parse()
	if flag.NArg() != 1 {
		/// TODO: brainfart - golang argv[0]?
		log.Fatalf("Expected: %s 'class-with-main'", "vmloader")
	}
	print("Initializing VM\n")
 	_, env, err := gojvm.NewJVM(0, gojvm.JvmConfig{[]string{cpBase,jrePath}})
 	if err != nil {
		log.Fatalf("err == %s", err)
	}
	klass := flag.Arg(0)
	inst, err := env.GetClassStr(klass)

 	if err != nil {
		log.Fatalf("Couldn't find class %s", err)
	}
	log.Printf("Got instance: %+v", inst)
	err = inst.CallVoid(env, true, "main", []string{})
 	if err != nil {
		log.Fatalf("Couldn't call main: %v", err)
	}
	
	//strt, _ := gojvm.TypeOf("")
	//ctx.FindMethod(inst, "toString", strt)
}
