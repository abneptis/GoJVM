package gojvm

import (
	"fmt"
)

type Error struct {
	Code	int
	Message	string
}


var ErrUnimplemented  = Error{-400, "Unimplemented functionality"}
var ErrUnknownClass  = Error{-403, "Unknown class"}
var ErrUnknownMethod = Error{-404, "Unknown method"}

func (self Error)String()(string){
	return fmt.Sprintf("(%d) %q", self.Code, self.Message)
}