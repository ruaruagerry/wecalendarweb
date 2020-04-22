package divination

import (
	"wecalendarweb/server"
)

func init() {
	server.RegisterGetHandle("/divination/record/get", divinationRecordGetHandle)
	server.RegisterGetHandle("/divination/record/add", divinationRecordAddHandle)
	server.RegisterGetHandle("/divination/record/del", divinationRecordDelHandle)
	server.RegisterGetHandle("/divination/record/setbest", divinationRecordSetBestHandle)
	server.RegisterGetHandle("/divination/getbest", divinationRecordGetBestHandle)
}
