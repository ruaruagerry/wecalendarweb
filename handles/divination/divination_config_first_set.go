package divination

import (
	"encoding/json"
	"wecalendarweb/gconst"
	"wecalendarweb/pb"
	"wecalendarweb/rconst"
	"wecalendarweb/server"

	"github.com/golang/protobuf/proto"
)

type divinationConfigFirstSetReq struct {
	First bool `json:"first"`
}

func divinationConfigFirstSetHandle(c *server.StupidContext) {
	log := c.Log.WithField("func", "divination.divinationConfigFirstSetHandle")

	httpRsp := pb.HTTPResponse{
		Result: proto.Int32(int32(gconst.UnknownError)),
	}
	defer c.WriteJSONRsp(&httpRsp)

	// req
	req := &divinationConfigFirstSetReq{}
	if err := json.Unmarshal(c.Body, req); err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("请求信息解析失败")
		log.Errorf("code:%d msg:%s json Unmarshal err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	log.Info("divinationConfigFirstSetHandle enter, req:", string(c.Body))

	conn := c.RedisConn
	// playerid := c.UserID

	// redis multi set
	conn.Send("MULTI")
	conn.Send("HSET", rconst.HashDivinationConfig, rconst.FieldDivinationFirst, req.First)
	_, err := conn.Do("EXEC")
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
		httpRsp.Msg = proto.String("统一存储缓存操作失败")
		log.Errorf("code:%d msg:%s exec err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	httpRsp.Result = proto.Int32(int32(gconst.Success))

	log.Info("divinationConfigFirstSetHandle rsp, result:", httpRsp.GetResult())

	return
}
