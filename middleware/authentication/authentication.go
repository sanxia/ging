package authentication

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

import (
	"github.com/gin-gonic/gin"
)

import (
	"github.com/sanxia/ging"
	"github.com/sanxia/ging/result"
)

/* ================================================================================
 * 表单认证模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */
type (
	fnValidate          func(ctx *gin.Context, formExtend FormsAuthenticationExtend, userIdentity *ging.UserIdentity) bool
	FormsAuthentication struct {
		Validate fnValidate
		Extend   FormsAuthenticationExtend
	}

	FormsAuthenticationExtend struct {
		EncryptKey        string                     //加密key
		Role              string                     //角色（多个之间用逗号分隔）
		PassUrls          []string                   //直接通过的url
		AuthenticationUrl string                     //认证url
		Cookie            *FormsAuthenticationCookie //cookie
		IsJson            bool                       //是否json响应数据
		IsEnabled         bool                       //是否启用验证
	}

	FormsAuthenticationCookie struct {
		Name   string
		Path   string
		Domain string
		MaxAge int
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 创建新的表单验证实例
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewFormAuthentication(forms FormsAuthentication) (*FormsAuthentication, error) {
	if len(forms.Extend.EncryptKey) != 32 {
		return nil, errors.New("表单认证Key长度必须是32bytes")
	}

	return &forms, nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 身份验证
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (forms *FormsAuthentication) Validation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//允许指定模式的Url跳过验证
		isPass := strings.HasPrefix(ctx.Request.URL.Path, forms.Extend.AuthenticationUrl)
		if !isPass {
			for _, passUrl := range forms.Extend.PassUrls {
				isPass = strings.HasPrefix(ctx.Request.URL.Path, passUrl)
				if isPass {
					break
				}
			}
		}

		if !forms.Extend.IsEnabled || isPass {
			ctx.Next()
			return
		} else if userIdentity, err := forms.parseUserIdentity(ctx); err != nil {
			forms.errorHandler(ctx)
			return
		} else if isSuccess := forms.Validate(ctx, forms.Extend, userIdentity); !isSuccess {
			forms.errorHandler(ctx)
			return
		} else {
			//传递验证用户标识
			ctx.Set(ging.UserIdentityKey, *userIdentity)
			ctx.Next()
		}
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 解析用户标识
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (forms *FormsAuthentication) parseUserIdentity(ctx *gin.Context) (*ging.UserIdentity, error) {
	//解析cookie
	userIdentity := new(ging.UserIdentity)
	if cookie, err := ctx.Request.Cookie(forms.Extend.Cookie.Name); err != nil {
		return nil, err
	} else if err := userIdentity.DecryptAES([]byte(forms.Extend.EncryptKey), cookie.Value); err != nil {
		return nil, err
	}

	return userIdentity, nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 错误处理
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (forms *FormsAuthentication) errorHandler(ctx *gin.Context) {
	errorResult := map[string]interface{}{
		"Code": 299,
		"Msg":  "用户未认证",
		"Data": nil,
	}
	authUrl := forms.Extend.AuthenticationUrl
	requestUrl := ctx.Request.URL.RequestURI()

	//认证失败的处理
	if forms.Extend.IsJson {
		result.JsonResult(ctx, errorResult).Render()
	} else {
		authUrl += "?returnurl=" + url.QueryEscape(requestUrl)
		result.RedirectResult(ctx, authUrl).Render()
	}
	ctx.Abort()
	return
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登陆
 * user: 用户域模型
 * isPersistence: 是否持久化登陆信息
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (forms *FormsAuthentication) Logon(ctx *gin.Context, userIdentity *ging.UserIdentity) bool {
	userIdentityTicket, err := userIdentity.EncryptAES([]byte(forms.Extend.EncryptKey))
	if err != nil {
		return false
	}

	path := "/"
	if len(forms.Extend.Cookie.Path) > 0 {
		path = forms.Extend.Cookie.Path
	}

	cookie := http.Cookie{
		Name:   forms.Extend.Cookie.Name,
		Value:  userIdentityTicket,
		MaxAge: forms.Extend.Cookie.MaxAge,
		Path:   path,
		Domain: forms.Extend.Cookie.Domain,
	}
	http.SetCookie(ctx.Writer, &cookie)

	return true
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登出
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (forms *FormsAuthentication) Logoff(ctx *gin.Context) bool {
	path := "/"
	if len(forms.Extend.Cookie.Path) > 0 {
		path = forms.Extend.Cookie.Path
	}

	cookie := http.Cookie{
		Name:   forms.Extend.Cookie.Name,
		Value:  "_",
		MaxAge: 0,
		Path:   path,
		Domain: forms.Extend.Cookie.Domain,
	}
	http.SetCookie(ctx.Writer, &cookie)
	return true
}
