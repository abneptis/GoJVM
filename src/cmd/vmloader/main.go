package main

import (
	"gojvm"
	"log"
)

func main(){
	print("Initializing VM\n")
 	ctx, err := gojvm.InitializeJVM(0,[]string{".","/usr/lib/jvm/java-6-sun/jre/lib"})
 	if err != nil {
		log.Fatalf("err == %s", err)
	}
	log.Printf("Ctx: %+v", ctx)
	inst, err := ctx.Env.NewStringObject("hello larry")
 	if err != nil {
		log.Fatalf("err == %s", err)
	}
	log.Printf("Instance: %+v", inst)
	i, err := inst.CallInt(false, "length")
	log.Printf("length = %d, err = %+v", i, err)
	o, _, err := inst.CallString(false, "concat", ", it's a pleasure to meet you!")
	log.Printf("toString = %v, err = %+v", o, err)
	
	//strt, _ := gojvm.TypeOf("")
	//ctx.FindMethod(inst, "toString", strt)
}
