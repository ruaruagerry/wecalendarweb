package divination

import (
	"encoding/json"
	"wecalendarweb/gconst"
	"wecalendarweb/gfunc"
	"wecalendarweb/pb"
	"wecalendarweb/rconst"
	"wecalendarweb/server"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
)

type divinationRecordReq struct {
	NowData      string `json:"nowdata"`
	DivinationID int64  `json:"divinationid"`
}

func divinationRecordDelHandle(c *server.StupidContext) {
	log := c.Log.WithField("func", "divination.divinationRecordDelHandle")

	httpRsp := pb.HTTPResponse{
		Result: proto.Int32(int32(gconst.UnknownError)),
	}
	defer c.WriteJSONRsp(&httpRsp)

	// req
	req := &divinationRecordReq{}
	if err := json.Unmarshal(c.Body, req); err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("请求信息解析失败")
		log.Errorf("code:%d msg:%s json Unmarshal err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	log.Info("divinationRecordDelHandle enter, req:", string(c.Body))

	conn := c.RedisConn
	// playerid := c.UserID

	// redis multi get
	// 获取吐槽信息
	conn.Send("MULTI")
	conn.Send("HGET", rconst.HashDivinationPrefix+req.NowData, req.DivinationID)
	redisMDArray, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
		httpRsp.Msg = proto.String("吐槽信息统一获取缓存操作失败")
		log.Errorf("code:%d msg:%s redisMDArray Values err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	divinationbyte, _ := redis.Bytes(redisMDArray[0], nil)

	divination := &rconst.Divination{}
	err = json.Unmarshal(divinationbyte, divination)
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("吐槽解析失败")
		log.Errorf("code:%d msg:%s databyte unmarshal err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	// redis multi set
	conn.Send("MULTI")
	conn.Send("HDET", rconst.HashDivinationPrefix+req.NowData, req.DivinationID)
	conn.Send("ZREM", rconst.ZSetDivinationRecordPrefix+req.NowData, req.DivinationID)
	conn.Send("ZINCRBY", rconst.ZSetDivinationRank, -1, divination.PlayerID)
	conn.Send("EXPIRE", rconst.ZSetDivinationRank, gfunc.TomorrowZeroRemain())
	_, err = conn.Do("EXEC")
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
		httpRsp.Msg = proto.String("统一存储缓存操作失败")
		log.Errorf("code:%d msg:%s exec err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	httpRsp.Result = proto.Int32(int32(gconst.Success))

	log.Info("divinationRecordDelHandle rsp, result:", httpRsp.GetResult())

	return
}
