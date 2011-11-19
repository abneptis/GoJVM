include /usr/share/go/src/Make.inc
TARG=gojvm

#CFILES=jvm_helpers.c\
#
CGOFILES=\
	arglist.c.go\
	class.c.go\
	consts.c.go\
	context.c.go\
	environ.c.go\
	jvm_init_args.c.go\
	jvm_helpers.c.go\
	method_sig_helpers.c.go\
	object.c.go\

CGO_OFILES=\
	jvm_env_helpers.o\
	jvm_init_helpers.o\

GOFILES=\
	class_name.go\
	error.go\
	method_sig.go\

CLEANFILES+=\
	java/org/golang/ext/gojvm/*.class\
	java/org/golang/ext/gojvm/testing/*.class\

TESTING_JAVA=\
	java/org/golang/ext/gojvm/testing/Cleaner.class\
	java/org/golang/ext/gojvm/testing/Pathos.class\
	java/org/golang/ext/gojvm/testing/Trivial.class\

DIST_JAVA=\
	java/org/golang/ext/gojvm/Invokable.class\

include /usr/share/go/src/Make.pkg

java_classes: $(TESTING_JAVA) $(DIST_JAVA)

%.class: %.java
	javac $<

#gojava.o: jvm_helpers.c
#	gcc $(CGO_CFLAGS) $(_CGO_CFLAGS_$(GOARCH)) -fPIC $(CFLAGS) -c $^ -o $@
