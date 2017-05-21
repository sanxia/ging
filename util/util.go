package util

import (
	"crypto/md5"
	"io"
	"os"
	"strings"
)

/* ================================================================================
 * 帮助函数
 * author: 美丽的地球啊
 * ================================================================================ */

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * 用指定的字符串分隔源字符串为字符串切片
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func StringToStringSlice(sourceString string, args ...string) []string {
	result := make([]string, 0)

	if len(sourceString) == 0 {
		return result
	}

	splitString := ","
	if len(args) == 1 {
		splitString = args[0]
	}

	stringSlice := strings.Split(sourceString, splitString)
	for _, v := range stringSlice {
		if v != "" {
			result = append(result, v)
		}
	}

	return result
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
