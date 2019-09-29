package ging

import (
	"github.com/sanxia/glib"
)

/* ================================================================================
 * Message
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	message struct {
		Payload   *MessagePayload `form:"payload" json:"payload"`     //消息内容
		Type      string          `form:"type" json:"type"`           //消息类型
		Timestamp int64           `form:"timestamp" json:"timestamp"` //unix时间戳
	}

	//消息内容
	MessagePayload struct {
		UserId     string `form:"user_id" json:"user_id"`         //用户id
		TargetId   string `form:"target_id" json:"target_id"`     //目标id
		TargetType string `form:"target_type" json:"target_type"` //目标类型
		Action     string `form:"action" json:"action"`           //动作
		Extend     string `form:"extend" json:"extend"`           //扩展信息
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 初始化Message
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewMessage() *message {
	return &message{
		Payload:   &MessagePayload{},
		Timestamp: glib.UnixTimestamp(),
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 对象转成Json字符串
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *message) ToJson() string {
	jsonString, err := glib.ToJson(s)
	if err != nil {
		return ""
	}

	return jsonString
}
