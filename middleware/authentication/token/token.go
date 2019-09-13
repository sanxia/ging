package token

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
 * Token认证模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	validateHandler     func(ctx *gin.Context, extend TokenExtend, tokenIdentity ging.IToken) bool
	tokenAuthentication struct {
		Validate validateHandler
		Extend   TokenExtend
	}

	TokenExtend struct {
		Roles        []string     //角色（多个之间用逗号分隔）
		Cookie       *ging.Cookie //Cookie
		AuthorizeUrl string       //认证url
		DefaultUrl   string       //认证通过默认返回url
		PassUrls     []string     //直接通过的url
		EncryptKey   string       //加密key
		IsCookie     bool         //是否Cookie模式（0:header | 1: cookie）
		IsRefresh    bool         //是否滑动刷新
		IsEnabled    bool         //是否启用验证
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登入
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (tokenAuth *tokenAuthentication) Logon(ctx *gin.Context, payload *ging.TokenPayload) {
	tokenIdentity := ging.NewToken(tokenAuth.Extend.EncryptKey)
	tokenIdentity.SetPayload(payload)
	tokenIdentity.SetAuthenticated(true)

	tokenAuth.SaveToken(ctx, tokenIdentity, true)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登出
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (tokenAuth *tokenAuthentication) Logoff(ctx *gin.Context) {
	tokenAuth.ClearToken(ctx)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 身份验证
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (tokenAuth *tokenAuthentication) Validation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestPath := ctx.Request.URL.Path
		isPass := strings.HasPrefix(requestPath, tokenAuth.Extend.AuthorizeUrl)
		if !isPass {
			//允许指定模式的Url跳过验证
			for _, passUrl := range tokenAuth.Extend.PassUrls {
				isPass = strings.HasPrefix(requestPath, passUrl)
				if isPass {
					break
				}
			}
		}

		if !tokenAuth.Extend.IsEnabled || isPass {
			//如果未启用认证或跳过url
			var currentToken ging.IToken

			if isPass {
				//默认返回地址处理
				if isReturnUrl := tokenAuth.defaultReturnUrl(ctx); isReturnUrl {
					return
				}

				if tokenIdentity, err := tokenAuth.parseTokenIdentity(ctx); err == nil {
					currentToken = tokenIdentity
				}
			}

			tokenAuth.SaveToken(ctx, currentToken, tokenAuth.Extend.IsRefresh)
		} else {
			if tokenIdentity, err := tokenAuth.parseTokenIdentity(ctx); err == nil {
				if !tokenIdentity.IsAuthenticated() {
					if isSuccess := tokenAuth.Validate(ctx, tokenAuth.Extend, tokenIdentity); !isSuccess {
						tokenAuth.ErrorHandler(ctx)
						return
					}
				}

				tokenAuth.SaveToken(ctx, tokenIdentity, tokenAuth.Extend.IsRefresh)
			} else {
				tokenAuth.ErrorHandler(ctx)
				return
			}
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 刷新Token
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (tokenAuth *tokenAuthentication) SaveToken(ctx *gin.Context, tokenIdentity ging.IToken, isRefresh bool) {
	if tokenIdentity != nil && isRefresh {
		//默认有效期15分钟
		maxAge := 900
		if tokenAuth.Extend.Cookie.MaxAge > 0 {
			maxAge = tokenAuth.Extend.Cookie.MaxAge
		}

		expiresDate := time.Now().Add(time.Duration(maxAge) * time.Second)
		tokenIdentity.SetExpires(expiresDate.Unix())

		//Token标识令牌写入客户端响应
		tokenTicket := tokenIdentity.GetToken()
		if tokenAuth.Extend.IsCookie {
			tokenAuth.SaveTokenCookie(ctx, tokenTicket)
		} else {
			tokenAuth.SaveTokenHeader(ctx, tokenTicket)
		}
	}

	//传递Token标识
	ctx.Set(ging.TOKEN_IDENTITY, tokenIdentity)

	ctx.Next()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 写入Token令牌到客户端响应Cookie
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (tokenAuth *tokenAuthentication) SaveTokenCookie(ctx *gin.Context, tokenTicket string) {
	//默认有效期15分钟
	maxAge := 900
	if tokenAuth.Extend.Cookie.MaxAge > 0 {
		maxAge = tokenAuth.Extend.Cookie.MaxAge
	}

	tokenName := ging.TOKEN_IDENTITY
	if len(tokenAuth.Extend.Cookie.Name) > 0 {
		tokenName = tokenAuth.Extend.Cookie.Name
	}

	path := "/"
	if len(tokenAuth.Extend.Cookie.Path) > 0 {
		path = tokenAuth.Extend.Cookie.Path
	}

	tokenCookie := http.Cookie{
		Name:     tokenName,
		Value:    tokenTicket,
		Path:     path,
		Domain:   tokenAuth.Extend.Cookie.Domain,
		MaxAge:   maxAge,
		HttpOnly: tokenAuth.Extend.Cookie.IsHttpOnly,
		Secure:   tokenAuth.Extend.Cookie.IsSecure,
	}

	http.SetCookie(ctx.Writer, &tokenCookie)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 写入Token令牌到客户端响应头
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (tokenAuth *tokenAuthentication) SaveTokenHeader(ctx *gin.Context, tokenTicket string) {
	tokenName := ging.TOKEN_IDENTITY
	if len(tokenAuth.Extend.Cookie.Name) > 0 {
		tokenName = tokenAuth.Extend.Cookie.Name
	}

	ctx.Header(tokenName, tokenTicket)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 清除Token
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (tokenAuth *tokenAuthentication) ClearToken(ctx *gin.Context) {
	tokenCookieName := ging.TOKEN_IDENTITY
	if len(tokenAuth.Extend.Cookie.Name) > 0 {
		tokenCookieName = tokenAuth.Extend.Cookie.Name
	}

	if tokenAuth.Extend.IsCookie {
		path := "/"
		if len(tokenAuth.Extend.Cookie.Path) > 0 {
			path = tokenAuth.Extend.Cookie.Path
		}

		tokenCookie := http.Cookie{
			Name:     tokenCookieName,
			Value:    "",
			Path:     path,
			Domain:   tokenAuth.Extend.Cookie.Domain,
			MaxAge:   -1,
			HttpOnly: tokenAuth.Extend.Cookie.IsHttpOnly,
			Secure:   tokenAuth.Extend.Cookie.IsSecure,
		}
		http.SetCookie(ctx.Writer, &tokenCookie)
	} else {
		ctx.Header(tokenCookieName, "")
	}

	//清空Token标识
	ctx.Set(ging.TOKEN_IDENTITY, nil)

	ctx.Next()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 错误处理
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (tokenAuth *tokenAuthentication) ErrorHandler(ctx *gin.Context) {
	requestUrl := ctx.Request.URL.RequestURI()
	if requestUrl == "" || requestUrl == "/" || requestUrl == "/#/" || requestUrl == "#/" {
		requestUrl = ctx.Request.URL.RequestURI()
	}

	//认证失败处理
	if ging.IsAjax(ctx) {
		result.JsonResult(ctx, ging.NewError(199, "身份未认证")).Render()
	} else {
		authorizeUrl := fmt.Sprintf("%s?returnurl=%s", tokenAuth.Extend.AuthorizeUrl, glib.UrlEncode(requestUrl))
		result.RedirectResult(ctx, authorizeUrl).Render()
	}

	ctx.Abort()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 默认返回地址处理
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (tokenAuth *tokenAuthentication) defaultReturnUrl(ctx *gin.Context) bool {
	isReturnUrl := false
	if ctx.Request.URL.Path == tokenAuth.Extend.AuthorizeUrl {
		returnUrl := ctx.DefaultQuery("returnurl", "")
		if len(returnUrl) == 0 {
			redirectUrl := fmt.Sprintf("%s?returnurl=%s", ctx.Request.URL.Path, glib.UrlEncode(tokenAuth.Extend.DefaultUrl))
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
func (tokenAuth *tokenAuthentication) parseTokenIdentity(ctx *gin.Context) (ging.IToken, error) {
	tokenIdentity := ging.NewToken(tokenAuth.Extend.EncryptKey)
	tokenName := ging.TOKEN_IDENTITY
	tokenValue := ""

	if len(tokenAuth.Extend.Cookie.Name) > 0 {
		tokenName = tokenAuth.Extend.Cookie.Name
	}

	//从请求头获取 > 从查询参数获取取 > 最后从请求体获取 > 从Cookie获取
	if token, isOk := ctx.Request.Header[tokenName]; !isOk {
		if token, isOk := ctx.Request.URL.Query()[tokenName]; !isOk {
			if token := ctx.PostForm(tokenName); len(token) == 0 {
				if cookieToken, err := ctx.Request.Cookie(tokenName); err == nil {
					tokenValue = cookieToken.Value
				}
			} else {
				tokenValue = token
			}
		} else {
			tokenValue = token[0]
		}
	} else {
		tokenValue = token[0]
	}

	if len(tokenValue) == 0 {
		return nil, errors.New("token identity error")
	}

	if err := tokenIdentity.ParseToken(tokenValue); err != nil {
		return nil, err
	}

	return tokenIdentity, nil
}
