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

	cookieStore struct {
		*sessions.CookieStore
	}

	CookieOption struct {
		Name       string `json:"name"`
		Path       string `json:"path"`
		Domain     string `json:"domain"`
		MaxAge     int    `json:"max_age"`
		IsHttpOnly bool   `json:"is_http_only"`
		IsSecure   bool   `json:"is_secure"`
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
		Secure:   cookie.IsSecure,
		HttpOnly: cookie.IsHttpOnly,
	}
}
