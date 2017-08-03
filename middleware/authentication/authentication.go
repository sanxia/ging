package authentication

import (
	"errors"
	"log"
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
	fnValidate          func(ctx *gin.Context, formExtend FormsExtend, userIdentity *ging.UserIdentity) bool
	FormsAuthentication struct {
		Validate fnValidate
		Extend   FormsExtend
	}

	FormsExtend struct {
		Cookie     *FormsCookie //cookie
		Role       string       //角色（多个之间用逗号分隔）
		LogonUrl   string       //认证url
		PassUrls   []string     //直接通过的url
		EncryptKey string       //加密key
		IsJson     bool         //是否json响应
		IsEnabled  bool         //是否启用验证
	}

	FormsCookie struct {
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
	currentUserIdentity := ging.UserIdentity{
		IsAuthenticated: false,
	}

	return func(ctx *gin.Context) {
		//允许指定模式的Url跳过验证
		isPass := strings.HasPrefix(ctx.Request.URL.Path, forms.Extend.LogonUrl)
		if !isPass {
			for _, passUrl := range forms.Extend.PassUrls {
				isPass = strings.HasPrefix(ctx.Request.URL.Path, passUrl)
				if isPass {
					break
				}
			}
		}

		//如果没有启用登陆或跳过url
		if !forms.Extend.IsEnabled || isPass {
			log.Println("authentication url pass")
			if isPass {
				if userIdentity, err := forms.parseUserIdentity(ctx); err == nil {
					currentUserIdentity = *userIdentity
					if userIdentity.UserId > 0 {
						currentUserIdentity.IsAuthenticated = true
					}
				}
			}
			ctx.Set(ging.UserIdentityKey, currentUserIdentity)
			ctx.Next()
			return
		} else {
			log.Println("authentication parseUserIdentity")
			if userIdentity, err := forms.parseUserIdentity(ctx); err != nil {
				log.Printf("authentication parseUserIdentity error: %v", err)
				forms.errorHandler(ctx)
				ctx.Set(ging.UserIdentityKey, currentUserIdentity)
				return
			} else {
				log.Printf("authentication userIdentity.IsAuthenticated: %v", userIdentity.IsAuthenticated)
				if !userIdentity.IsAuthenticated {
					isSuccess := forms.Validate(ctx, forms.Extend, userIdentity)
					log.Printf("authentication Validate isSuccess %v", isSuccess)
					if !isSuccess {
						forms.errorHandler(ctx)
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
	logonUrl := forms.Extend.LogonUrl
	requestUrl := ctx.Request.URL.RequestURI()

	//认证失败的处理
	if forms.Extend.IsJson {
		result.JsonResult(ctx, errorResult).Render()
	} else {
		logonUrl += "?returnurl=" + url.QueryEscape(requestUrl)
		result.RedirectResult(ctx, logonUrl).Render()
	}
	ctx.Abort()
	return
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登陆
 * userIdentity: 用户标示符域模型
 * isRemember: 是否持久化登陆信息
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (forms *FormsAuthentication) Logon(ctx *gin.Context, userIdentity *ging.UserIdentity, isRemember bool) bool {
	userIdentity.IsAuthenticated = true
	ticket, err := userIdentity.EncryptAES([]byte(forms.Extend.EncryptKey))
	if err != nil {
		return false
	}

	path := "/"
	if len(forms.Extend.Cookie.Path) > 0 {
		path = forms.Extend.Cookie.Path
	}

	cookie := http.Cookie{
		Name:   forms.Extend.Cookie.Name,
		Value:  ticket,
		Path:   path,
		Domain: forms.Extend.Cookie.Domain,
	}

	if isRemember {
		cookie.MaxAge = forms.Extend.Cookie.MaxAge
	} else {
		//关闭浏览器即失效
		cookie.MaxAge = 0
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

	//删除cookie
	cookie := http.Cookie{
		Name:   forms.Extend.Cookie.Name,
		Value:  "",
		MaxAge: -1,
		Path:   path,
		Domain: forms.Extend.Cookie.Domain,
	}
	http.SetCookie(ctx.Writer, &cookie)

	return true
}
