include /usr/share/go/src/Make.inc
TARG=gojvm/environment


CFLAGS=-I../include/
DEPS=../types

JAVA_BASE=../../../java/

#CFILES=jvm_helpers.c\
#
CGOFILES=\
	arglist.c.go\
	environ.c.go\
	object.c.go\
	class.c.go\
	method_sig_helpers.c.go\

CGO_OFILES=\
  jvm_env_helpers.o\
  jvm_value_helpers.o\


GOFILES=\
	method_signatures.go

CLEANFILES+=\


include /usr/share/go/src/Make.pkg

