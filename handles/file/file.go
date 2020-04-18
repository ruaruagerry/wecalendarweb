package login

import (
	"strings"
	"weagentweb/assets"
	"weagentweb/server"
)

func onServeFile(ctx *server.StupidContext) {
	fileName := "index.html"
	r := ctx.GetHTTPRequest()
	if r.URL.Path != "/" {
		fileName = strings.Replace(r.URL.Path, "/", "", 1)
	}
	data, err := assets.ReadFile(fileName)
	if err != nil {
		ctx.Log.Errorf("ReadFile %s err:%v", "index.html", err)
		ctx.WriteContext(nil)
		return
	}
	ctx.WriteContext(data)
	return
}
