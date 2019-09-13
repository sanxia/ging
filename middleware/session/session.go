package session

import (
	"net/http"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/sanxia/ging/middleware/session/store"
)

/* ================================================================================
 * Session接口模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	ISession interface {
		SessionId() string
		Values() map[interface{}]interface{}
		Get(key interface{}) interface{}
		Set(key interface{}, val interface{})
		AddFlash(value interface{}, vars ...string)
		Flashes(vars ...string) []interface{}
		Save() error
		Delete(key interface{})
		Clear()
		Options(store.CookieOption)
	}

	session struct {
		name    string
		request *http.Request
		writer  http.ResponseWriter
		session *sessions.Session
		store   store.IStore
		written bool
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 初始化会话接口
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewSession(ctx *gin.Context) ISession {
	name := ctx.MustGet(sessionName).(string)
	store := ctx.MustGet(sessionStore).(store.IStore)

	return &session{
		name:    name,
		request: ctx.Request,
		writer:  ctx.Writer,
		store:   store,
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取Session
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *session) Session() *sessions.Session {
	if s.session == nil {
		if session, err := s.store.Get(s.request, s.name); err != nil {
			panic(err)
		} else {
			s.session = session
		}
	}

	return s.session
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取会话Id
 * redis store 存储模式下有效
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *session) SessionId() string {
	return s.Session().ID
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取会话值字典
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *session) Values() map[interface{}]interface{} {
	return s.Session().Values
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取会话值
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *session) Get(key interface{}) interface{} {
	return s.Session().Values[key]
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置会话值
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *session) Set(key interface{}, value interface{}) {
	s.Session().Values[key] = value
	s.written = true
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 添加会话值
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *session) AddFlash(value interface{}, vars ...string) {
	s.Session().AddFlash(value, vars...)
	s.written = true
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取指定的会话值，然后删除
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *session) Flashes(vars ...string) []interface{} {
	s.written = true

	return s.Session().Flashes(vars...)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 保存会话
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *session) Save() error {
	if s.written {
		err := s.Session().Save(s.request, s.writer)

		if err == nil {
			s.written = false
		}
		return err
	}

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 删除指定会话Key
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *session) Delete(key interface{}) {
	delete(s.Session().Values, key)
	s.written = true
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 清除所有会话Key
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *session) Clear() {
	for key := range s.Session().Values {
		s.Delete(key)
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置会话选项
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *session) Options(cookie store.CookieOption) {
	s.Session().Options = &sessions.Options{
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		MaxAge:   cookie.MaxAge,
		Secure:   cookie.IsSecure,
		HttpOnly: cookie.IsHttpOnly,
	}
}
