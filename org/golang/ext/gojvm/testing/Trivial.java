package org.golang.ext.gojvm.testing;

class Trivial {
	String	ConstructorUsed;
	Trivial(){
		ConstructorUsed = new String("()V");
	}
	Trivial(int i){
		ConstructorUsed = new String("(I)V");
	}
	Trivial(long i){
		ConstructorUsed = new String("(J)V");
	}
	Trivial(String i){
		ConstructorUsed = new String("(Ljava/lang/String;)V");
	}

	String getConstructorUsed(){
		return this.ConstructorUsed;
	}
}
