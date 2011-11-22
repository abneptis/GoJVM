include /usr/share/go/src/Make.inc
TARG=gojvm


JAVA_BASE=../../../java/

DEPS=\
	environment\

CGOFILES=\
	consts.c.go\
	jvm_init_args.c.go\
	jvm_helpers.c.go\

CGO_OFILES=\
	jvm_init_helpers.o\
#	jvm_env_helpers.o\

GOFILES=\
	error.go\

CLEANFILES+=\
	$(JAVA_BASE)/org/golang/ext/gojvm/*.class\
	$(JAVA_BASE)/org/golang/ext/gojvm/testing/*.class\

TESTING_JAVA=\
	$(JAVA_BASE)/org/golang/ext/gojvm/testing/Cleaner.class\
	$(JAVA_BASE)/org/golang/ext/gojvm/testing/Pathos.class\
	$(JAVA_BASE)/org/golang/ext/gojvm/testing/Trivial.class\
	$(JAVA_BASE)/org/golang/ext/gojvm/testing/Native.class\

DIST_JAVA=\
	$(JAVA_BASE)/org/golang/ext/gojvm/Invokable.class\

include /usr/share/go/src/Make.pkg

java_classes: $(TESTING_JAVA) $(DIST_JAVA)

%.class: %.java
	javac $<
