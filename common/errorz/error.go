package errorz

import (
	"fmt"
	"runtime"
)

const (
	STOP_METHOD = "github.com/gin-gonic/gin.(*Context).Next"
)

type MyError struct {
	code  int
	msg   string
	err   error
	stack []recorde // 记录错误链
}

func (e *MyError) Error() string {
	return fmt.Sprintf("code: %d, msg：%s", e.code, e.msg)
}

func (e *MyError) Unwrap() error {
	return e.err
}

func (e MyError) GetCode() int {
	return e.code
}

func (e MyError) GetMsg() string {
	return e.msg
}

func Code(code int) error {
	msg := GetMsgWithCode(code)
	return &MyError{
		code:  code,
		msg:   msg,
		stack: caller(),
	}
}

func CodeMsg(code int, msg string) error {
	return &MyError{
		code:  code,
		msg:   msg,
		stack: caller(),
	}
}

func CodeError(code int, err error) error {
	msg := GetMsgWithCode(code)
	return &MyError{
		code:  code,
		msg:   msg,
		err:   err,
		stack: caller(),
	}
}

func CodeMsgError(code int, msg string, err error) error {
	return &MyError{
		code:  code,
		msg:   msg,
		err:   err,
		stack: caller(),
	}
}

func Error(err error) error {
	return &MyError{
		err:   err,
		stack: caller(),
	}
}

// 获取错误码相应的错误信息
func GetMsgWithCode(code int) string {
	if v, ok := respMsg[code]; ok {
		return v
	}
	return ""
}

// 获取错误码相应的请求码
func GetHttpCodeWithCode(code int) int {
	if v, ok := httpCodeMap[code]; ok {
		return v
	}
	return 200
}

type recorde struct {
	File   string
	Line   int
	Method string
}

// 记录错误链，方便定位错误
func caller() []recorde {
	var list []recorde
	var hasMap = make(map[string]bool)
	for i := 2; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		method := runtime.FuncForPC(pc).Name()
		if !ok || method == STOP_METHOD {
			break
		}
		// 递归的避免重复打印
		key := fmt.Sprintf("%s:%d", file, line)
		if _, ok := hasMap[key]; ok {
			continue
		}

		list = append(list, recorde{
			file,
			line,
			method,
		})
	}
	var errStack []recorde
	for i := len(list) - 1; i >= 0; i-- {
		errStack = append(errStack, list[i])
	}

	return errStack
}

type ErrorCaller struct {
	Caller string `json:"caller"`
	Code   int    `json:"code,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Error  string `json:"error,omitempty"`
}

// 打印错误栈
func GetErrorCallerList(err error) []ErrorCaller {
	var errStack []ErrorCaller
	length := 0
	for {
		if err == nil {
			break
		}
		if myErr, ok := err.(*MyError); ok {
			for i := length; i < len(myErr.stack); i++ {
				v := myErr.stack[i]
				called := fmt.Sprintf("%s:%d ———— %s", v.File, v.Line, v.Method)
				errStack = append(errStack, ErrorCaller{
					Caller: called,
				})
				length++
			}
			if length > 0 {
				errStack[length-1].Msg = myErr.Error()
			}
			err = myErr.err
		} else {
			break
		}
	}
	if err != nil && length > 0 {
		errStack[length-1].Error = err.Error()
	}
	return errStack
}
