package auth

import (
	"weagentweb/gconst"
	"weagentweb/pb"
	"weagentweb/server"
	"weagentweb/tables"

	"github.com/golang/protobuf/proto"
)

func webAddHandle(c *server.StupidContext) {
	log := c.Log.WithField("func", "auth.webAddHandle")

	httpRsp := pb.HTTPResponse{
		Result: proto.Int32(int32(gconst.UnknownError)),
	}
	defer c.WriteJSONRsp(&httpRsp)

	account := c.Query.Get("account")
	password := c.Query.Get("password")
	nick := c.Query.Get("nick")
	protrait := c.Query.Get("protrait")
	key := c.Query.Get("key")

	log.Infof("webAddHandle enter")
	if account == "" || password == "" || nick == "" {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParamNil))
		httpRsp.Msg = proto.String("参数为空")
		log.Errorf("code:%d msg:%s param nil", httpRsp.GetResult(), httpRsp.GetMsg())
		return
	}

	db := c.DbConn

	if key != webkey {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParam))
		httpRsp.Msg = proto.String("兄弟，别试了，回去吃屎吧")
		log.Errorf("code:%d msg:%s webkey err, key:%s", httpRsp.GetResult(), httpRsp.GetMsg(), key)
		return
	}

	getpassword, err := getPassword(password)
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("计算密码失败")
		log.Errorf("code:%d msg:%s getPassword err, err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}
	log.Info("getpassword:", getpassword)

	webaccount := &tables.Webaccount{
		Account:  account,
		Password: getpassword,
		Nick:     nick,
		Gender:   1,
		Portrait: protrait,
		Role:     []string{"admin"},
	}
	_, err = db.Insert(webaccount)
	if err != nil {
		httpRsp.Result = proto.Int32(int32(gconst.ErrParse))
		httpRsp.Msg = proto.String("插入用户信息失败")
		log.Errorf("code:%d msg:%s db insert err:%s", httpRsp.GetResult(), httpRsp.GetMsg(), err.Error())
		return
	}

	httpRsp.Result = proto.Int32(int32(gconst.Success))

	log.Info("webAddHandle rsp, result:", httpRsp.GetResult())

	return
}
