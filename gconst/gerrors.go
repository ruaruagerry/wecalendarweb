package gconst

// Error 错误
type Error int32

/* 返回客户端错误码 */
/* 区间：0-10000 */
const (
	// Success 成功 0
	Success = Error(iota)

	// UnknownError 未知错误 1
	UnknownError = Error(1)
	// ErrParam 参数错误 2
	ErrParam = Error(2)
	// ErrParamNil 请求参数为空 3
	ErrParamNil = Error(3)
	// ErrParse 解析失败 4
	ErrParse = Error(4)
	// ErrDB 数据库操作失败 5
	ErrDB = Error(5)
	// ErrRedis 缓存操作失败 6
	ErrRedis = Error(6)
	// ErrTableConfig 表配置错误 7
	ErrTableConfig = Error(7)
	// ErrTokenEmpty token is empty 8
	ErrTokenEmpty = Error(8)
	// ErrTokenDecrypt token decrypt failed 9
	ErrTokenDecrypt = Error(9)
	// ErrTokenFormat token format is invalid 10
	ErrTokenFormat = Error(10)
	// ErrTokenExpired token expired 11
	ErrTokenExpired = Error(11)
	// ErrCreateUUID 生成uuid失败
	ErrCreateUUID = Error(12)
	// ErrHTTP http请求失败
	ErrHTTP = Error(13)
	// ErrHTTPTooFast http请求太快
	ErrHTTPTooFast = Error(14)
	// ErrPassword 密码错误
	ErrPassword = Error(15)
)

var errMsg = map[Error]string{
	UnknownError:    "未知错误",
	ErrParam:        "参数错误",
	ErrParamNil:     "请求参数为空",
	ErrParse:        "解析失败",
	ErrDB:           "数据库操作失败",
	ErrRedis:        "缓存操作失败",
	ErrTableConfig:  "表配置错误",
	ErrTokenEmpty:   "token为空",
	ErrTokenDecrypt: "token解析失败",
	ErrTokenFormat:  "token格式错误",
	ErrTokenExpired: "token已过期",
	ErrCreateUUID:   "生成uuid失败",
	ErrHTTP:         "http请求失败",
	ErrHTTPTooFast:  "http请求太快",
	ErrPassword:     "密码错误",
}

// String 获得错误码描述信息
func (e Error) String() string {
	v, ok := errMsg[e]
	if !ok {
		return "未定义错误描述"
	}

	return v
}
