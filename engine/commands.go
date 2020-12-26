package engine

import (
	"fmt"
	"strings"
)

type Command interface {
	Execute(handler IHandler)
}

type PrintCmd struct {
	Message string
}

func (pc PrintCmd) Execute(handler IHandler) {
	fmt.Println(pc.Message)
}

type DeleteCmd struct {
	Str    string
	Symbol byte
}

func (dc DeleteCmd) Execute(handler IHandler) {
	r := strings.NewReplacer(string(dc.Symbol), "")
	newString := r.Replace(dc.Str)
	var printCmd Command = &PrintCmd{newString}
	handler.Post(printCmd)
}
