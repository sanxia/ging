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
 * author  : 美丽的地球啊
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
func NewRedisStore(ip string, port int, password, keyPrefix string, keyPairs ...[]byte) IRedisStore {
	redisIp := fmt.Sprintf("%s:%d", ip, port)
	s, err := NewRediStore(10, "tcp", redisIp, password, keyPairs...)
	if err != nil {
		panic(fmt.Sprintf("connect redis error: %v", err))
	}

	if len(keyPrefix) > 0 {
		s.SetKeyPrefix(keyPrefix)
	}

	return &redisStore{s}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置Redis存储选项
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (c *redisStore) Options(options Options) {
	c.RediStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
