package gojvm

import (
	"strings"
)

/* A unified representation (with compare) of both foo.bar and foo/bar class names */
type ClassName []string

func NewClassName(s string) ClassName {
	return ClassName(strings.FieldsFunc(s, func(ch int) bool {
		return (ch == '.' || ch == '/')
	}))
}

func (self ClassName) String() string {
	return "<javaClass:" + self.AsName() + " />"
}

func (self ClassName) AsName() string {
	return strings.Join(self, ".")
}

func (self ClassName) AsPath() string {
	return strings.Join(self, "/")
}

func (self ClassName) Cmp(rhs ClassName) (out int) {
	for i, ls := range self {
		if len(rhs) >= i+1 {
			if ls != rhs[i] {
				if ls < rhs[i] {
					out = -1
				} else {
					out = 1
				}
			}
		} else {
			out = 1
		}
		if out != 0 {
			break
		}
	}
	return
}
