package session

import (
	"github.com/gorilla/sessions"
)

/* ================================================================================
 * Cookie存储接口模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */
type (
	ICookieStore interface {
		IStore
		Options(Options)
	}

	cookieStore struct {
		*sessions.CookieStore
	}
)

func NewCookieStore(keyPairs ...[]byte) ICookieStore {
	return &cookieStore{sessions.NewCookieStore(keyPairs...)}
}

func (c *cookieStore) Options(options Options) {
	c.CookieStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
