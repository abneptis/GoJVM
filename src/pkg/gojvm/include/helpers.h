#include</usr/lib/jvm/java-6-sun-1.6.0.26/include/jni.h>
#include<string.h>
#include<stdlib.h>
#include<libio.h>
#include<unistd.h>


#ifndef HELPERS_H
#define HELPERS_H

// VarArgs for variadic
typedef struct {
	int length;
	jvalue *values;
}	ArgList,	*ArgListPtr;

// env wrappers
// exception wrappers

jboolean envExceptionCheck(JNIEnv *);
jthrowable envExceptionOccurred(JNIEnv *);
void	envExceptionDescribe(JNIEnv*);
void	envExceptionClear(JNIEnv*);


// ref calls

jobject envNewLocalRef(JNIEnv *env, jobject ref) ;
void envDeleteLocalRef(JNIEnv *env, jobject obj) ;


jclass 		envFindClass(JNIEnv *, const char *);
jmethodID envGetMethodID(JNIEnv *, jobject, const char *, const char *);
jmethodID envGetStaticMethodID(JNIEnv *env, jclass jobj, const char *meth, const char *sig);

jclass		envGetObjectClass(JNIEnv *, jobject);
jstring		envNewStringUTF(JNIEnv *, const char *);
jsize			envGetStringUTFLength(JNIEnv *, jstring);
const	char	*envGetStringUTFChars(JNIEnv *, jstring, jboolean *);
void			envReleaseStringUTFChars(JNIEnv *, jstring, const char *);



jobjectArray	envNewObjectArray(JNIEnv *env, jsize, jclass, jobject);
void 					envSetObjectArrayElement(JNIEnv *env, jobjectArray array, jsize index, jobject val);

jbyteArray	envNewByteArray(JNIEnv *env, jsize len);
void 				envSetByteArrayRegion(JNIEnv *env, jbyteArray array, jsize start, jsize len, const void *buf); 

jfloat		envCallFloatMethodA(JNIEnv *, jobject, jmethodID, void *);
jfloat		envCallStaticFloatMethodA(JNIEnv *, jobject, jmethodID, void *);

jdouble		envCallDoubleMethodA(JNIEnv *, jobject, jmethodID, void *);
jdouble		envCallStaticDoubleMethodA(JNIEnv *, jobject, jmethodID, void *);

jboolean	envCallBoolMethodA(JNIEnv *, jobject, jmethodID, void *);
jboolean	envCallStaticBoolMethodA(JNIEnv *, jobject, jmethodID, void *);

jshort		envCallShortMethodA(JNIEnv *, jobject, jmethodID, void *);
jshort		envCallStaticShortMethodA(JNIEnv *, jclass, jmethodID, void *);

jint			envCallIntMethodA(JNIEnv *, jobject, jmethodID, void *);
jint			envCallStaticIntMethodA(JNIEnv *, jclass, jmethodID, void *);

jlong			envCallLongMethodA(JNIEnv *, jobject, jmethodID, void *);
jlong			envCallStaticLongMethodA(JNIEnv *, jclass, jmethodID, void *);

jobject		envCallObjectMethodA(JNIEnv *, jobject, jmethodID, void *);
jobject		envCallStaticObjectMethodA(JNIEnv *, jclass, jmethodID, void *);

void			envCallVoidMethodA(JNIEnv *, jobject, jmethodID, void *);
void			envCallStaticVoidMethodA(JNIEnv *, jclass, jmethodID, void *);

jint			envGetArrayLength(JNIEnv *, jobject);
jobject		envNewGlobalRef(JNIEnv *, jobject);

jobject		envNewObjectA(JNIEnv *, jclass, jmethodID, void *);
jobject		envNewObjectALP(JNIEnv *, jclass, jmethodID, ArgListPtr);

jboolean	envIsSameObject(JNIEnv *, jobject, jobject);

jbyte			*envGetByteArrayElements(JNIEnv *, jobject, jboolean *);
void			envReleaseByteArrayElements(JNIEnv *, jobject, jbyte *, jint); 


jvalue getArg(ArgListPtr, int);


// internal helpers
int		addStringArgument(JavaVMInitArgs *args, const char *string);
// vm Calls
// env is actually a void **, but we allow void to make CGo easier
// cleaner solutions welcome! :)
jint	newJVMContext(JavaVM **, void *, JavaVMInitArgs *);
jint  vmAttachCurrentThread(JavaVM *jvm, void *env, void *args);
jint 	vmDetachCurrentThread(JavaVM *jvm);


jint	envGetJavaVM(JNIEnv	*, JavaVM **);
jint  envRegisterNative(JNIEnv *, jclass , char *, char *, int); 
jint  envUnregisterNatives(JNIEnv *, jclass); 



// jvalue conversions
jvalue  boolValue(jboolean v);
jvalue  byteValue(jbyte v);
jvalue  charValue(jchar v);
jvalue  shortValue(jshort v);
jvalue  intValue(jint v);
jvalue  longValue(jlong v);
jvalue  floatValue(jfloat v);
jvalue  doubleValue(jdouble v);
jvalue  objValue(jobject v);



void  voidCallback0(JNIEnv *, jobject, void* );
void  voidCallback1(JNIEnv *, jobject, void* );

void *generified0(JNIEnv *env, jobject obj, ...);
void *generified1(JNIEnv *env, jobject obj, ...);
void *generified2(JNIEnv *env, jobject obj, ...);
void *generified3(JNIEnv *env, jobject obj, ...);
void *generified4(JNIEnv *env, jobject obj, ...);
void *generified5(JNIEnv *env, jobject obj, ...);
void *generified6(JNIEnv *env, jobject obj, ...);
void *generified7(JNIEnv *env, jobject obj, ...);
void *generified8(JNIEnv *env, jobject obj, ...);
void *generified9(JNIEnv *env, jobject obj, ...);


jboolean  valBool(jvalue v) ;
jbyte     valByte(jvalue v) ;
jchar     valChar(jvalue v) ;
jshort    valShort(jvalue v) ;
jint      valInt(jvalue v) ;
jlong     valLong(jvalue v) ;
jfloat    valFloat(jvalue v) ;
jdouble   valDouble(jvalue v) ;
jobject   valObject(jvalue v) ;



#endif
