package errors

// Error 异常类型
type Error struct {
	Code int64  // 错误码
	Msg  string // 错误信息
}

func (c *Error) Error() string {
	return c.Msg
}

// New 创那家异常类型
func New(code int64, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

var (
	// OK 正常
	OK = New(0, "OK")

	// ErrSystem 系统错误
	ErrSystem = New(-9999, "系统错误")
)
