package types

import (
	"strings"
)

/* A unified representation (with compare) of both foo.bar and foo/bar class names */
type ClassName []string

/*
	Parse a string by splitting on '.' and '/';  $Refs are left intact;
	returns a new ClassName instance
*/
func NewClassName(s string) ClassName {
	return ClassName(strings.FieldsFunc(s, func(ch int) bool {
		return (ch == '.' || ch == '/')
	}))
}

/* 
	returns a generalized (but deliberately machine-useless) string representation
	of a ClassName
*/
func (self ClassName) String() string {
	return "<javaClass:" + self.AsName() + " />"
}

// returns a '.' joined representation of the className
func (self ClassName) AsName() string {
	return strings.Join(self, ".")
}

// returns a '/' joined representation of the className
func (self ClassName) AsPath() string {
	return strings.Join(self, "/")
}

// returns {-1,0,1} for a comparsion of two classnames.
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
			out = -1
		}
		if out != 0 {
			break
		}
	}
	return
}
