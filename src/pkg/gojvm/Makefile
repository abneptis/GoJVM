include /usr/share/go/src/Make.inc
TARG=gojvm


JAVA_BASE=../../../java
CGO_CFLAGS=-Iinclude/

DEPS=\
	types\
	$(JAVA_BASE)\

CGOFILES=\
	arglist.c.go\
	environ.c.go\
	globals.c.go\
	object.c.go\
	class.c.go\
	jvm.c.go\
	method_sig_helpers.c.go\

CGO_OFILES=\
  jvm_env_helpers.o\
  jvm_jvm_helpers.o\
  jvm_value_helpers.o\


GOFILES=\
	method_signatures.go\
	callback_descriptor.go\

CLEANFILES+=\

CLEANFILES+=\
	$(JAVA_BASE)/org/golang/ext/gojvm/*.class\
	$(JAVA_BASE)/org/golang/ext/gojvm/testing/*.class\

TESTING_JAVA=\
	$(JAVA_BASE)/org/golang/ext/gojvm/testing/Cleaner.class\
	$(JAVA_BASE)/org/golang/ext/gojvm/testing/Native.class\
	$(JAVA_BASE)/org/golang/ext/gojvm/testing/Pathos.class\
	$(JAVA_BASE)/org/golang/ext/gojvm/testing/Trivial.class\

include /usr/share/go/src/Make.pkg

java_classes: $(TESTING_JAVA) $(DIST_JAVA)
%.class: %.java
	javac	$<
