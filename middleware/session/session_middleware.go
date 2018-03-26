package session

import (
	"strings"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
)

/* ================================================================================
 * Session中间件模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
const (
	sessionName  = "sessionName"
	sessionStore = "sessionStore"
)

type SessionOption struct {
	StoreType  string
	EncryptKey string
	Cookie     *SessionCookieOption
	Redis      *SessionRedisOption
}

type SessionCookieOption struct {
	Name     string
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

type SessionRedisOption struct {
	Host     string
	Port     int
	Password string
	Prefix   string
}

func SessionMiddleware(sessionOption *SessionOption) gin.HandlerFunc {
	storeType := strings.ToLower(sessionOption.StoreType)
	if storeType == "cookie" {
		return CookieStoreSessionMiddleware(sessionOption)
	}
	return RedisStoreSessionMiddleware(sessionOption)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Cookie 存储
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func CookieStoreSessionMiddleware(sessionOption *SessionOption) gin.HandlerFunc {
	store := NewCookieStore([]byte(sessionOption.EncryptKey))

	//session id 的cookie存储属性设置
	options := Options{
		Name:     sessionOption.Cookie.Name,
		Path:     sessionOption.Cookie.Path,
		Domain:   sessionOption.Cookie.Domain,
		MaxAge:   sessionOption.Cookie.MaxAge,
		Secure:   sessionOption.Cookie.Secure,
		HttpOnly: sessionOption.Cookie.HttpOnly,
	}
	store.Options(options)

	return StoreSessionMiddleware(sessionOption.Cookie.Name, store)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Redis 存储
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func RedisStoreSessionMiddleware(sessionOption *SessionOption) gin.HandlerFunc {
	store := NewRedisStore(
		sessionOption.Redis.Host,
		sessionOption.Redis.Port,
		sessionOption.Redis.Password,
		sessionOption.Redis.Prefix,
		[]byte(sessionOption.EncryptKey),
	)

	//session id 的cookie存储属性设置
	options := Options{
		Name:     sessionOption.Cookie.Name,
		Path:     sessionOption.Cookie.Path,
		Domain:   sessionOption.Cookie.Domain,
		MaxAge:   sessionOption.Cookie.MaxAge,
		Secure:   sessionOption.Cookie.Secure,
		HttpOnly: sessionOption.Cookie.HttpOnly,
	}
	store.Options(options)

	return StoreSessionMiddleware(sessionOption.Cookie.Name, store)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 会话存储
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func StoreSessionMiddleware(name string, store IStore) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(sessionName, name)
		ctx.Set(sessionStore, store)
		defer context.Clear(ctx.Request)
		ctx.Next()
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取会话接口
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Get(ctx *gin.Context) ISession {
	name := ctx.MustGet(sessionName).(string)
	store := ctx.MustGet(sessionStore).(IStore)
	return &session{
		name:    name,
		request: ctx.Request,
		writer:  ctx.Writer,
		store:   store,
	}
}
