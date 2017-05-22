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
 * qq: 2091938785
 * email: 2091938785@qq.com
 * author: 美丽的地球啊
 * ================================================================================ */
const (
	UserIdentityKey string = "__ging_u_id__"
)

type UserIdentity struct {
	UserId   uint64
	Nickname string
	Role     string
}

func (c *UserIdentity) Serialize() ([]byte, error) {
	var b bytes.Buffer
	e := gob.NewEncoder(&b)
	if err := e.Encode(c); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (c *UserIdentity) Deserialize(by []byte) error {
	var b bytes.Buffer
	b.Write(by)
	d := gob.NewDecoder(&b)
	err := d.Decode(&c)
	return err
}

func (c *UserIdentity) EncryptAES(key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	text, err := c.Serialize()
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

func (c *UserIdentity) DecryptAES(key []byte, input string) error {
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

	return c.Deserialize(data)
}
