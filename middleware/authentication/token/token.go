package token

import (
	"errors"
	"log"
	"strings"
	"time"
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
 * Token认证模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
const (
	DefaultTokenName string = "__sx_token___"
)

type (
	fnValidate          func(ctx *gin.Context, extend TokenExtend, userIdentity *ging.UserIdentity) bool
	TokenAuthentication struct {
		Validate fnValidate
		Extend   TokenExtend
	}

	TokenExtend struct {
		Option       *TokenOption
		Role         string   //角色（多个之间用逗号分隔）
		AuthorizeUrl string   //认证url
		PassUrls     []string //直接通过的url
		EncryptKey   string   //加密key
		IsEnabled    bool     //是否启用验证
	}

	TokenOption struct {
		Name   string //token名称
		MaxAge int    //最大有效时长，单位：秒
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 创建新的Token验证实例
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewTokenAuthentication(tokenAuth TokenAuthentication) (*TokenAuthentication, error) {
	if len(tokenAuth.Extend.EncryptKey) != 32 {
		return nil, errors.New("表单认证Key长度必须是32bytes")
	}

	return &tokenAuth, nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 身份验证
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (tokenAuth *TokenAuthentication) Validation() gin.HandlerFunc {
	currentUserIdentity := ging.UserIdentity{
		UserId:          0,
		IsAuthenticated: false,
	}

	return func(ctx *gin.Context) {
		//允许指定模式的Url跳过验证
		isPass := strings.HasPrefix(ctx.Request.URL.Path, tokenAuth.Extend.AuthorizeUrl)
		if !isPass {
			for _, passUrl := range tokenAuth.Extend.PassUrls {
				isPass = strings.HasPrefix(ctx.Request.URL.Path, passUrl)
				if isPass {
					break
				}
			}
		}

		//如果未启用认证或跳过url
		if !tokenAuth.Extend.IsEnabled || isPass {
			log.Println("authentication url pass")
			if isPass {
				if userIdentity, err := tokenAuth.parseUserIdentity(ctx); err == nil {
					currentUserIdentity = *userIdentity
					log.Println("authentication url pass currentUserIdentity: %v", currentUserIdentity)
					if userIdentity.UserId > 0 {
						currentUserIdentity.IsAuthenticated = true
					}
				} else {
					currentUserIdentity.UserId = 0
					currentUserIdentity.IsAuthenticated = false
				}
			}
			ctx.Set(ging.UserIdentityKey, currentUserIdentity)
			ctx.Next()
			return
		} else {
			if userIdentity, err := tokenAuth.parseUserIdentity(ctx); err != nil {
				log.Printf("authentication parseUserIdentity error: %v", err)
				tokenAuth.errorHandler(ctx)
				ctx.Set(ging.UserIdentityKey, currentUserIdentity)
				return
			} else {
				log.Printf("authentication userIdentity: %v", userIdentity)

				//时间是否过期
				if time.Now().After(time.Unix(userIdentity.Expires, 0)) {
					tokenAuth.errorHandler(ctx)
					return
				}

				if !userIdentity.IsAuthenticated {
					isSuccess := tokenAuth.Validate(ctx, tokenAuth.Extend, userIdentity)
					log.Printf("authentication Validate isSuccess %v", isSuccess)

					if !isSuccess {
						tokenAuth.errorHandler(ctx)
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
func (tokenAuth *TokenAuthentication) parseUserIdentity(ctx *gin.Context) (*ging.UserIdentity, error) {
	userIdentity := new(ging.UserIdentity)

	tokenName := DefaultTokenName
	tokenValue := ""

	if tokenAuth.Extend.Option.Name != "" {
		tokenName = tokenAuth.Extend.Option.Name
	}

	//从请求头获取>从查询参数获取取>最后从请求体获取>从Cookie获取
	if token, isOk := ctx.Request.Header[tokenName]; !isOk {
		if token, ok := ctx.Request.URL.Query()[tokenName]; !ok {
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
		return nil, errors.New("empty token error")
	}

	if err := userIdentity.DecryptAES([]byte(tokenAuth.Extend.EncryptKey), tokenValue); err != nil {
		return nil, err
	}

	return userIdentity, nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 错误处理
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (tokenAuth *TokenAuthentication) errorHandler(ctx *gin.Context) {
	requestUrl := ctx.Request.URL.RequestURI()
	if returnUrl == "" || returnUrl == "/" || returnUrl == "/#/" || returnUrl == "#/" {
		returnUrl = ctx.Request.URL.RequestURI()
	}

	//认证失败的处理
	if ging.IsAjax(ctx) {
		result.JsonResult(ctx, ging.NewError(199, "身份未认证")).Render()
	} else {
		authorizeUrl := fmt.Sprintf("%s?returnurl=%s", tokenAuth.Extend.AuthorizeUrl, glib.UrlEncode(requestUrl))
		result.RedirectResult(ctx, authorizeUrl).Render()
	}

	ctx.Abort()

	return
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登陆
 * userIdentity: 用户标示符域模型
 * isRemember: 是否持久化登陆信息
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (tokenAuth *TokenAuthentication) Logon(ctx *gin.Context, userIdentity *ging.UserIdentity) string {
	//当前时间戳加上秒数
	maxAge := 900
	if tokenAuth.Extend.Option.MaxAge > 0 {
		maxAge = tokenAuth.Extend.Option.MaxAge
	}
	expires := time.Now().Add(time.Duration(maxAge) * time.Second).Unix()

	userIdentity.Expires = expires
	userIdentity.IsAuthenticated = true
	ticket, err := userIdentity.EncryptAES([]byte(tokenAuth.Extend.EncryptKey))
	if err != nil {
		return ""
	}

	tokenName := DefaultTokenName
	if tokenAuth.Extend.Option.Name != "" {
		tokenName = tokenAuth.Extend.Option.Name
	}

	ctx.Header(tokenName, ticket)

	return ticket
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 登出
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (tokenAuth *TokenAuthentication) Logoff(ctx *gin.Context) {
	tokenName := DefaultTokenName
	if tokenAuth.Extend.Option.Name != "" {
		tokenName = tokenAuth.Extend.Option.Name
	}

	ctx.Header(tokenName, "")
}
