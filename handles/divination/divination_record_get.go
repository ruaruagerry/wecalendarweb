package divination

import (
	"encoding/json"
	"time"
	"wecalendarweb/gconst"
	"wecalendarweb/pb"
	"wecalendarweb/rconst"
	"wecalendarweb/server"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
)

type divinationRecordGetReq struct {
	NowData string `json:"nowdata"`
	Start   int32  `json:"start"`
	End     int32  `json:"end"`
}

type divinationRecordGetItem struct {
	PlayerID     string `json:"playerid"`
	NickName     string `json:"nickname"`
	DivinationID int64  `json:"divinationid"`
	Time         string `json:"time"`
	Content      string `json:"content"`
}

type divinationRecordGetRsp struct {
	Records []*divinationRecordGetItem `json:"records"`
}

func divinationRecordGetHandle(c *server.StupidContext) {
	log := c.Log.WithField("func", "divination.divinationRecordGetHandle")

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

	log.Info("divinationRecordGetHandle enter, req:", string(c.Body))

	rsp := divinationRecordGetRsp{}

	conn := c.RedisConn

	// redis multi get
	conn.Send("MULTI")
	conn.Send("ZRANGE", rconst.ZSetDivinationRecordPrefix+req.NowData, req.Start, req.End)
	redisMDArray, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
		httpRsp.Msg = proto.String("统一获取缓存操作失败")
		log.Errorf("code:%d msg:%s redisMDArray Values err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	// do something
	divinationids, _ := redis.Ints(redisMDArray[0], nil)

	if len(divinationids) > 0 {
		// 先获取吐槽内容
		conn.Send("MULTI")
		for _, v := range divinationids {
			conn.Send("HGET", rconst.HashDivinationPrefix+req.NowData, v)
		}
		redisMDArray, err = redis.Values(conn.Do("EXEC"))
		if err != nil {
			httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
			httpRsp.Msg = proto.String("统一获取缓存操作失败")
			log.Errorf("code:%d msg:%s redisMDArray Values err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
			return
		}

		dslice := [][]byte{}
		for i := range divinationids {
			tmpbytes, _ := redis.Bytes(redisMDArray[i], nil)
			dslice = append(dslice, tmpbytes)
		}

		playerids := []string{}
		for _, v := range dslice {
			tmp := &rconst.Divination{}
			err := json.Unmarshal(v, tmp)
			if err != nil {
				httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
				httpRsp.Msg = proto.String("吐槽解析失败")
				log.Errorf("code:%d msg:%s divination Unmarshal err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
				return
			}

			tmprsp := &divinationRecordGetItem{
				PlayerID:     tmp.PlayerID,
				DivinationID: tmp.DivinationID,
				Time:         time.Unix(tmp.Time, 0).Format("2006-01-02 15:04:05"),
				Content:      tmp.Content,
			}

			rsp.Records = append(rsp.Records, tmprsp)
			playerids = append(playerids, tmp.PlayerID)
		}

		// 再获取玩家信息
		conn.Send("MULTI")
		for _, v := range playerids {
			conn.Send("HGET", rconst.HashAccountPrefix+v, rconst.FieldAccName)
		}
		redisMDArray, err = redis.Values(conn.Do("EXEC"))
		if err != nil {
			httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
			httpRsp.Msg = proto.String("统一获取缓存操作失败")
			log.Errorf("code:%d msg:%s redisMDArray Values err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
			return
		}

		for i := range playerids {
			name, _ := redis.String(redisMDArray[i], nil)
			if name != "" {
				rsp.Records[i].NickName = name
			}
		}
	}

	// rsp
	data, err := json.Marshal(rsp)
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("返回信息marshal解析失败")
		log.Errorf("code:%d msg:%s json marshal err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}
	httpRsp.Result = proto.Int32(int32(gconst.Success))
	httpRsp.Data = data

	log.Info("divinationRecordGetHandle rsp, rsp:", string(data))

	return
}
