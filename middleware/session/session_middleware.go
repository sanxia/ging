package session

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
	_sessionName             = "__session_name___"
	_sessionStoreName        = "__session_store__"
	SESSION_IDENTITY  string = "__ging_s__"
)

type (
	SessionExtend struct {
		EncryptSecret string
		Cookie        store.CookieOption
		Redis         store.RedisOption
		IsRedis       bool
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Session会话中间件
 * option: option扩展数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func SessionMiddleware(extend SessionExtend) gin.HandlerFunc {
	if extend.IsRedis {
		return sessionRedisStore(extend)
	}

	return sessionCookieStore(extend)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * cookie store
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func sessionCookieStore(extend SessionExtend) gin.HandlerFunc {
	cookieStore := store.NewCookieStore([]byte(extend.EncryptSecret))
	cookieStore.Options(extend.Cookie)

	return sessionStore(extend.Cookie.Name, cookieStore)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * redis store
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func sessionRedisStore(extend SessionExtend) gin.HandlerFunc {
	redisStore := store.NewRedisStore(
		extend.Redis.Host,
		extend.Redis.Port,
		extend.Redis.Password,
		extend.Redis.Prefix,
		[]byte(extend.EncryptSecret),
		extend.Redis.Db,
	)

	redisStore.Options(extend.Cookie)

	return sessionStore(extend.Cookie.Name, redisStore)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * store
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func sessionStore(name string, storeImpl store.IStore) gin.HandlerFunc {
	if len(name) == 0 {
		name = SESSION_IDENTITY
	}

	return func(ctx *gin.Context) {
		ctx.Set(_sessionName, name)
		ctx.Set(_sessionStoreName, storeImpl)

		defer context.Clear(ctx.Request)

		ctx.Next()
	}
}
