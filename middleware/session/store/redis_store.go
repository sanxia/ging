package store

import (
	"fmt"
)

import (
	"github.com/gorilla/sessions"
)

/* ================================================================================
 * Redis存储接口模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	IRedisStore interface {
		IStore
		Options(*CookieOption)
	}

	RedisOption struct {
		Host     string
		Port     int
		Password string
		Prefix   string
	}

	redisStore struct {
		*RedisStoreImpl
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Session会话Cookie存储
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewRedisStore(ip string, port int, password, prefixKey string, encryptKey []byte) IRedisStore {
	host := fmt.Sprintf("%s:%d", ip, port)
	redisStoreImpl, err := NewRedisStoreImpl(10, "tcp", host, password, encryptKey)
	if err != nil {
		panic(fmt.Sprintf("connect redis error: %v", err))
	}

	if len(prefixKey) > 0 {
		redisStoreImpl.SetKeyPrefix(prefixKey)
	}

	return &redisStore{redisStoreImpl}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置Redis存储选项
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisStore) Options(cookie *CookieOption) {
	s.RedisStoreImpl.Options = &sessions.Options{
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		MaxAge:   cookie.MaxAge,
		Secure:   cookie.IsSecure,
		HttpOnly: cookie.IsHttpOnly,
	}

	s.SetMaxAge(cookie.MaxAge)
}
