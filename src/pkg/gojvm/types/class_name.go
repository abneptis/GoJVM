package types

import (
	"strings"
)

/* A unified representation (with compare) of both foo.bar and foo/bar class names */
type Name []string

var JavaLangString	=	Name{"java","lang","String"}
var JavaLangObject	=	Name{"java","lang","Object"}
var JavaLangThrowable	=	Name{"java","lang","Throwable"}

/*
	Parse a string by splitting on '.' and '/';  $Refs are left intact;
	returns a new Name instance
*/
func NewName(s string) Name {
	return Name(strings.FieldsFunc(s, func(ch int) bool {
		return (ch == '.' || ch == '/')
	}))
}

/* 
	returns a generalized (but deliberately machine-useless) string representation
	of a Name
*/
func (self Name) String() string {
	return "<javaClass:" + self.AsName() + " />"
}

// returns a '.' joined representation of the className
func (self Name) AsName() string {
	return strings.Join(self, ".")
}

// returns a '/' joined representation of the className
func (self Name) AsPath() string {
	return strings.Join(self, "/")
}

// returns {-1,0,1} for a comparsion of two classnames.
func (self Name) Cmp(rhs Name) (out int) {
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
