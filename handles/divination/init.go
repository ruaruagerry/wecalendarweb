package divination

import (
	"wecalendarweb/server"
)

func init() {
	server.RegisterPostHandle("/divination/record/get", divinationRecordGetHandle)
	server.RegisterPostHandle("/divination/record/add", divinationRecordAddHandle)
	server.RegisterPostHandle("/divination/record/del", divinationRecordDelHandle)
	server.RegisterPostHandle("/divination/record/setbest", divinationRecordSetBestHandle)
	server.RegisterPostHandle("/divination/getbest", divinationRecordGetBestHandle)
	server.RegisterPostHandle("/divination/record/count", divinationRecordCountHandle)
	server.RegisterPostHandle("/divination/config/first/set", divinationConfigFirstSetHandle)
	server.RegisterGetHandle("/divination/config/first/get", divinationConfigFirstGetHandle)
}
