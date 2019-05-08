package serializer

import (
	"errors"
	"fmt"
	"strings"
)

import (
	"github.com/sanxia/glib"
)

/* ================================================================================
 * 会话自定义编解码
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 自定义编码
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func CustomEncode(name, sourceData string, encryptKey string) (string, error) {
	if len(sourceData) == 0 || len(encryptKey) == 0 {
		return "", errors.New("encode err")
	}

	hashKey := encryptKey[:16]

	data := glib.ToBase64(sourceData)
	hash := glib.HmacSha256(fmt.Sprintf("%s|%s", name, data), hashKey)

	return glib.ToBase64(fmt.Sprintf("%s|%s", data, hash)), nil
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 自定义解码
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func CustomDecode(name, encodeData string, encryptKey string) (string, error) {
	decodeData, err := glib.FromBase64(encodeData)
	if err != nil {
		return "", err
	}

	datas := strings.Split(decodeData, "|")
	if len(datas) != 2 {
		return "", errors.New("decode err")
	}

	hashKey := encryptKey[:16]
	hash := glib.HmacSha256(fmt.Sprintf("%s|%s", name, datas[0]), hashKey)
	if hash != datas[1] {
		return "", errors.New("decode err")
	}

	return glib.FromBase64(datas[0])
}
