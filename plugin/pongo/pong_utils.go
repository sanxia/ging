package pongo

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
