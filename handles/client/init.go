package client

import (
	"wecalendarweb/server"
)

func init() {
	server.RegisterGetHandle("/client/config/get", getHandle)
	server.RegisterPostHandle("/client/config/set", setHandle)
}
