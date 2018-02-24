package common

/* ================================================================================
 * 地址信息数据域结构
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type Gps struct {
	Longitude float64 `form:"longitude" json:"longitude"`
	Latitude  float64 `form:"latitude" json:"latitude"`
	Location  string  `form:"location" json:"location"`
}

type Address struct {
	Contact  string `form:"contact" json:"contact"`
	Mobile   string `form:"mobile" json:"mobile"`
	Address  string `form:"address" json:"address"`
	PostCode string `form:"post_code" json:"post_code"`
}
