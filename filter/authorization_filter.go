package filter

import (
	"fmt"
	"log"
	"net/url"
	"time"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/sanxia/ging"
	"github.com/sanxia/ging/result"
)

/* ================================================================================
 * 授权过滤器
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	IAuthorization interface {
		Authorize(ging.IToken) bool
	}

	authorizationFilter struct {
		ging.Filter
		AuthorizationOption
	}

	AuthorizationOption struct {
		Authorization IAuthorization
		AuthorizeUrl  string
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 实例化授权过滤器
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewAuthorizationFilter(option *AuthorizationOption) ging.IActionFilter {
	return &authorizationFilter{
		Filter: ging.Filter{
			Name: "authorization_filter",
		},
		AuthorizationOption: AuthorizationOption{
			Authorization: option.Authorization,
			AuthorizeUrl:  option.AuthorizeUrl,
		},
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 动作执行之前
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *authorizationFilter) Before(ctx *gin.Context) ging.IActionResult {
	log.Printf("[%s] Before %v", s.Name, time.Now())

	var userToken ging.IToken
	var actionResult ging.IActionResult

	//判断是否ajax请求
	isAjax := ging.IsAjax(ctx)

	//获取当前用户标识
	if userIdentity, isOk := ctx.Get(ging.USER_IDENTITY); userIdentity != nil && isOk {
		if tokenIdentity, isOk := userIdentity.(*ging.Token); isOk {
			userToken = tokenIdentity
		}
	}

	//跳转到认证地址
	requestUrl := ctx.Request.URL.RequestURI()
	authUrl := fmt.Sprintf("%s?returnurl=%s", s.AuthorizeUrl, url.QueryEscape(requestUrl))

	if userToken != nil {
		if s.Authorization != nil {
			if !s.Authorization.Authorize(userToken) {
				if isAjax {
					actionResult = getJsonResult(ctx, ging.NewError(191, "操作权限未被许可"))
				} else {
					actionResult = result.RedirectResult(ctx, authUrl)
				}
			}
		}
	} else {
		if isAjax {
			actionResult = getJsonResult(ctx, ging.NewError(199, "身份未认证"))
		} else {
			actionResult = result.RedirectResult(ctx, authUrl)
		}
	}

	return actionResult
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 动作执行之后
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *authorizationFilter) After(ctx *gin.Context) {
	log.Printf("[%s] After %v", s.Name, time.Now())
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取json结果
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func getJsonResult(ctx *gin.Context, err *ging.CustomError) ging.IActionResult {
	jsonResult := new(ging.Result)
	jsonResult.SetError(err)
	return result.JsonResult(ctx, jsonResult, true)
}
