package session

import (
	"net/http"
)

import (
	"github.com/gorilla/sessions"
)

/* ================================================================================
 * Session接口模块
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
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
		Options(Options)
	}

	IStore interface {
		sessions.Store
	}

	Options struct {
		Path     string
		Domain   string
		MaxAge   int
		Secure   bool
		HttpOnly bool
	}

	session struct {
		name    string
		request *http.Request
		writer  http.ResponseWriter
		store   IStore
		session *sessions.Session
		written bool
	}
)

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
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *session) SessionId() string {
	return s.Session().ID
}

func (s *session) Values() map[interface{}]interface{} {
	return s.Session().Values
}

func (s *session) Get(key interface{}) interface{} {
	return s.Session().Values[key]
}

func (s *session) Set(key interface{}, value interface{}) {
	s.Session().Values[key] = value
	s.written = true
}

func (s *session) AddFlash(value interface{}, vars ...string) {
	s.Session().AddFlash(value, vars...)
	s.written = true
}

func (s *session) Flashes(vars ...string) []interface{} {
	s.written = true
	return s.Session().Flashes(vars...)
}

func (s *session) Save() error {
	if s.Written() {
		err := s.Session().Save(s.request, s.writer)
		if err == nil {
			s.written = false
		}
		return err
	}
	return nil
}

func (s *session) Delete(key interface{}) {
	delete(s.Session().Values, key)
	s.written = true
}

func (s *session) Clear() {
	for key := range s.Session().Values {
		s.Delete(key)
	}
}

func (s *session) Written() bool {
	return s.written
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置会话选项
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *session) Options(options Options) {
	s.Session().Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
