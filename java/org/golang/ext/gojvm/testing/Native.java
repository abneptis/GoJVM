package org.golang.ext.gojvm.testing;

class Native {
	public	native	void 	NativePing();
	public	native	int  	NativeInt();

	public native		void		NativeComplex(Object a, Object b, int i);

	public	native	boolean	NativeBool();
	public	native	short 	NativeShort();
	public	native	long		NativeLong();
	public	native	float 	NativeFloat();
	public	native	double	NativeDouble();
	public	native	String	NativeString();
}
