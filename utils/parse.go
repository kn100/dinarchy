package utils

import (
	"errors"
	"strings"
)

type CreateRequest struct {
	Name       string
	Cronstring string
	Message    string
}

// 0 - name, 1 - message, 2 - cronstring
func ParseCreate(str string) (CreateRequest, error) {
	argsplit := split(str)
	for i, arg := range argsplit {
		argsplit[i] = strings.Trim(arg, " ")
	}
	if strings.Contains(argsplit[0], " ") {
		return CreateRequest{}, errors.New("name cannot contain spaces")
	}
	// TODO: Crons can have commas... this is gonna be nasty
	cr := CreateRequest{Name: argsplit[0], Message: argsplit[1], Cronstring: argsplit[2]}
	return cr, nil
}

func split(cmd string) []string {
	splitFunc := func(c rune) bool {
		return c == ','
	}
	return strings.FieldsFunc(cmd, splitFunc)
}
