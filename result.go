package ging

import (
	"fmt"
)

/* ================================================================================
 * Json结果数据域结构
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	Result struct {
		Code int32
		Msg  string
		Data interface{}
	}

	PagingResult struct {
		Result
		Paging *Paging
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置错误状态信息
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (result *Result) SetError(err error) {
	if customErr, ok := err.(*CustomError); ok {
		result.Code = customErr.Code
		result.Msg = customErr.Msg
	} else {
		msg := fmt.Sprintf("%s", err.Error())
		result.Code = 111
		result.Msg = msg
	}
}
