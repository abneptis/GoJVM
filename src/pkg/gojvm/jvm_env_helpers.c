#include "helpers.h"


/* Exception handlers */
jboolean envExceptionCheck(JNIEnv *env) {
  return (*env)->ExceptionCheck(env);
}

jthrowable envExceptionOccurred(JNIEnv *env) {
	return (*env)->ExceptionOccurred(env);
}

void  envExceptionDescribe(JNIEnv* env){ (*env)->ExceptionDescribe(env); }
void  envExceptionClear(JNIEnv* env) { (*env)->ExceptionClear(env); }


/* 'Local' ref handlers */
jobject envNewLocalRef(JNIEnv *env, jobject obj) {
  return (*env)->NewLocalRef(env, obj);
}

void envDeleteLocalRef(JNIEnv *env, jobject obj) {
  (*env)->DeleteLocalRef(env, obj);
}   

jclass envFindClass(JNIEnv *env, const char *string){
	return  (*env)->FindClass(env, string);
}

jmethodID envGetMethodID(JNIEnv *env, jobject jobj, const char *meth, const char *sig){
	return  (*env)->GetMethodID(env, jobj, meth, sig);
}

jmethodID envGetStaticMethodID(JNIEnv *env, jclass jobj, const char *meth, const char *sig){
	return  (*env)->GetStaticMethodID(env, jobj, meth, sig);
}

jclass	envGetObjectClass(JNIEnv *env, jobject jobj){
	return (*env)->GetObjectClass(env, jobj);
}

jint    envCallIntMethodA(JNIEnv *env, jobject o, jmethodID m, void *val){
	return (*env)->CallIntMethodA(env,o,m,val);
}

jint envCallIntMethodV(JNIEnv *env, jclass o, jmethodID m, va_list args){
	return (*env)->CallIntMethodV(env,o,m,args);
}

jobject    envCallObjectMethodA(JNIEnv *env, jobject o, jmethodID m, void *val){
	return (*env)->CallObjectMethodA(env,o,m,val);
}

jobject    envCallStaticObjectMethodA(JNIEnv *env, jclass o, jmethodID m, void *val){
	return (*env)->CallStaticObjectMethodA(env,o,m,val);
}

jint			envCallStaticIntMethodA(JNIEnv *env, jclass o, jmethodID m, void *val){
	return (*env)->CallStaticIntMethodA(env,o,m,val);
}

jboolean	envCallBoolMethodA(JNIEnv *env, jclass o, jmethodID m, void *val){
	return (*env)->CallBooleanMethodA(env,o,m,val);
}

void	envCallStaticVoidMethodA(JNIEnv *env, jclass o, jmethodID m, void *val){
	(*env)->CallStaticVoidMethodA(env,o,m,val);
}

void	envCallStaticVoidMethodV(JNIEnv *env, jclass o, jmethodID m, va_list args){
	(*env)->CallStaticVoidMethodV(env,o,m,args);
}

void	envCallVoidMethodA(JNIEnv *env, jobject o, jmethodID m, void *val){
	(*env)->CallVoidMethodA(env,o,m,val);
}

jint    envGetArrayLength(JNIEnv *env, jobject o){
	return (*env)->GetArrayLength(env,o);
}

jbyte     *envGetByteArrayElements(JNIEnv *env, jobject o, jboolean *b){
	return (*env)->GetByteArrayElements(env,o, b);
}

void      envReleaseByteArrayElements(JNIEnv *env, jobject o, jbyte *bts, jint mode){
	(*env)->ReleaseByteArrayElements(env,o, bts, mode);
}


jboolean  envIsSameObject(JNIEnv *env, jobject o, jobject o2){
	return (*env)->IsSameObject(env,o, o2);
}

jobject	envNewGlobalRef(JNIEnv *env, jobject o){
	return (*env)->NewGlobalRef(env,o);
}

jobject	envNewObjectA(JNIEnv *env, jobject o, jmethodID meth, void *jv){
	return (*env)->NewObjectA(env,o, meth, jv);
}

jobject	envNewObjectALP(JNIEnv *env, jobject o, jmethodID meth, ArgListPtr args){
	return (*env)->NewObjectA(env,o, meth, NULL);
}

jbyteArray  envNewByteArray(JNIEnv *env, jsize len){
	return (*env)->NewByteArray(env,len);
}

jstring   envNewStringUTF(JNIEnv *env, const char *s){
	return (*env)->NewStringUTF(env, s);
}

jsize     envGetStringUTFLength(JNIEnv *env, jstring s){
	return (*env)->GetStringUTFLength(env, s);
}

const char *envGetStringUTFChars(JNIEnv *env, jstring s, jboolean *jb){
	return (*env)->GetStringUTFChars(env, s, jb);
}

void      envReleaseStringUTFChars(JNIEnv *env, jstring s, const char *jb){
	(*env)->ReleaseStringUTFChars(env, s, jb);
}




#include <stdio.h>
void        envSetByteArrayRegion(JNIEnv *env, jbyteArray array, jsize start, jsize len, const void *buf){
	//printf("bufp: %s\n", buf);
	(*env)->SetByteArrayRegion(env, array, start, len, buf);
}

