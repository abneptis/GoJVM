package org.golang.ext.gojvm.testing;

class Cleaner {
	class Cleanable {
		Cleaner parent;
		Cleanable(Cleaner daddy){
			parent = daddy;
		}
		protected void finalize() throws Throwable {
			parent.deadKid(this);			
		}
	}
	int deadKids = 0;
	void deadKid(Cleanable kid){
		//System.err.println("Child died");
		deadKids++;
	}
	int getDeadKids(){ return deadKids; }
	
	Cleanable NewChild(){
		return new Cleanable(this);
	}
}