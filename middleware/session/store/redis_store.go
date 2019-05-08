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
		*RediStore
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Session会话Cookie存储
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewRedisStore(ip string, port int, password, prefixKey string, encryptKey []byte) IRedisStore {
	host := fmt.Sprintf("%s:%d", ip, port)
	s, err := NewRediStore(10, "tcp", host, password, encryptKey)
	if err != nil {
		panic(fmt.Sprintf("connect redis error: %v", err))
	}

	if len(prefixKey) > 0 {
		s.SetKeyPrefix(prefixKey)
	}

	return &redisStore{s}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置Redis存储选项
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *redisStore) Options(cookie *CookieOption) {
	s.RediStore.Options = &sessions.Options{
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		MaxAge:   cookie.MaxAge,
		Secure:   cookie.Secure,
		HttpOnly: cookie.HttpOnly,
	}

	s.SetMaxAge(cookie.MaxAge)
}
