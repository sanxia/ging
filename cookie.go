package ging

/* ================================================================================
 * Cookie
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type (
	Cookie struct {
		Name       string `json:"name"`
		Path       string `json:"path"`
		Domain     string `json:"domain"`
		MaxAge     int    `json:"max_age"`
		IsHttpOnly bool   `json:"is_http_only"`
		IsSecure   bool   `json:"is_secure"`
	}
)
