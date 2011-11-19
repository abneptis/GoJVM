#include "helpers.h"

/* string is duplicated into args, and may be freed after calling, 0 on success. */
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



jint	newJVMContext(JavaVM **jvm, void **env, JavaVMInitArgs *args){
	jint out = JNI_CreateJavaVM(jvm, (void **)env, args);
	return out;
}


