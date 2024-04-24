package util

import (
	"crypto/md5"
	"fmt"
	"time"
)

// md5
func Md5(str string) string {
	md5Str := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return md5Str
}

func RepeatExecWhenError[T any](n, sleepTime int, f func() (T, error)) (t T, err error) {
	for i := 0; i < n; i++ {
		t, err = f()
		if err == nil {
			return
		}
		if i != n-1 {
			time.Sleep(time.Duration(sleepTime * i * int(time.Millisecond)))
		}
	}
	return
}
