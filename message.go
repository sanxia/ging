package ging

import (
	"github.com/sanxia/glib"
)

/* ================================================================================
 * Message数据域结构
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	Message struct {
		Payload   string //消息内容
		Type      string //类型码
		Timestamp int64  //unix时间戳，单位秒
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 对象转成Json字符串
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (message *Message) ToJson() string {
	jsonString, err := glib.ToJson(message)
	if err != nil {
		return ""
	}

	return jsonString
}
