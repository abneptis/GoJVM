TESTING_JAVA=\
  org/golang/ext/gojvm/testing/Cleaner.class\
  org/golang/ext/gojvm/testing/Native.class\
  org/golang/ext/gojvm/testing/Pathos.class\
  org/golang/ext/gojvm/testing/Trivial.class\

java_classes: $(TESTING_JAVA)
install: java_classes

clean:
	rm -f org/golang/ext/gojvm/testing/*.class

%.class: %.java
	javac $<

