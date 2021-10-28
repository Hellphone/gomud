package models

type UserState int

// TODO: maybe store states for each command to be executed in
// возможно, для каждой функции нужно хранить состояния, в которых их можно использовать
const (
	NotLoggedIn UserState = iota
	LoggingIn
	LoggedIn
	InGame
)
