package auth

import (
	"encoding/json"
	"strconv"
	"weagentweb/gconst"
	"weagentweb/pb"
	"weagentweb/server"
	"weagentweb/tables"

	"github.com/golang/protobuf/proto"
)

type getinfoRsp struct {
	ID        string   `json:"id"`
	NickName  string   `json:"nickname"`
	Gender    int32    `json:"gender"`
	AvatarURL string   `json:"avatarurl"`
	Role      []string `json:"role"`
}

func webGetinfoHandle(c *server.StupidContext) {
	log := c.Log.WithField("func", "auth.webGetinfoHandle")

	httpRsp := pb.HTTPResponse{
		Result: proto.Int32(int32(gconst.UnknownError)),
	}
	defer c.WriteJSONRsp(&httpRsp)

	log.Info("webGetinfoHandle enter", string(c.Body))

	db := c.DbConn
	playerid := c.UserID

	playeridint, err := strconv.ParseInt(playerid, 10, 64)
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("ID解析失败")
		log.Errorf("code:%d msg:%s ParseInt err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	webaccount := tables.Webaccount{ID: playeridint}
	has, err := db.Get(&webaccount)
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrDB))
		httpRsp.Msg = proto.String("查询账号信息失败")
		log.Errorf("code:%d msg:%s get webaccount err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	if !has {
		httpRsp.Result = proto.Int32(int32(gconst.ErrDB))
		httpRsp.Msg = proto.String("用户不存在")
		log.Errorf("code:%d msg:%s not find account, playerid:%s", httpRsp.GetResult(), httpRsp.GetMsg(), playerid)
		return
	}

	// rsp
	rsp := &getinfoRsp{
		ID:        playerid,
		NickName:  webaccount.Nick,
		Gender:    webaccount.Gender,
		AvatarURL: webaccount.Portrait,
		Role:      webaccount.Role,
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

	log.Info("webGetinfoHandle rsp, rsp:", string(data))

	return
}
