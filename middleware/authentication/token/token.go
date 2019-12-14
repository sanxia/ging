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
 * Token Auth
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
		Roles         []string    //角色（多个之间用逗号分隔）
		Cookie        ging.Cookie //Cookie
		AuthorizeUrl  string      //认证url
		DefaultUrl    string      //认证通过默认返回url
		PassUrls      []string    //直接通过的url
		EncryptSecret string      //加密秘匙
		IsCookie      bool        //是否Cookie模式（0:header | 1: cookie）
		IsSliding     bool        //是否滑动失效期
		IsDisabled    bool        //是否禁用验证
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * logon
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *tokenAuthentication) Logon(ctx *gin.Context, payload *ging.TokenPayload) {
	//ua finger
	if userAgent, isOk := ctx.Request.Header["User-Agent"]; isOk {
		payload.UserAgent = glib.Sha256(userAgent[0])
	}

	tokenIdentity := ging.NewToken(s.Extend.EncryptSecret)
	tokenIdentity.SetPayload(payload)
	tokenIdentity.SetAuthenticated(true)

	s.setToken(ctx, tokenIdentity)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * logoff
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *tokenAuthentication) Logoff(ctx *gin.Context) {
	s.clearToken(ctx)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * validation
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *tokenAuthentication) Validation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestPath := ctx.Request.URL.Path
		isPass := strings.HasPrefix(requestPath, s.Extend.AuthorizeUrl)

		if !isPass {
			//enable url pass
			for _, passUrl := range s.Extend.PassUrls {
				isPass = strings.HasPrefix(requestPath, passUrl)
				if isPass {
					break
				}
			}
		}

		if s.Extend.IsDisabled || isPass {
			if isPass {
				if isReturnUrl := s.defaultReturnUrl(ctx); isReturnUrl {
					return
				}
			}

			if tokenIdentity, err := s.parseTokenIdentity(ctx); err == nil {
				s.setToken(ctx, tokenIdentity)
			}

			ctx.Next()
		} else {
			if tokenIdentity, err := s.parseTokenIdentity(ctx); err == nil {
				if !tokenIdentity.IsAuthenticated() {
					if isSuccess := s.Validate(ctx, s.Extend, tokenIdentity); !isSuccess {
						s.ErrorHandler(ctx)
						return
					}
				}

				s.setToken(ctx, tokenIdentity)

				ctx.Next()
			} else {
				s.ErrorHandler(ctx)
				return
			}
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * save token
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *tokenAuthentication) setToken(ctx *gin.Context, tokenIdentity ging.IToken) {
	if tokenIdentity != nil && s.Extend.IsSliding {
		//default 15 minutes
		maxAge := 15 * 30
		if s.Extend.Cookie.MaxAge > 0 {
			maxAge = s.Extend.Cookie.MaxAge
		}

		expiresDate := time.Now().Add(time.Duration(maxAge) * time.Second)
		tokenIdentity.SetExpire(expiresDate.Unix())

		//save token to response
		tokenTicket := tokenIdentity.GetToken()
		if s.Extend.IsCookie {
			s.setTokenForCookie(ctx, tokenTicket)
		} else {
			s.setTokenForHeader(ctx, tokenTicket)
		}
	}

	//传递Token标识
	ctx.Set(ging.TOKEN_IDENTITY, tokenIdentity)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * save token to cookie
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *tokenAuthentication) setTokenForCookie(ctx *gin.Context, tokenTicket string) {
	//有效期默认15分钟
	maxAge := 15 * 30
	if s.Extend.Cookie.MaxAge > 0 {
		maxAge = s.Extend.Cookie.MaxAge
	}

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
		Value:    tokenTicket,
		Path:     path,
		Domain:   s.Extend.Cookie.Domain,
		MaxAge:   maxAge,
		HttpOnly: s.Extend.Cookie.IsHttpOnly,
		Secure:   s.Extend.Cookie.IsSecure,
	}

	http.SetCookie(ctx.Writer, &tokenCookie)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * save token to http header
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *tokenAuthentication) setTokenForHeader(ctx *gin.Context, tokenTicket string) {
	tokenName := ging.TOKEN_IDENTITY
	if len(s.Extend.Cookie.Name) > 0 {
		tokenName = s.Extend.Cookie.Name
	}

	ctx.Header(tokenName, tokenTicket)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * claer Token
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *tokenAuthentication) clearToken(ctx *gin.Context) {
	tokenCookieName := ging.TOKEN_IDENTITY
	if len(s.Extend.Cookie.Name) > 0 {
		tokenCookieName = s.Extend.Cookie.Name
	}

	if s.Extend.IsCookie {
		path := "/"
		if len(s.Extend.Cookie.Path) > 0 {
			path = s.Extend.Cookie.Path
		}

		tokenCookie := http.Cookie{
			Name:     tokenCookieName,
			Value:    "",
			Path:     path,
			Domain:   s.Extend.Cookie.Domain,
			MaxAge:   -1,
			HttpOnly: s.Extend.Cookie.IsHttpOnly,
			Secure:   s.Extend.Cookie.IsSecure,
		}
		http.SetCookie(ctx.Writer, &tokenCookie)
	} else {
		ctx.Header(tokenCookieName, "")
	}

	//清空Token标识
	ctx.Set(ging.TOKEN_IDENTITY, nil)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * error process
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *tokenAuthentication) ErrorHandler(ctx *gin.Context) {
	requestUrl := ctx.Request.URL.RequestURI()
	if requestUrl == "" || requestUrl == "/" || requestUrl == "/#/" || requestUrl == "#/" {
		requestUrl = ctx.Request.URL.RequestURI()
	}

	//认证失败处理
	if ging.IsAjax(ctx) {
		result.JsonResult(ctx, ging.NewError(119, "身份未认证")).Render()
	} else {
		authorizeUrl := fmt.Sprintf("%s?returnurl=%s", s.Extend.AuthorizeUrl, glib.UrlEncode(requestUrl))
		result.RedirectResult(ctx, authorizeUrl).Render()
	}

	ctx.Abort()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * defaut return url
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *tokenAuthentication) defaultReturnUrl(ctx *gin.Context) bool {
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
func (s *tokenAuthentication) parseTokenIdentity(ctx *gin.Context) (ging.IToken, error) {
	tokenIdentity := ging.NewToken(s.Extend.EncryptSecret)
	tokenName := ging.TOKEN_IDENTITY
	tokenValue := ""

	if len(s.Extend.Cookie.Name) > 0 {
		tokenName = s.Extend.Cookie.Name
	}

	//req header > query param > form > cookie
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

	//ua
	userAgent := ""
	if userAgentHeaders, isOk := ctx.Request.Header["User-Agent"]; isOk {
		userAgent = glib.Sha256(userAgentHeaders[0])
	}

	if err := tokenIdentity.ParseToken(tokenValue, userAgent); err != nil {
		return nil, err
	}

	return tokenIdentity, nil
}
