package ging

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

import (
	"github.com/sanxia/glib"
)

/* ================================================================================
 * Token标识数据域结构
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
const (
	USER_IDENTITY string = "__ging_u__"
)

type (
	IToken interface {
		GetToken() string
		ParseToken(token string) error

		GetPayload() *TokenPayload
		SetPayload(payload *TokenPayload)

		SetExpires(expires int64)
		SetAuthenticated(isAuthenticated bool)

		IsAuthenticated() bool
		IsExpired() bool
		IsValid() bool
	}
)

type (
	token struct {
		payload   *TokenPayload //载荷
		signature string        //签名指纹
		secret    string        //秘匙
	}

	TokenPayload struct {
		UserId          string   `json:"sub"`              //用户id
		Username        string   `json:"username"`         //用户名
		Nickname        string   `json:"nickname"`         //用户昵称
		Avatar          string   `json:"avatar"`           //用户图像
		Roles           []string `json:"roles"`            //角色名集合
		Start           int64    `json:"iat"`              //签发时间（距离1970-1-1的秒数）
		Expires         int64    `json:"exp"`              //过期时间（距离1970-1-1的秒数）
		IsAuthenticated bool     `json:"is_authenticated"` //是否已验证
	}

	Cookie struct {
		Name     string
		Path     string
		Domain   string
		MaxAge   int
		HttpOnly bool
		Secure   bool
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 初始化Token标识
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewToken(secret string) *token {
	return &token{
		payload: &TokenPayload{},
		secret:  secret,
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取Token字符串
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *token) GetToken() string {
	tokenTicket := ""
	s.signature = s.tokenSignature()

	if payload, err := s.payload.Serialize(); err == nil {
		tokenTicket = fmt.Sprintf("%s.%s", glib.ToBase64(payload, true), glib.ToBase64(s.signature, true))
	}

	return tokenTicket
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Parse Token
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *token) ParseToken(token string) error {
	if len(token) == 0 {
		return errors.New("token format error")
	}

	tokens := glib.StringToStringSlice(token, ".")
	if len(tokens) < 2 || len(tokens) > 3 {
		return errors.New("token format error")
	}

	tokenPayload := ""
	tokenSignature := ""
	if len(tokens) == 2 {
		tokenPayload = tokens[0]
		tokenSignature = tokens[1]
	} else {
		tokenPayload = tokens[1]
		tokenSignature = tokens[2]
	}

	if payload, err := glib.FromBase64(tokenPayload, true); err != nil {
		return err
	} else {
		if err := s.payload.Deserialize(payload); err != nil {
			return err
		}
	}

	if signature, err := glib.FromBase64(tokenSignature, true); err != nil {
		return err
	} else {
		s.signature = signature
	}

	//是否有效签名
	if isValid := s.IsValid(); !isValid {
		return errors.New("token signature error")
	}

	//是否过期
	if isExpired := s.IsExpired(); isExpired {
		return errors.New("token expired error")
	}

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取Payload
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *token) GetPayload() *TokenPayload {
	return s.payload
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置Payload
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *token) SetPayload(payload *TokenPayload) {
	s.payload = payload
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置Expires
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *token) SetExpires(expires int64) {
	s.payload.Start = time.Now().Unix()

	if expires <= 0 {
		s.payload.Expires = s.payload.Start
	} else {
		s.payload.Expires = expires
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 设置是否已认证
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *token) SetAuthenticated(isAuthenticated bool) {
	s.payload.IsAuthenticated = isAuthenticated
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 获取Token签名
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *token) tokenSignature() string {
	signature := ""
	if len(s.secret) > 0 {
		if payload, err := s.payload.Serialize(); err == nil {
			signature = glib.HmacSha256(glib.ToBase64(payload, true), s.secret)
		}
	}

	return signature
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Token是否已认证
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *token) IsAuthenticated() bool {
	isAuthenticated := s.payload.IsAuthenticated
	return isAuthenticated
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Token签名是否有效
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *token) IsValid() bool {
	isValid := false
	if signature := s.tokenSignature(); strings.ToLower(s.signature) == strings.ToLower(signature) {
		isValid = true
	}

	return isValid
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Token是否已过期
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *token) IsExpired() bool {
	isExpired := false

	if s.payload.Start <= 0 || s.payload.Expires <= 0 {
		isExpired = true
	}

	//签发日期是否大于等于失效日期
	if s.payload.Start >= s.payload.Expires {
		isExpired = true
	}

	//当前日期是否大于失效日期
	if time.Now().After(time.Unix(s.payload.Expires, 0)) {
		isExpired = true
	}

	return isExpired
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * UserId string to UserId int64
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *TokenPayload) UserIdInt64() int64 {
	return glib.StringToInt64(s.UserId)
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * TokenPayload json序列化
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *TokenPayload) Serialize() (string, error) {
	jsonString, err := glib.ToJson(s)
	if err != nil {
		return "", err
	}

	return jsonString, nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * TokenPayload json反序列化
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *TokenPayload) Deserialize(payload string) error {
	return glib.FromJson(payload, &s)
}
