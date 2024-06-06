package errorz

import (
	"gorm.io/gorm"
)

func GlobalError(err error) (code int, msg string) {
	var oneErr *MyError
	myErr, ok := err.(*MyError)
	temp := myErr
	if ok {
		oneErr = myErr
	}
	for ok {
		myErr = temp
		if myErr.err != nil {
			err = myErr.err
			temp, ok = err.(*MyError)
		} else {
			break
		}
	}
	switch err {
	case gorm.ErrRecordNotFound:
		code, msg = RECORD_NOT_FOUND, GetMsgWithCode(RECORD_NOT_FOUND)
	case myErr:
		code, msg = myErr.code, myErr.msg
	default:
		if oneErr != nil {
			code, msg = oneErr.code, oneErr.msg
		} else {
			code, msg = RESP_UNKNOWN_ERR, err.Error()
		}
	}
	return
}
