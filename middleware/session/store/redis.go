package store

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/sessions"
	"github.com/sanxia/ging/middleware/session/serializer"
	"github.com/sanxia/glib"
)

/* ================================================================================
 * Redis存储
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	RediStore struct {
		Pool          *redis.Pool
		Options       *sessions.Options
		DefaultMaxAge int
		serializer    serializer.ISessionSerializer
		maxLength     int
		keyPrefix     string
		encryptKey    string
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 实例化RediStore
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewRediStore(size int, network, address, password string, encryptKey []byte) (*RediStore, error) {
	return NewRediStoreWithPool(&redis.Pool{
		MaxIdle:     size,
		IdleTimeout: 240 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			return dial(network, address, password)
		},
	}, encryptKey)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 实例化RedisPool
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewRediStoreWithPool(pool *redis.Pool, encryptKey []byte) (*RediStore, error) {
	rs := &RediStore{
		Pool: pool,
		Options: &sessions.Options{
			Path: "/",
		},
		DefaultMaxAge: 600,
		encryptKey:    string(encryptKey),
		maxLength:     1024000,
		serializer:    serializer.GobSerializer{},
	}

	_, err := rs.ping()

	return rs, err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 连接redis
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func dial(network, address, password string) (redis.Conn, error) {
	c, err := redis.Dial(network, address)
	if err != nil {
		return nil, err
	}

	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return nil, err
		}
	}

	/*
		if _, err := c.Do("SELECT", "0"); err != nil {
			c.Close()
			return nil, err
		}
	*/

	return c, err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 实例化会话
 * IStore接口实现
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RediStore) New(request *http.Request, name string) (*sessions.Session, error) {
	var err error

	session := sessions.NewSession(s, name)
	options := *s.Options
	session.Options = &options
	session.IsNew = true

	if c, errCookie := request.Cookie(name); errCookie == nil {
		decodeData, err := serializer.CustomDecode(name, c.Value, s.encryptKey)
		if err == nil {
			session.ID = decodeData
		}

		if err == nil {
			ok, err := s.load(session)
			session.IsNew = !(err == nil && ok)
		}
	}

	return session, err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置Redis Key 前缀
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RediStore) SetKeyPrefix(keyPrefix string) {
	s.keyPrefix = keyPrefix
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置会话最大有效时间（单位：秒）
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RediStore) SetMaxAge(maxAge int) {
	s.Options.MaxAge = maxAge
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置最大长度
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RediStore) SetMaxLength(maxLength int) {
	if maxLength >= 0 {
		s.maxLength = maxLength
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 关闭redis链接
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RediStore) Close() error {
	return s.Pool.Close()
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取会话
 * IStore接口实现
 * 第一次会回调当前store的New，实例化Store
 * 随后会从Register缓存sessions map[string]sessionInfo的里取出sessionInfo对象里的s值即sessions.Session对象
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RediStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 保存会话
 * IStore接口实现
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RediStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if session.Options.MaxAge < 0 {
		if err := s.delete(session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
	} else {
		if session.ID == "" {
			session.ID = glib.ToBase64(fmt.Sprintf("%d%s", time.Now().UnixNano(), glib.RandString(16)))
		}

		if err := s.save(session); err != nil {
			return err
		}

		encodeData, err := serializer.CustomEncode(session.Name(), session.ID, s.encryptKey)
		if err != nil {
			return err
		}

		http.SetCookie(w, sessions.NewCookie(session.Name(), encodeData, session.Options))
	}

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 检测服务器是否活着
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RediStore) ping() (bool, error) {
	conn := s.Pool.Get()
	defer conn.Close()
	data, err := conn.Do("PING")
	if err != nil || data == nil {
		return false, err
	}
	return (data == "PONG"), nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 从Redis读取会话
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RediStore) load(session *sessions.Session) (bool, error) {
	conn := s.Pool.Get()
	defer conn.Close()

	if err := conn.Err(); err != nil {
		return false, err
	}

	data, err := conn.Do("GET", fmt.Sprintf("%s%s", s.keyPrefix, session.ID))
	if err != nil {
		return false, err
	}

	if data == nil {
		return false, nil
	}

	sessionValue, err := redis.Bytes(data, err)
	if err != nil {
		return false, err
	}

	return true, s.serializer.Deserialize(sessionValue, session)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 会话存入Redis
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RediStore) save(session *sessions.Session) error {
	sessionValue, err := s.serializer.Serialize(session)
	if err != nil {
		return err
	}

	if s.maxLength != 0 && len(sessionValue) > s.maxLength {
		return errors.New("SessionStore: the value to store is too big")
	}

	conn := s.Pool.Get()
	defer conn.Close()

	if err = conn.Err(); err != nil {
		return err
	}

	age := session.Options.MaxAge
	if age == 0 {
		age = s.DefaultMaxAge
	}

	_, err = conn.Do("SETEX", fmt.Sprintf("%s%s", s.keyPrefix, session.ID), age, sessionValue)

	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 从Redis删除会话
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *RediStore) delete(session *sessions.Session) error {
	conn := s.Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", fmt.Sprintf("%s%s", s.keyPrefix, session.ID)); err != nil {
		return err
	}

	return nil
}
