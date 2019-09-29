package ging

/* ================================================================================
 * Custom Error
 * qq group: 582452342
 * email   : 2091938785@qq.com
 * author  : 美丽的地球啊 - mliu
 * ================================================================================ */
type CustomError struct {
	Code int32
	Msg  string
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * instantiate custom error
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewCustomError(msg string) *CustomError {
	return &CustomError{
		Code: 111,
		Msg:  msg,
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * instantiate custom error
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func NewError(code int32, msg string) *CustomError {
	return &CustomError{
		Code: code,
		Msg:  msg,
	}
}

/* ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 * error interface implementation
 * ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ */
func (err CustomError) Error() string {
	return err.Msg
}
