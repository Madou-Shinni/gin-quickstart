package constant

type RspCode int64

const (
	TokenKey = "claims" // 保存在gin上下文中的token key

	FilePathPrefix      = "./static/files/"      // 文件目录前缀
	FileChunkPathPrefix = "./static/fileschunk/" // 文件分片目录前缀

	FileChunkHkey   = "filechunkid:"    // redis HKey 分片文件
	FileChunkHFiled = "filechunkindex:" // redis HFiled 分片文件
)

// 错误码
const (
	CODE_NO_PERMISSIONS    RspCode = 302         // 无权限
	CODE_ERR_MSG           RspCode = 1000 + iota // 消息处理失败
	CODE_ERR_BUSY                                // 系统繁忙
	CODE_INVALID_PARAMETER                       // 无效的参数
	CODE_ADD_FAILED                              // 添加失败
	CODE_DELETE_FAILED                           // 删除失败
	CODE_UPDATE_FAILED                           // 修改失败
	CODE_FIND_FAILED                             // 查询失败
)

var codeMsgMap = map[RspCode]string{
	CODE_ERR_MSG:           "消息处理失败！",
	CODE_ERR_BUSY:          "系统繁忙!",
	CODE_INVALID_PARAMETER: "无效的参数！",
	CODE_NO_PERMISSIONS:    "您没有访问权限！",
	CODE_ADD_FAILED:        "添加失败！",
	CODE_DELETE_FAILED:     "删除失败！",
	CODE_UPDATE_FAILED:     "修改失败！",
	CODE_FIND_FAILED:       "查询失败",
}

func (c RspCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CODE_ERR_BUSY]
	}

	return msg
}
