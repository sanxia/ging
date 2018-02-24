package ging

/* ================================================================================
 * 错误数据域结构
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type CustomError struct {
	Code int32
	Msg  string
}

func (err CustomError) Error() string {
	return err.Msg
}
