package ging

import (
	"log"
	"strings"
)

import (
	"github.com/gin-gonic/gin"
)

/* ================================================================================
 * 扩展数据
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 动作过滤器接口
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type (
	IActionFilter interface {
		Before(*gin.Context) IActionResult
		After(*gin.Context)
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 过滤器数据结构
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
type Filter struct {
	Name string
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 判断是否ajax请求
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func IsAjax(ctx *gin.Context) bool {
	var isAjax bool

	//判断是否ajax请求
	xRequestHeader := ctx.Request.Header.Get("x-requested-with")
	if strings.ToLower(xRequestHeader) == "xmlhttprequest" {
		isAjax = true
	}

	log.Printf("IsAjax xRequestHeader: %s", xRequestHeader)

	return isAjax
}
