package money

import (
	"weagentweb/server"
)

func init() {
	server.RegisterPostHandle("/money/getout/record", getoutRecordHandle)             // 查看所有提现记录
	server.RegisterPostHandle("/money/getout/result", getoutResultHandle)             // 提现审核
	server.RegisterPostHandle("/money/getout/playerrecord", getoutPlayerRecordHandle) // 查找玩家提现记录
	server.RegisterGetHandle("/money/getout/count", getoutCountHandle)                // 获取提现总数
}
