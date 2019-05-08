package store

import (
	"github.com/gorilla/sessions"
)

/* ================================================================================
 * Cookie存储接口模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	IStore interface {
		sessions.Store
	}

	ICookieStore interface {
		IStore
		Options(*CookieOption)
	}

	CookieOption struct {
		Name     string
		Path     string
		Domain   string
		MaxAge   int
		HttpOnly bool
		Secure   bool
	}

	cookieStore struct {
		*sessions.CookieStore
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Session会话Cookie存储
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewCookieStore(keyPairs ...[]byte) ICookieStore {
	return &cookieStore{sessions.NewCookieStore(keyPairs...)}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置Cookie存储选项
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *cookieStore) Options(cookie *CookieOption) {
	s.CookieStore.Options = &sessions.Options{
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		MaxAge:   cookie.MaxAge,
		Secure:   cookie.Secure,
		HttpOnly: cookie.HttpOnly,
	}
}
