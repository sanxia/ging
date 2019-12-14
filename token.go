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
 * Token identity
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
const (
	TOKEN_IDENTITY string = "__ging_t__"
)

type (
	IToken interface {
		GetToken() string
		ParseToken(token, userAgent string) error

		GetPayload() *TokenPayload
		SetPayload(payload *TokenPayload)

		SetExpire(expire int64)
		SetAuthenticated(isAuthenticated bool)

		IsAuthenticated() bool
		IsExpired() bool
		IsValid() bool
	}
)

type (
	Token struct {
		payload   *TokenPayload //载荷
		signature string        //签名
		secret    string        //秘匙
	}

	TokenPayload struct {
		Owner           string                 `json:"iss"`              //签发者
		Domain          string                 `json:"aud"`              //接收域
		UserId          string                 `json:"sub"`              //用户id
		Username        string                 `json:"username"`         //用户名
		Nickname        string                 `json:"nickname"`         //用户昵称
		Avatar          string                 `json:"avatar"`           //用户图像
		Roles           []string               `json:"roles"`            //角色名集合
		UserAgent       string                 `json:"ua"`               //客户端代理数据
		Extend          map[string]interface{} `json:"extend"`           //扩展数据
		IsAuthenticated bool                   `json:"is_authenticated"` //是否已验证
		Start           int64                  `json:"iat"`              //签发时间（距离1970-1-1的秒数）
		Expire          int64                  `json:"exp"`              //过期时间（距离1970-1-1的秒数）
	}
)

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * instantiating token
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewToken(secret string) *Token {
	return &Token{
		payload: &TokenPayload{},
		secret:  secret,
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get token string
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Token) GetToken() string {
	tokenTicket := ""
	s.signature = s.tokenSignature()

	if payload, err := s.payload.Serialize(); err == nil {
		tokenTicket = fmt.Sprintf("%s.%s", glib.ToBase64(payload, true), glib.ToBase64(s.signature, true))
	}

	return tokenTicket
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * parse token
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Token) ParseToken(token, userAgent string) error {
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

	//ua是否匹配
	if isUserAgent := s.payload.UserAgent == userAgent; !isUserAgent {
		return errors.New("token ua error")
	}

	return nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get payload
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Token) GetPayload() *TokenPayload {
	return s.payload
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * set payload
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Token) SetPayload(payload *TokenPayload) {
	s.payload = payload
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * set expire
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Token) SetExpire(expire int64) {
	s.payload.Start = time.Now().Unix()

	if expire <= 0 {
		s.payload.Expire = s.payload.Start
	} else {
		s.payload.Expire = expire
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * set is certified
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Token) SetAuthenticated(isAuthenticated bool) {
	s.payload.IsAuthenticated = isAuthenticated
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * get token signature
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Token) tokenSignature() string {
	signature := ""
	if len(s.secret) > 0 {
		if payload, err := s.payload.Serialize(); err == nil {
			signature = glib.HmacSha256(glib.ToBase64(payload, true), s.secret)
		}
	}

	return signature
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * is certified
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Token) IsAuthenticated() bool {
	isAuthenticated := s.payload.IsAuthenticated
	return isAuthenticated
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * token signature is valid
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Token) IsValid() bool {
	isValid := false

	if signature := s.tokenSignature(); strings.ToLower(s.signature) == strings.ToLower(signature) {
		isValid = true
	}

	return isValid
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * has token expired
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *Token) IsExpired() bool {
	isExpired := false

	if s.payload.Start <= 0 || s.payload.Expire <= 0 {
		isExpired = true
	}

	//签发日期是否大于等于失效日期
	if s.payload.Start >= s.payload.Expire {
		isExpired = true
	}

	//当前日期是否大于失效日期
	if time.Now().After(time.Unix(s.payload.Expire, 0)) {
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
 * TokenPayload json serialization
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *TokenPayload) Serialize() (string, error) {
	jsonString, err := glib.ToJson(s)
	if err != nil {
		return "", err
	}

	return jsonString, nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * TokenPayload json deserialization
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *TokenPayload) Deserialize(payload string) error {
	return glib.FromJson(payload, &s)
}
