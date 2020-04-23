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

type divinationRecordCountReq struct {
	NowData string `json:"nowdata"`
}

type divinationRecordCountRsp struct {
	Count int32 `json:"count"`
}

func divinationRecordCountHandle(c *server.StupidContext) {
	log := c.Log.WithField("func", "divination.divinationRecordCountHandle")

	httpRsp := pb.HTTPResponse{
		Result: proto.Int32(int32(gconst.UnknownError)),
	}
	defer c.WriteJSONRsp(&httpRsp)

	// req
	req := &divinationRecordGetReq{}
	if err := json.Unmarshal(c.Body, req); err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("请求信息解析失败")
		log.Errorf("code:%d msg:%s json Unmarshal err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	log.Info("divinationRecordCountHandle enter, req:", string(c.Body))

	conn := c.RedisConn
	// playerid := c.UserID

	// redis multi get
	conn.Send("MULTI")
	conn.Send("ZCARD", rconst.ZSetDivinationRecordPrefix+req.NowData)
	redisMDArray, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
		httpRsp.Msg = proto.String("统一获取缓存操作失败")
		log.Errorf("code:%d msg:%s redisMDArray Values err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	// do something
	count, _ := redis.Int(redisMDArray[0], nil)

	// rsp
	rsp := &divinationRecordCountRsp{
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

	log.Info("divinationRecordCountHandle rsp, rsp:", string(data))

	return
}
