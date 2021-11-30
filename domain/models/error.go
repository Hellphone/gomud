package models

import "errors"

var ErrorCommandNotFound = errors.New("command not found")
var ErrorUserNotFound = errors.New("user not found")
var ErrorNotEnoughArguments = errors.New("not enough arguments")
