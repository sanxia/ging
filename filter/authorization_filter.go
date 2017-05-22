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
	"github.com/sanxia/ging/util"
)

/* ================================================================================
 * 授权过滤器
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */
type authorizationFilter struct {
	Filter
	AuthorizationOption
}

type AuthorizationOption struct {
	Authorization *func(*ging.UserIdentity) bool
	Role          string
	AuthUrl       string
}

func AuthorizationFilter(authorizationOption *AuthorizationOption) IActionFilter {
	return &authorizationFilter{
		Filter: Filter{
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
func (s *authorizationFilter) Before(ctx *gin.Context) result.IActionResult {
	log.Printf("[%s] Before %v", s.Name, time.Now())

	requestUrl := ctx.Request.URL.RequestURI()
	var actionResult ging.IActionResult
	var userIdentity *ging.UserIdentity
	if identity, isOk := ctx.Get(ging.UserIdentityKey); isOk {
		user := identity.(ging.UserIdentity)
		userIdentity = &user
	}

	//跳转到登录地址
	authUrl := s.AuthUrl
	authUrl += "?returnurl=" + url.QueryEscape(requestUrl)
	if userIdentity != nil {
		//Authorization优先于Role
		if s.Authorization != nil {
			authorization := *s.Authorization
			if !authorization(userIdentity) {
				actionResult = result.RedirectResult(ctx, authUrl)
			}
		} else {
			isInRole := util.IsInRole(userIdentity.Role, s.Role)
			if !isInRole {
				actionResult = result.RedirectResult(ctx, authUrl)
			}
		}
	} else {
		actionResult = result.RedirectResult(ctx, authUrl)
	}

	return actionResult
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 动作执行之后
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *authorizationFilter) After(ctx *gin.Context) {
	log.Printf("[%s] After %v", s.Name, time.Now())
}
