package divination

import (
	"encoding/json"
	"time"
	"wecalendarweb/gconst"
	"wecalendarweb/gfunc"
	"wecalendarweb/pb"
	"wecalendarweb/rconst"
	"wecalendarweb/server"

	"github.com/golang/protobuf/proto"
	"github.com/gomodule/redigo/redis"
)

type divinationRecordAddReq struct {
	NowData  string `json:"nowdata"`  // 日期
	Content  string `json:"content"`  // 内容
	Name     string `json:"name"`     // 昵称
	Portrait string `json:"portrait"` // 头像
}

func divinationRecordAddHandle(c *server.StupidContext) {
	log := c.Log.WithField("func", "divination.divinationRecordAddHandle")

	httpRsp := pb.HTTPResponse{
		Result: proto.Int32(int32(gconst.UnknownError)),
	}
	defer c.WriteJSONRsp(&httpRsp)

	// req
	req := &divinationRecordAddReq{}
	if err := json.Unmarshal(c.Body, req); err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("请求信息解析失败")
		log.Errorf("code:%d msg:%s json Unmarshal err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	log.Info("divinationRecordAddHandle enter, req:", string(c.Body))

	conn := c.RedisConn
	// playerid := c.UserID
	nowtime := time.Now()

	// redis multi get
	conn.Send("MULTI")
	conn.Send("GET", rconst.StringDivinationID)
	redisMDArray, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
		httpRsp.Msg = proto.String("统一获取缓存操作失败")
		log.Errorf("code:%d msg:%s redisMDArray Values err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	// do something
	divinationid, _ := redis.Int64(redisMDArray[0], nil)
	divinationid++

	// redis multi set
	conn.Send("MULTI")
	data := rconst.Divination{
		DivinationID: divinationid,
		Time:         nowtime.Unix(),
		Content:      req.Content,
		Name:         req.Name,
		Portrait:     req.Portrait,
		Noname:       false,
	}
	databyte, err := json.Marshal(data)
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("吐槽解析错误")
		log.Errorf("code:%d msg:%s marshal err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}
	conn.Send("HSET", rconst.HashDivinationPrefix+req.NowData, divinationid, databyte)
	conn.Send("ZADD", rconst.ZSetDivinationRecordPrefix+req.NowData, nowtime.Unix(), divinationid)
	conn.Send("SETEX", rconst.StringDivinationID, gfunc.TomorrowZeroRemain(), divinationid)
	_, err = conn.Do("EXEC")
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
		httpRsp.Msg = proto.String("统一存储缓存操作失败")
		log.Errorf("code:%d msg:%s exec err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	httpRsp.Result = proto.Int32(int32(gconst.Success))

	log.Info("divinationRecordAddHandle rsp, result:", httpRsp.GetResult())

	return
}
