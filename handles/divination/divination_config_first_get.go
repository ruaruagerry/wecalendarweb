package divination

import (
	"encoding/json"
	"wecalendarweb/gconst"
	"wecalendarweb/pb"
	"wecalendarweb/rconst"
	"wecalendarweb/server"

	"github.com/golang/protobuf/proto"
	"github.com/gomodule/redigo/redis"
)

type divinationConfigFirstGetRsp struct {
	First bool `json:"first"`
}

func divinationConfigFirstGetHandle(c *server.StupidContext) {
	log := c.Log.WithField("func", "divination.divinationConfigFirstGetHandle")

	httpRsp := pb.HTTPResponse{
		Result: proto.Int32(int32(gconst.UnknownError)),
	}
	defer c.WriteJSONRsp(&httpRsp)

	log.Info("divinationConfigFirstGetHandle enter")

	conn := c.RedisConn
	// playerid := c.UserID

	// redis multi get
	conn.Send("MULTI")
	conn.Send("HGET", rconst.HashDivinationConfig, rconst.FieldDivinationFirst)
	redisMDArray, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
		httpRsp.Msg = proto.String("统一获取缓存操作失败")
		log.Errorf("code:%d msg:%s redisMDArray Values err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	// do something
	first, _ := redis.Bool(redisMDArray[0], nil)

	// rsp
	rsp := &divinationConfigFirstGetRsp{
		First: first,
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

	log.Info("divinationConfigFirstGetHandle rsp, rsp:", string(data))

	return
}
