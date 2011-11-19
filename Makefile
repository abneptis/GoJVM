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

#CLEANFILES+=\
#	tester\
#	tester.o\
#

include /usr/share/go/src/Make.pkg

tester:	../../cmd/vmloader/main.go
	$(GC)	-o	$@.6 $<
	$(LD)	-o $@	$@.6


#gojava.o: jvm_helpers.c
#	gcc $(CGO_CFLAGS) $(_CGO_CFLAGS_$(GOARCH)) -fPIC $(CFLAGS) -c $^ -o $@
