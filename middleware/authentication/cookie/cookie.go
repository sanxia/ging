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
		Roles         []string    //角色（多个之间用逗号分隔）
		Cookie        ging.Cookie //cookie
		AuthorizeUrl  string      //认证url
		DefaultUrl    string      //认证通过默认返回url
		PassUrls      []string    //直接通过的url
		EncryptSecret string      //加密秘匙
		IsSliding     bool        //是否滑动失效期
		IsDisabled    bool        //是否禁用验证
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * logon
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *cookieAuthentication) Logon(ctx *gin.Context, payload *ging.TokenPayload) bool {
	tokenIdentity := ging.NewToken(s.Extend.EncryptSecret)
	tokenIdentity.SetPayload(payload)
	tokenIdentity.SetAuthenticated(true)

	s.SaveCookie(ctx, tokenIdentity)

	return true
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * logoff
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *cookieAuthentication) Logoff(ctx *gin.Context) {
	s.ClearCookie(ctx)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * validation
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *cookieAuthentication) Validation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestPath := ctx.Request.URL.Path
		isPass := strings.HasPrefix(requestPath, s.Extend.AuthorizeUrl)
		if !isPass {
			//允许指定模式的Url跳过验证
			for _, passUrl := range s.Extend.PassUrls {
				isPass = strings.HasPrefix(requestPath, passUrl)
				if isPass {
					break
				}
			}
		}

		if s.Extend.IsDisabled || isPass {
			//如果未启用认证或跳过url
			var currentToken ging.IToken

			if isPass {
				//默认返回地址处理
				if isReturnUrl := s.defaultReturnUrl(ctx); isReturnUrl {
					return
				}

				if tokenIdentity, err := s.parseTokenIdentity(ctx); err == nil {
					currentToken = tokenIdentity
				}
			}

			s.SaveCookie(ctx, currentToken)
		} else {
			if tokenIdentity, err := s.parseTokenIdentity(ctx); err == nil {
				if !tokenIdentity.IsAuthenticated() {
					if isSuccess := s.Validate(ctx, s.Extend, tokenIdentity); !isSuccess {
						s.ErrorHandler(ctx)
						return
					}
				}

				s.SaveCookie(ctx, tokenIdentity)
			} else {
				s.ErrorHandler(ctx)
				return
			}
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * save cookie
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *cookieAuthentication) SaveCookie(ctx *gin.Context, tokenIdentity ging.IToken) {
	if tokenIdentity != nil && s.Extend.IsSliding {
		//默认有效期15分钟
		maxAge := 900
		if s.Extend.Cookie.MaxAge > 0 {
			maxAge = s.Extend.Cookie.MaxAge
		}

		expiresDate := time.Now().Add(time.Duration(maxAge) * time.Second)
		tokenIdentity.SetExpires(expiresDate.Unix())

		tokenName := ging.TOKEN_IDENTITY
		if len(s.Extend.Cookie.Name) > 0 {
			tokenName = s.Extend.Cookie.Name
		}

		path := "/"
		if len(s.Extend.Cookie.Path) > 0 {
			path = s.Extend.Cookie.Path
		}

		tokenCookie := http.Cookie{
			Name:     tokenName,
			Value:    tokenIdentity.GetToken(),
			Path:     path,
			Domain:   s.Extend.Cookie.Domain,
			MaxAge:   maxAge,
			HttpOnly: s.Extend.Cookie.HttpOnly,
			Secure:   s.Extend.Cookie.Secure,
		}

		http.SetCookie(ctx.Writer, &tokenCookie)
	}

	//传递Token标识
	ctx.Set(ging.TOKEN_IDENTITY, tokenIdentity)
	ctx.Next()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * clear cookie
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *cookieAuthentication) ClearCookie(ctx *gin.Context) {
	tokenName := ging.TOKEN_IDENTITY
	if len(s.Extend.Cookie.Name) > 0 {
		tokenName = s.Extend.Cookie.Name
	}

	path := "/"
	if len(s.Extend.Cookie.Path) > 0 {
		path = s.Extend.Cookie.Path
	}

	tokenCookie := http.Cookie{
		Name:     tokenName,
		Value:    "",
		Path:     path,
		Domain:   s.Extend.Cookie.Domain,
		MaxAge:   -1,
		HttpOnly: s.Extend.Cookie.HttpOnly,
		Secure:   s.Extend.Cookie.Secure,
	}

	http.SetCookie(ctx.Writer, &tokenCookie)

	ctx.Set(ging.TOKEN_IDENTITY, nil)

	ctx.Next()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * error process
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *cookieAuthentication) ErrorHandler(ctx *gin.Context) {
	requestUrl := ctx.Request.URL.RequestURI()
	if requestUrl == "" || requestUrl == "/" || requestUrl == "/#/" || requestUrl == "#/" {
		requestUrl = ctx.Request.URL.RequestURI()
	}

	if ging.IsAjax(ctx) {
		result.JsonResult(ctx, ging.NewError(199, "身份未认证")).Render()
	} else {
		authorizeUrl := fmt.Sprintf("%s?returnurl=%s", s.Extend.AuthorizeUrl, glib.UrlEncode(requestUrl))
		result.RedirectResult(ctx, authorizeUrl).Render()
	}

	ctx.Abort()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * default return url
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *cookieAuthentication) defaultReturnUrl(ctx *gin.Context) bool {
	isReturnUrl := false
	if ctx.Request.URL.Path == s.Extend.AuthorizeUrl {
		returnUrl := ctx.DefaultQuery("returnurl", "")
		if len(returnUrl) == 0 {
			redirectUrl := fmt.Sprintf("%s?returnurl=%s", ctx.Request.URL.Path, glib.UrlEncode(s.Extend.DefaultUrl))
			result.RedirectResult(ctx, redirectUrl).Render()

			ctx.Abort()

			isReturnUrl = true
		}
	}

	return isReturnUrl
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * parse token
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *cookieAuthentication) parseTokenIdentity(ctx *gin.Context) (ging.IToken, error) {
	tokenIdentity := ging.NewToken(s.Extend.EncryptSecret)
	tokenName := ging.TOKEN_IDENTITY
	tokenValue := ""

	if len(s.Extend.Cookie.Name) > 0 {
		tokenName = s.Extend.Cookie.Name
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
