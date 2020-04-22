package login

import (
	"wecalendarweb/server"
)

func init() {
	server.RegisterGetHandleNoUserID("/", onServeFile)
	server.RegisterGetHandleNoUserID("/static/css/:name", onServeFile)
	server.RegisterGetHandleNoUserID("/static/js/:name", onServeFile)
	server.RegisterGetHandleNoUserID("/static/fonts/:name", onServeFile)
}
