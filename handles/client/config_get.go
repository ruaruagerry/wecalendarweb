package client

import (
	"encoding/json"
	"wecalendarweb/gconst"
	"wecalendarweb/pb"
	"wecalendarweb/rconst"
	"wecalendarweb/server"

	"github.com/golang/protobuf/proto"
	"github.com/gomodule/redigo/redis"
)

type getRsp struct {
	Version string `json:"version"`
	Ad      bool   `json:"ad"`
}

func getHandle(c *server.StupidContext) {
	log := c.Log.WithField("func", "client.getHandle")

	httpRsp := pb.HTTPResponse{
		Result: proto.Int32(int32(gconst.UnknownError)),
	}
	defer c.WriteJSONRsp(&httpRsp)

	log.Info("getHandle enter")

	conn := c.RedisConn
	// playerid := c.UserID

	// redis multi get
	conn.Send("MULTI")
	conn.Send("HGET", rconst.HashClient, rconst.FieldClientLastestVersion)
	conn.Send("HGET", rconst.HashClient, rconst.FieldClientAd)
	redisMDArray, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
		httpRsp.Msg = proto.String("统一获取缓存操作失败")
		log.Errorf("code:%d msg:%s redisMDArray Values err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	lastestversion, _ := redis.String(redisMDArray[0], nil)
	ad, _ := redis.Bool(redisMDArray[1], nil)

	// rsp
	rsp := &getRsp{
		Version: lastestversion,
		Ad:      ad,
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

	log.Info("getHandle rsp, rsp:", string(data))

	return
}
