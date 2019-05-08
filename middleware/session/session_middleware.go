package session

import (
	"strings"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
	"github.com/sanxia/ging/middleware/session/store"
)

/* ================================================================================
 * Session中间件模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
const (
	sessionName             = "__session_name___"
	sessionStore            = "__session_store__"
	SESSION_IDENTITY string = "__ging_s__"
)

type (
	SessionExtend struct {
		StoreType  string
		EncryptKey string
		Cookie     *store.CookieOption
		Redis      *store.RedisOption
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Session会话中间件
 * option: option扩展数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func SessionMiddleware(extend *SessionExtend) gin.HandlerFunc {
	storeType := "cookie"
	if len(extend.StoreType) > 0 {
		storeType = strings.ToLower(extend.StoreType)
	}

	if storeType == "cookie" {
		return CookieStore(extend)
	}

	return RedisStore(extend)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Cookie存储
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func CookieStore(extend *SessionExtend) gin.HandlerFunc {
	cookieStore := store.NewCookieStore([]byte(extend.EncryptKey))
	cookieStore.Options(extend.Cookie)

	return SessionStore(extend.Cookie.Name, cookieStore)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Redis存储
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func RedisStore(extend *SessionExtend) gin.HandlerFunc {
	redisStore := store.NewRedisStore(
		extend.Redis.Host,
		extend.Redis.Port,
		extend.Redis.Password,
		extend.Redis.Prefix,
		[]byte(extend.EncryptKey),
	)

	redisStore.Options(extend.Cookie)

	return SessionStore(extend.Cookie.Name, redisStore)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 会话存储
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func SessionStore(name string, storeImpl store.IStore) gin.HandlerFunc {
	if len(name) == 0 {
		name = SESSION_IDENTITY
	}

	return func(ctx *gin.Context) {
		ctx.Set(sessionName, name)
		ctx.Set(sessionStore, storeImpl)

		defer context.Clear(ctx.Request)

		ctx.Next()
	}
}
