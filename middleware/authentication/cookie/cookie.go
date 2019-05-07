package cookie

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/sanxia/ging"
	"github.com/sanxia/ging/result"
	"github.com/sanxia/glib"
)

/* ================================================================================
 * Cookie认证模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	validateHandler      func(ctx *gin.Context, extend CookieExtend, tokenIdentity ging.IToken) bool
	cookieAuthentication struct {
		Validate validateHandler
		Extend   CookieExtend
	}

	CookieExtend struct {
		Roles        []string     //角色（多个之间用逗号分隔）
		Cookie       *ging.Cookie //cookie
		AuthorizeUrl string       //认证url
		DefaultUrl   string       //认证通过默认返回url
		PassUrls     []string     //直接通过的url
		EncryptKey   string       //加密key
		IsRefresh    bool         //是否滑动刷新
		IsEnabled    bool         //是否启用验证
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登入
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (cookieAuth *cookieAuthentication) Logon(ctx *gin.Context, payload *ging.TokenPayload) bool {
	tokenIdentity := ging.NewToken(cookieAuth.Extend.EncryptKey)
	tokenIdentity.SetPayload(payload)
	tokenIdentity.SetAuthenticated(true)

	cookieAuth.SaveCookie(ctx, tokenIdentity, true)

	return true
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登出
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (cookieAuth *cookieAuthentication) Logoff(ctx *gin.Context) {
	cookieAuth.ClearCookie(ctx)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 身份验证
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (cookieAuth *cookieAuthentication) Validation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestPath := ctx.Request.URL.Path
		isPass := strings.HasPrefix(requestPath, cookieAuth.Extend.AuthorizeUrl)
		if !isPass {
			//允许指定模式的Url跳过验证
			for _, passUrl := range cookieAuth.Extend.PassUrls {
				isPass = strings.HasPrefix(requestPath, passUrl)
				if isPass {
					break
				}
			}
		}

		if !cookieAuth.Extend.IsEnabled || isPass {
			//如果未启用认证或跳过url
			var currentToken ging.IToken

			if isPass {
				//默认返回地址处理
				if isReturnUrl := cookieAuth.defaultReturnUrl(ctx); isReturnUrl {
					return
				}

				if tokenIdentity, err := cookieAuth.parseTokenIdentity(ctx); err == nil {
					currentToken = tokenIdentity
				}
			}

			cookieAuth.SaveCookie(ctx, currentToken, cookieAuth.Extend.IsRefresh)
		} else {
			if tokenIdentity, err := cookieAuth.parseTokenIdentity(ctx); err == nil {
				if !tokenIdentity.IsAuthenticated() {
					if isSuccess := cookieAuth.Validate(ctx, cookieAuth.Extend, tokenIdentity); !isSuccess {
						cookieAuth.ErrorHandler(ctx)
						return
					}
				}

				cookieAuth.SaveCookie(ctx, tokenIdentity, cookieAuth.Extend.IsRefresh)
			} else {
				cookieAuth.ErrorHandler(ctx)
				return
			}
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 保存Cookie
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (cookieAuth *cookieAuthentication) SaveCookie(ctx *gin.Context, tokenIdentity ging.IToken, isRefresh bool) {
	if tokenIdentity != nil && isRefresh {
		//默认有效期15分钟
		maxAge := 900
		if cookieAuth.Extend.Cookie.MaxAge > 0 {
			maxAge = cookieAuth.Extend.Cookie.MaxAge
		}

		expiresDate := time.Now().Add(time.Duration(maxAge) * time.Second)
		tokenIdentity.SetExpires(expiresDate.Unix())

		tokenName := ging.USER_IDENTITY
		if len(cookieAuth.Extend.Cookie.Name) > 0 {
			tokenName = cookieAuth.Extend.Cookie.Name
		}

		path := "/"
		if len(cookieAuth.Extend.Cookie.Path) > 0 {
			path = cookieAuth.Extend.Cookie.Path
		}

		tokenCookie := http.Cookie{
			Name:     tokenName,
			Value:    tokenIdentity.GetToken(),
			Path:     path,
			Domain:   cookieAuth.Extend.Cookie.Domain,
			MaxAge:   maxAge,
			HttpOnly: cookieAuth.Extend.Cookie.HttpOnly,
			Secure:   cookieAuth.Extend.Cookie.Secure,
		}

		http.SetCookie(ctx.Writer, &tokenCookie)
	}

	//传递Token标识
	ctx.Set(ging.USER_IDENTITY, tokenIdentity)
	ctx.Next()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 清除Cookie
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (cookieAuth *cookieAuthentication) ClearCookie(ctx *gin.Context) {
	tokenName := ging.USER_IDENTITY
	if len(cookieAuth.Extend.Cookie.Name) > 0 {
		tokenName = cookieAuth.Extend.Cookie.Name
	}

	path := "/"
	if len(cookieAuth.Extend.Cookie.Path) > 0 {
		path = cookieAuth.Extend.Cookie.Path
	}

	tokenCookie := http.Cookie{
		Name:     tokenName,
		Value:    "",
		Path:     path,
		Domain:   cookieAuth.Extend.Cookie.Domain,
		MaxAge:   -1,
		HttpOnly: cookieAuth.Extend.Cookie.HttpOnly,
		Secure:   cookieAuth.Extend.Cookie.Secure,
	}

	http.SetCookie(ctx.Writer, &tokenCookie)

	//清空Token标识
	ctx.Set(ging.USER_IDENTITY, nil)
	ctx.Next()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 错误处理
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (cookieAuth *cookieAuthentication) ErrorHandler(ctx *gin.Context) {
	requestUrl := ctx.Request.URL.RequestURI()
	if requestUrl == "" || requestUrl == "/" || requestUrl == "/#/" || requestUrl == "#/" {
		requestUrl = ctx.Request.URL.RequestURI()
	}

	//认证失败处理
	if ging.IsAjax(ctx) {
		result.JsonResult(ctx, ging.NewError(199, "身份未认证")).Render()
	} else {
		authorizeUrl := fmt.Sprintf("%s?returnurl=%s", cookieAuth.Extend.AuthorizeUrl, glib.UrlEncode(requestUrl))
		result.RedirectResult(ctx, authorizeUrl).Render()
	}

	ctx.Abort()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 默认返回地址处理
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (cookieAuth *cookieAuthentication) defaultReturnUrl(ctx *gin.Context) bool {
	isReturnUrl := false
	if ctx.Request.URL.Path == cookieAuth.Extend.AuthorizeUrl {
		returnUrl := ctx.DefaultQuery("returnurl", "")
		if len(returnUrl) == 0 {
			redirectUrl := fmt.Sprintf("%s?returnurl=%s", ctx.Request.URL.Path, glib.UrlEncode(cookieAuth.Extend.DefaultUrl))
			result.RedirectResult(ctx, redirectUrl).Render()

			ctx.Abort()

			isReturnUrl = true
		}
	}

	return isReturnUrl
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 解析Token标识
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (cookieAuth *cookieAuthentication) parseTokenIdentity(ctx *gin.Context) (ging.IToken, error) {
	tokenIdentity := ging.NewToken(cookieAuth.Extend.EncryptKey)
	tokenName := ging.USER_IDENTITY
	tokenValue := ""

	if len(cookieAuth.Extend.Cookie.Name) > 0 {
		tokenName = cookieAuth.Extend.Cookie.Name
	}

	if cookieToken, err := ctx.Request.Cookie(tokenName); err == nil {
		tokenValue = cookieToken.Value
	}

	if len(tokenValue) == 0 {
		return nil, errors.New("token identity error")
	}

	if err := tokenIdentity.ParseToken(tokenValue); err != nil {
		return nil, err
	}

	return tokenIdentity, nil
}
