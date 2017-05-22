package util

import (
	"strings"
)

/* ================================================================================
 * 帮助函数
 * author: 美丽的地球啊
 * ================================================================================ */

func IsInRole(userRole, allowRole string) bool {
	isInRole := false

	if len(userRole) == 0 || len(allowRole) == 0 {
		return false
	}

	roles := StringToStringSlice(allowRole)
	currentRoles := StringToStringSlice(userRole)
	for _, role := range roles {
		for _, currentRole := range currentRoles {
			if currentRole == role {
				isInRole = true
				break
			}
		}
	}
	return isInRole
}

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
