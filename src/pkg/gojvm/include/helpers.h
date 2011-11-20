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


jboolean	envCallBoolMethodA(JNIEnv *, jobject, jmethodID, void *);
jobject		envCallObjectMethodA(JNIEnv *, jobject, jmethodID, void *);
jint			envCallIntMethodA(JNIEnv *, jobject, jmethodID, void *);
jint			envCallIntMethodV(JNIEnv *, jobject, jmethodID, va_list);
void			envCallVoidMethodA(JNIEnv *, jobject, jmethodID, void *);

jobject		envCallStaticObjectMethodA(JNIEnv *, jclass, jmethodID, void *);
jint			envCallStaticIntMethodA(JNIEnv *, jclass, jmethodID, void *);
void			envCallStaticVoidMethodA(JNIEnv *, jclass, jmethodID, void *);

jint			envGetArrayLength(JNIEnv *, jobject);
jobject		envNewGlobalRef(JNIEnv *, jobject);

jobject		envNewObjectA(JNIEnv *, jclass, jmethodID, void *);
jobject		envNewObjectALP(JNIEnv *, jclass, jmethodID, ArgListPtr);

jboolean	envIsSameObject(JNIEnv *, jobject, jobject);

jbyte			*envGetByteArrayElements(JNIEnv *, jobject, jboolean *);
void			envReleaseByteArrayElements(JNIEnv *, jobject, jbyte *, jint); 


// internal helpers
int		addStringArgument(JavaVMInitArgs *args, const char *string);
// vm Calls
// env is actually a void **, but we allow void to make CGo easier
// cleaner solutions welcome! :)
jint	newJVMContext(JavaVM **, void *, JavaVMInitArgs *);
jint  vmAttachCurrentThread(JavaVM *jvm, void *env, void *args);
jint 	vmDetachCurrentThread(JavaVM *jvm);






jvalue  boolValue(jboolean v);
jvalue  byteValue(jbyte v);
jvalue  charValue(jchar v);
jvalue  shortValue(jshort v);
jvalue  intValue(jint v);
jvalue  longValue(jlong v);
jvalue  floatValue(jfloat v);
jvalue  doubleValue(jdouble v);
jvalue  objValue(jobject v);




#endif
