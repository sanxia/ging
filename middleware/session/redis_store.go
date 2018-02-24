package session

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
type IRedisStore interface {
	IStore
	Options(Options)
}

type redisStore struct {
	*RediStore
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取Redis存储
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewRedisStore(ip string, port int, password, prefixKey string, keyPairs ...[]byte) IRedisStore {
	redisIp := fmt.Sprintf("%s:%d", ip, port)
	s, err := NewRediStore(10, "tcp", redisIp, password, keyPairs...)
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
func (s *redisStore) Options(options Options) {
	s.RediStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}

	s.RediStore.SetMaxAge(options.MaxAge)
}
