package store

import (
	"fmt"
)

import (
	"github.com/gorilla/sessions"
)

/* ================================================================================
 * redis store interface
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	IRedisStore interface {
		IStore
		Options(CookieOption)
	}

	RedisOption struct {
		Host     string
		Port     int
		Password string
		Prefix   string
		Db       int
	}

	redisStore struct {
		*RedisStoreImpl
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Session会话Cookie存储
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewRedisStore(ip string, port int, password, prefixKey string, encryptKey []byte, dbIndex int) IRedisStore {
	host := fmt.Sprintf("%s:%d", ip, port)
	redisStoreImpl, err := NewRedisStoreImpl("tcp", host, password, encryptKey, 10, dbIndex)
	if err != nil {
		panic(fmt.Sprintf("connect redis error: %v", err))
	}

	if len(prefixKey) > 0 {
		redisStoreImpl.SetPrefix(prefixKey)
	}

	return &redisStore{redisStoreImpl}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置Redis存储选项
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisStore) Options(cookie CookieOption) {
	s.RedisStoreImpl.Options = &sessions.Options{
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		MaxAge:   cookie.MaxAge,
		Secure:   cookie.IsSecure,
		HttpOnly: cookie.IsHttpOnly,
	}

	s.SetMaxAge(cookie.MaxAge)
}
