#include "helpers.h"

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

jboolean	valBool(jvalue v) { return v.z; }
jbyte			valByte(jvalue v) { return v.b; }
jchar			valChar(jvalue v) { return v.c; }
jshort		valShort(jvalue v) { return v.s; }
jint			valInt(jvalue v) { return v.i; }
jlong			valLong(jvalue v) { return v.j; }
jfloat		valFloat(jvalue v) { return v.f; }
jdouble		valDouble(jvalue v) { return v.d; }
jobject		valObject(jvalue v) { return v.l; }

