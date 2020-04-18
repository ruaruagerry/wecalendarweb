package money

import (
	"encoding/json"
	"weagentweb/gconst"
	"weagentweb/pb"
	"weagentweb/server"
	"weagentweb/tables"

	"github.com/golang/protobuf/proto"
)

type getoutCountRsp struct {
	Count int32 `json:"count"`
}

func getoutCountHandle(c *server.StupidContext) {
	log := c.Log.WithField("func", "money.getoutCountHandle")

	httpRsp := pb.HTTPResponse{
		Result: proto.Int32(int32(gconst.UnknownError)),
	}
	defer c.WriteJSONRsp(&httpRsp)

	log.Info("getoutCountHandle enter")

	db := c.DbConn

	count, err := db.Count(&tables.Getoutrecord{})
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrDB))
		httpRsp.Msg = proto.String("查询提现记录失败")
		log.Errorf("code:%d msg:%s get getoutrecords err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	// rsp
	rsp := &getoutCountRsp{
		Count: int32(count),
	}
	data, err := json.Marshal(rsp)
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("返回信息marshal解析失败")
		log.Errorf("code:%d msg:%s json marshal err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}
	httpRsp.Result = proto.Int32(int32(gconst.Success))
	httpRsp.Data = data

	log.Info("getoutCountHandle rsp, rsp:", string(data))

	return
}
