package pongo

import (
	"crypto/md5"
	"io"
	"os"
)

import (
	"github.com/flosch/pongo2"
)

/* ================================================================================
 * Pongo模版引擎帮助模块
 * author: 美丽的地球啊
 * ================================================================================ */
func Render(templateString string, data interface{}) (string, error) {
	tpl, err := pongo2.FromString(templateString)
	if err != nil {
		panic(err)
	}

	result, err := tpl.Execute(data.(map[string]interface{}))
	if err != nil {
		panic(err)
	}

	return result, err
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * Md5哈希
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func Md5(data string) string {
	m := md5.New()
	io.WriteString(m, data)
	return hex.EncodeToString(m.Sum(nil))
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 判断文件是否存在
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func FileIsExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
