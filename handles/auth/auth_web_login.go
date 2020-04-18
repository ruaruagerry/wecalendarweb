package auth

import (
	"encoding/json"
	"fmt"
	"weagentweb/gconst"
	"weagentweb/pb"
	"weagentweb/server"
	"weagentweb/tables"

	"github.com/golang/protobuf/proto"
)

type webLoginReq struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type webLoginRsp struct {
	Token string `json:"token"`
}

func webLoginHandle(c *server.StupidContext) {
	log := c.Log.WithField("func", "auth.webLoginHandle")

	httpRsp := pb.HTTPResponse{
		Result: proto.Int32(int32(gconst.UnknownError)),
	}
	defer c.WriteJSONRsp(&httpRsp)

	// req
	req := &webLoginReq{}
	if err := json.Unmarshal(c.Body, req); err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("请求信息解析失败")
		log.Errorf("code:%d msg:%s json Unmarshal err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	log.Info("webLoginHandle enter, req:", string(c.Body))

	if req.Password == "" || req.Account == "" {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParamNil))
		httpRsp.Msg = proto.String("请求参数为空")
		log.Errorf("code:%d msg:%s req param nil", httpRsp.GetResult(), httpRsp.GetMsg())
		return
	}

	// 获得当前密码
	password, err := getPassword(req.Password)
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("计算密码失败")
		log.Errorf("code:%d msg:%s getPassword err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	db := c.DbConn

	webaccount := tables.Webaccount{Account: req.Account}
	has, err := db.Get(&webaccount)
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrDB))
		httpRsp.Msg = proto.String("查询账号信息失败")
		log.Errorf("code:%d msg:%s get webaccount err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	if !has {
		httpRsp.Result = proto.Int32(int32(gconst.ErrDB))
		httpRsp.Msg = proto.String("账号不存在")
		log.Errorf("code:%d msg:%s not find account, account:%s", httpRsp.GetResult(), httpRsp.GetMsg(), req.Account)
		return
	}

	// do something
	if password != webaccount.Password {
		httpRsp.Result = proto.Int32(int32(gconst.ErrPassword))
		httpRsp.Msg = proto.String("密码错误")
		log.Errorf("code:%d msg:%s password err", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	playerid := fmt.Sprintf("%d", webaccount.ID)

	// 生成token， 根据目前客户端的约定需要设置到header上
	token := server.GenTK(playerid)

	// rsp
	rsp := &webLoginRsp{
		Token: token,
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

	log.Info("webLoginHandle rsp, rsp:", string(data))

	return
}
