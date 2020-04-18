package auth

import "weagentweb/server"

func init() {
	server.RegisterPostHandleNoUserID("/auth/web/login", webLoginHandle)
	server.RegisterGetHandleNoUserID("/auth/web/add", webAddHandle)
	server.RegisterGetHandle("/auth/web/getinfo", webGetinfoHandle)
}
