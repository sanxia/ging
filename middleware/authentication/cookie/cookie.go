package cookie

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

import (
	"github.com/gin-gonic/gin"
)

import (
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
	fnValidate           func(ctx *gin.Context, extend CookieExtend, userIdentity *ging.UserIdentity) bool
	CookieAuthentication struct {
		Validate fnValidate
		Extend   CookieExtend
	}

	CookieExtend struct {
		Option       *CookieOption //forms cookie
		Role         string        //角色（多个之间用逗号分隔）
		AuthorizeUrl string        //身份授权url
		DefaultUrl   string        //认证通过默认返回url
		PassUrls     []string      //直接通过的url
		EncryptKey   string        //加密key
		IsEnabled    bool          //是否启用验证
	}

	CookieOption struct {
		Name     string
		Path     string
		Domain   string
		MaxAge   int
		HttpOnly bool
		Secure   bool
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 创建新的Cookie验证实例
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewCookieAuthentication(cookieAuth CookieAuthentication) (*CookieAuthentication, error) {
	if len(cookieAuth.Extend.EncryptKey) != 32 {
		return nil, errors.New("Cookie认证Key长度必须是32bytes")
	}

	return &cookieAuth, nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登陆
 * userIdentity: 用户标示符域模型
 * isRemember: 是否持久化登陆信息
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (cookieAuth *CookieAuthentication) Logon(ctx *gin.Context, userIdentity *ging.UserIdentity, isRemember bool) bool {
	userIdentity.IsAuthenticated = true
	ticket, err := userIdentity.EncryptAES([]byte(cookieAuth.Extend.EncryptKey))
	if err != nil {
		return false
	}

	path := "/"
	if len(cookieAuth.Extend.Option.Path) > 0 {
		path = cookieAuth.Extend.Option.Path
	}

	httpCookie := http.Cookie{
		Name:     cookieAuth.Extend.Option.Name,
		Value:    ticket,
		Path:     path,
		Domain:   cookieAuth.Extend.Option.Domain,
		HttpOnly: cookieAuth.Extend.Option.HttpOnly,
		Secure:   cookieAuth.Extend.Option.Secure,
	}

	if isRemember {
		httpCookie.MaxAge = cookieAuth.Extend.Option.MaxAge
	} else {
		httpCookie.MaxAge = 0
	}

	http.SetCookie(ctx.Writer, &httpCookie)

	return true
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登出
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (cookieAuth *CookieAuthentication) Logoff(ctx *gin.Context) {
	path := "/"
	if len(cookieAuth.Extend.Option.Path) > 0 {
		path = cookieAuth.Extend.Option.Path
	}

	//删除cookie
	httpCookie := http.Cookie{
		Name:     cookieAuth.Extend.Option.Name,
		Value:    "",
		MaxAge:   -1,
		Path:     path,
		Domain:   cookieAuth.Extend.Option.Domain,
		HttpOnly: cookieAuth.Extend.Option.HttpOnly,
		Secure:   cookieAuth.Extend.Option.Secure,
	}
	http.SetCookie(ctx.Writer, &httpCookie)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 身份验证
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (cookieAuth *CookieAuthentication) Validation() gin.HandlerFunc {
	currentUserIdentity := ging.UserIdentity{
		IsAuthenticated: false,
	}

	return func(ctx *gin.Context) {
		//允许指定模式的Url跳过验证
		isPass := strings.HasPrefix(ctx.Request.URL.Path, cookieAuth.Extend.AuthorizeUrl)
		if !isPass {
			for _, passUrl := range cookieAuth.Extend.PassUrls {
				isPass = strings.HasPrefix(ctx.Request.URL.Path, passUrl)
				if isPass {
					break
				}
			}
		}

		//如果未启用认证或跳过url
		if !cookieAuth.Extend.IsEnabled || isPass {
			log.Println("authentication url pass")
			if isPass {
				//默认返回地址处理
				if isReturnUrl := cookieAuth.defaultReturnUrl(ctx); isReturnUrl {
					return
				}

				//解析用户标识
				if userIdentity, err := cookieAuth.parseUserIdentity(ctx); err == nil {
					currentUserIdentity = *userIdentity
					log.Println("authentication url pass currentUserIdentity: %v", currentUserIdentity)
					if len(userIdentity.UserId) > 0 {
						currentUserIdentity.IsAuthenticated = true
					}
				} else {
					currentUserIdentity.IsAuthenticated = false
				}
			}
			ctx.Set(ging.UserIdentityKey, currentUserIdentity)
			ctx.Next()
			return
		} else {
			//解析用户标识
			if userIdentity, err := cookieAuth.parseUserIdentity(ctx); err != nil {
				log.Printf("authentication parseUserIdentity error: %v", err)
				cookieAuth.errorHandler(ctx)
				ctx.Set(ging.UserIdentityKey, currentUserIdentity)
				return
			} else {
				log.Printf("authentication userIdentity: %v", userIdentity)
				if !userIdentity.IsAuthenticated {
					isSuccess := cookieAuth.Validate(ctx, cookieAuth.Extend, userIdentity)
					log.Printf("authentication Validate isSuccess %v", isSuccess)
					if !isSuccess {
						cookieAuth.errorHandler(ctx)
						return
					}
				}

				//传递验证用户标识
				log.Printf("authentication Set UserIdentityKey: %v", userIdentity)
				ctx.Set(ging.UserIdentityKey, *userIdentity)
				ctx.Next()
			}
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 解析用户标识
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (cookieAuth *CookieAuthentication) parseUserIdentity(ctx *gin.Context) (*ging.UserIdentity, error) {
	//解析Cookie
	userIdentity := new(ging.UserIdentity)
	if httpCookie, err := ctx.Request.Cookie(cookieAuth.Extend.Option.Name); err != nil {
		return nil, err
	} else if err := userIdentity.DecryptAES([]byte(cookieAuth.Extend.EncryptKey), httpCookie.Value); err != nil {
		return nil, err
	}

	return userIdentity, nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 默认返回地址处理
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (cookieAuth *CookieAuthentication) defaultReturnUrl(ctx *gin.Context) bool {
	isReturnUrl := false
	if ctx.Request.URL.Path == cookieAuth.Extend.AuthorizeUrl {
		returnUrl := ctx.DefaultQuery("returnurl", "")
		if returnUrl == "" {
			redirectUrl := fmt.Sprintf("%s?returnurl=%s", ctx.Request.URL.Path, glib.UrlEncode(cookieAuth.Extend.DefaultUrl))
			log.Println("authentication defaultReturnUrl: ", redirectUrl)

			result.RedirectResult(ctx, redirectUrl).Render()
			ctx.Abort()

			isReturnUrl = true
		}
	}

	return isReturnUrl
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 错误处理
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (cookieAuth *CookieAuthentication) errorHandler(ctx *gin.Context) {
	returnUrl := ctx.Request.URL.RequestURI()
	if returnUrl == "" || returnUrl == "/" || returnUrl == "/#/" || returnUrl == "#/" {
		returnUrl = ctx.Request.URL.RequestURI()
	}

	//认证失败的处理
	if ging.IsAjax(ctx) {
		result.JsonResult(ctx, ging.NewError(199, "身份未认证")).Render()
	} else {
		redirectUrl := fmt.Sprintf("%s?returnurl=%s", cookieAuth.Extend.AuthorizeUrl, glib.UrlEncode(returnUrl))
		result.RedirectResult(ctx, redirectUrl).Render()
	}

	ctx.Abort()

	return
}
