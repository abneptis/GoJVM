#include "helpers.h"

/* string is duplicated into args, and may be freed after calling */
int addStringArgument(JavaVMInitArgs *args, const char *string){
	if (args == NULL ){
		return -1;
	}
	int opti = args->nOptions++;
	args->options = realloc(args->options, sizeof(JavaVMInitArgs) * args->nOptions);
	if (args -> options == NULL ) {
		return -1;
	}
	args->options[opti].optionString = strdup(string);
	args->options[opti].extraInfo = NULL;
	if (args->options[opti].optionString == NULL ){
		return -1;
	}
	return 0;
}



jint	newJVMContext(JavaVM **jvm, JNIEnv **env, JavaVMInitArgs *args){
	jint out = JNI_CreateJavaVM(jvm, (void **)env, args);
	return out;
}


/* arglist helpers */
/*
    jboolean z;
    jbyte    b;
    jchar    c;
    jshort   s;
    jint     i;
    jlong    j;
    jfloat   f;
    jdouble  d;
    jobject  l;
*/

// Ease the coercion of values to jvalue from CGO
jvalue	boolValue(jboolean v){ jvalue jv ={.z=v}; return jv; }
jvalue	byteValue(jbyte v){ jvalue jv ={.b=v}; return jv; }
jvalue	charValue(jchar v){ jvalue jv ={.c=v}; return jv; }
jvalue	shortValue(jshort v){ jvalue jv ={.s=v}; return jv; }
jvalue	intValue(jint v){ jvalue jv ={.i=v}; return jv; }
jvalue	longValue(jlong v){ jvalue jv ={.j=v}; return jv; }
jvalue	floatValue(jfloat v){ jvalue jv ={.f=v}; return jv; }
jvalue	doubleValue(jdouble v){ jvalue jv ={.d=v}; return jv; }
jvalue	objValue(jobject v){ jvalue jv ={.l=v}; return jv; }

