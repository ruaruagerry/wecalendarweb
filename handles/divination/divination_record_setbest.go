package divination

import (
	"encoding/json"
	"wecalendarweb/gconst"
	"wecalendarweb/pb"
	"wecalendarweb/rconst"
	"wecalendarweb/server"

	"github.com/golang/protobuf/proto"
)

type divinationRecordBestReq struct {
	NowData      string `json:"nowdata"`
	DivinationID int64  `json:"divinationid"`
}

func divinationRecordSetBestHandle(c *server.StupidContext) {
	log := c.Log.WithField("func", "divination.divinationRecordSetBestHandle")

	httpRsp := pb.HTTPResponse{
		Result: proto.Int32(int32(gconst.UnknownError)),
	}
	defer c.WriteJSONRsp(&httpRsp)

	// req
	req := &divinationRecordBestReq{}
	if err := json.Unmarshal(c.Body, req); err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("请求信息解析失败")
		log.Errorf("code:%d msg:%s json Unmarshal err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	log.Info("divinationRecordSetBestHandle enter, req:", string(c.Body))

	conn := c.RedisConn
	// playerid := c.UserID

	// redis multi set
	conn.Send("MULTI")
	conn.Send("SET", rconst.StringDivinationBestPrefix+req.NowData, req.DivinationID)
	_, err := conn.Do("EXEC")
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrRedis))
		httpRsp.Msg = proto.String("统一存储缓存操作失败")
		log.Errorf("code:%d msg:%s exec err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	httpRsp.Result = proto.Int32(int32(gconst.Success))

	log.Info("divinationRecordSetBestHandle rsp, rsp:", httpRsp.GetResult())

	return
}
