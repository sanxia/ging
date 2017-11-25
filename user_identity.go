package ging

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"io"
)

/* ================================================================================
 * 用户身份标识数据结构
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊
 * ================================================================================ */
const (
	UserIdentityKey string = "__ging_u_id__"
)

type UserIdentity struct {
	UserId          uint64 //用户id
	Username        string //用户名
	Nickname        string //用户昵称
	Avatar          string //用户图像
	Role            string //角色名（多个之间用逗号分隔）
	Expires         int64  //过期时间（距离1970-1-1的秒数）
	IsAuthenticated bool   //是否已验证
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 输出字典数据
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *UserIdentity) Dump() map[string]interface{} {
	data := make(map[string]interface{}, 0)
	data["id"] = s.UserId
	data["username"] = s.Username
	data["nickname"] = s.Nickname
	data["avatar"] = s.Avatar

	return data
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * gob序列化
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *UserIdentity) Serialize() ([]byte, error) {
	var b bytes.Buffer
	e := gob.NewEncoder(&b)
	if err := e.Encode(s); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * gob反序列化
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *UserIdentity) Deserialize(by []byte) error {
	var b bytes.Buffer
	b.Write(by)
	d := gob.NewDecoder(&b)
	err := d.Decode(&s)
	return err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * aes加密
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *UserIdentity) EncryptAES(key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	text, err := s.Serialize()
	if err != nil {
		return "", err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * aes解密
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (s *UserIdentity) DecryptAES(key []byte, input string) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	text, err := base64.URLEncoding.DecodeString(input)
	if err != nil {
		return err
	}
	if len(text) < aes.BlockSize {
		return errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return err
	}

	return s.Deserialize(data)
}
