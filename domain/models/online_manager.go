package models

import (
	"context"
	"net"
)

type OnlineManager struct {
	Connection net.Conn
	Context    context.Context
	User       *User
}
