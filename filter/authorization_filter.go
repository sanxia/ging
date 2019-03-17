package filter

import (
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
	"github.com/sanxia/ging/util"
)

/* ================================================================================
 * 授权过滤器
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type authorizationFilter struct {
	ging.Filter
	AuthorizationOption
}

type AuthorizationOption struct {
	Authorization *func(*ging.UserIdentity) bool
	Role          string
	AuthUrl       string
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 实例化授权过滤器
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewAuthorizationFilter(authorizationOption *AuthorizationOption) ging.IActionFilter {
	return &authorizationFilter{
		Filter: ging.Filter{
			Name: "authorization",
		},
		AuthorizationOption: AuthorizationOption{
			Authorization: authorizationOption.Authorization,
			Role:          authorizationOption.Role,
			AuthUrl:       authorizationOption.AuthUrl,
		},
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 动作执行之前
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *authorizationFilter) Before(ctx *gin.Context) ging.IActionResult {
	log.Printf("[%s] Before %v", s.Name, time.Now())

	var actionResult ging.IActionResult

	var userIdentity *ging.UserIdentity
	if identity, isOk := ctx.Get(ging.UserIdentityKey); isOk {
		user := identity.(ging.UserIdentity)
		userIdentity = &user
	}

	//跳转到登录地址
	requestUrl := ctx.Request.URL.RequestURI()

	authUrl := s.AuthUrl
	authUrl += "?returnurl=" + url.QueryEscape(requestUrl)
	isJson := len(s.AuthUrl) == 0

	if userIdentity != nil {
		if s.Authorization != nil {
			authorization := *s.Authorization
			if !authorization(userIdentity) {
				if isJson {
					actionResult = getJsonResult(ctx, ging.NewError(198, "操作权限未被许可"))
				} else {
					actionResult = result.RedirectResult(ctx, authUrl)
				}
			}
		} else {
			isInRole := util.IsInRole(userIdentity.Role, s.Role)
			if !isInRole {
				if isJson {
					actionResult = getJsonResult(ctx, ging.NewError(198, "操作权限未被许可"))
				} else {
					actionResult = result.RedirectResult(ctx, authUrl)
				}
			}
		}
	} else {
		if isJson {
			actionResult = getJsonResult(ctx, ging.NewError(199, "用户未认证"))
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
