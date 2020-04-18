package money

import (
	"encoding/json"
	"time"
	"weagentweb/gconst"
	"weagentweb/pb"
	"weagentweb/rconst"
	"weagentweb/server"
	"weagentweb/tables"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
)

type getoutResultReq struct {
	Rid    int64 `json:"rid"`
	Status int32 `json:"status"`
}

type getoutResultRsp struct {
	ResultTime time.Time `json:"resulttime"`
}

func getoutResultHandle(c *server.StupidContext) {
	log := c.Log.WithField("func", "money.getoutResultHandle")

	httpRsp := pb.HTTPResponse{
		Result: proto.Int32(int32(gconst.UnknownError)),
	}
	defer c.WriteJSONRsp(&httpRsp)

	// req
	req := &getoutResultReq{}
	if err := json.Unmarshal(c.Body, req); err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("请求信息解析失败")
		log.Errorf("code:%d msg:%s json Unmarshal err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	log.Info("getoutResultHandle enter, req:", string(c.Body))

	if req.Status == tables.GetoutStatusReview {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("请求状态错误")
		log.Errorf("code:%d msg:%s status err, status:%d", httpRsp.GetResult(), httpRsp.GetMsg(), req.Status)
		return
	}

	db := c.DbConn
	conn := c.RedisConn

	// do something
	getoutrecord := tables.Getoutrecord{}
	has, err := db.Where("rid = ?", req.Rid).Get(&getoutrecord)
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrDB))
		httpRsp.Msg = proto.String("查询提现记录失败")
		log.Errorf("code:%d msg:%s query getoutrecord err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	if !has {
		httpRsp.Result = proto.Int32(int32(gconst.ErrDB))
		httpRsp.Msg = proto.String("未找到提现记录")
		log.Errorf("code:%d msg:%s not find rid getoutrecord, rid:%s", httpRsp.GetResult(), httpRsp.GetMsg(), req.Rid)
		return
	}

	playerid := getoutrecord.ID

	// 检查
	conn.Send("MULTI")
	conn.Send("SETNX", rconst.StringLockMoneyGetoutResultPrefix+playerid, "1")
	conn.Send("EXPIRE", rconst.StringLockMoneyGetoutResultPrefix+playerid, rconst.LockTime)
	redisMDArray, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
		httpRsp.Msg = proto.String("请求锁获取缓存失败")
		log.Errorf("code:%d msg:%s, GET lock redis data error:(%s)", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}
	locktag, _ := redis.Int(redisMDArray[0], nil)
	if locktag == 0 {
		httpRsp.Result = proto.Int32(int32(gconst.ErrHTTPTooFast))
		httpRsp.Msg = proto.String("请求过于频繁")
		log.Errorf("code:%d msg:%s, request too fast", httpRsp.GetResult(), httpRsp.GetMsg())
		return
	}

	defer func() {
		conn.Do("DEL", rconst.StringLockMoneyGetoutResultPrefix+playerid)
	}()

	// redis multi get
	conn.Send("MULTI")
	conn.Send("HGET", rconst.HashMoneyPrefix+playerid, rconst.FieldMoneyGetout)
	redisMDArray, err = redis.Values(conn.Do("EXEC"))
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
		httpRsp.Msg = proto.String("统一获取缓存操作失败")
		log.Errorf("code:%d msg:%s redisMDArray Values err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	getout, _ := redis.Int64(redisMDArray[0], nil)

	// 更新
	getoutrecord.Status = req.Status
	getoutrecord.ResultTime = time.Now()
	_, err = db.Where("rid = ?", req.Rid).AllCols().Update(&getoutrecord)
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrDB))
		httpRsp.Msg = proto.String("更新提现记录失败")
		log.Errorf("code:%d msg:%s update getoutrecord err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	// 提现成功
	if req.Status == tables.GetoutStatusSuccess {
		getout += getoutrecord.GetoutMoney

		// redis multi set
		conn.Send("MULTI")
		conn.Send("HSET", rconst.HashMoneyPrefix+playerid, rconst.FieldMoneyGetout, getout)
		_, err = conn.Do("EXEC")
		if err != nil {
			httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
			httpRsp.Msg = proto.String("统一存储缓存操作失败")
			log.Errorf("code:%d msg:%s exec err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
			return
		}
	}

	// 通知发钱

	// rsp
	rsp := &getoutResultRsp{
		ResultTime: getoutrecord.ResultTime,
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

	log.Info("getoutResultHandle rsp, result:", httpRsp.GetResult())

	return
}
