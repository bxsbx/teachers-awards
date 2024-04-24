package util

import (
	"strconv"
	"time"
)

const (
	YMD_HMS = "2006-01-02 15:04:05"
	YMD     = "2006-01-02"
	HMS     = "15:04:05"
)

var timeLocal, _ = time.LoadLocation("Asia/Chongqing")

func NowTime() time.Time {
	return time.Now().In(timeLocal)
}

func TimeToLocation(t *time.Time) {
	if t != nil {
		format := t.Format(YMD_HMS)
		newT, _ := time.ParseInLocation(YMD_HMS, format, timeLocal)
		*t = newT
	}
}

func NowTimeStr() string {
	return time.Now().In(timeLocal).Format(YMD_HMS)
}

func NowDate() time.Time {
	dateStr := time.Now().In(timeLocal).Format(YMD)
	parse, _ := time.Parse(YMD, dateStr)
	return parse
}

func NowDateStr() string {
	return time.Now().In(timeLocal).Format(YMD)
}

func GetYmdByTime(time time.Time) (year, month, day int) {
	year = time.Year()
	month, _ = strconv.Atoi(time.Month().String())
	day = time.Day()
	return
}

func GetTomorrowDate() *time.Time {
	now := time.Now().In(timeLocal)
	dateStr := now.AddDate(0, 0, 1).Format(YMD)
	tomorrowDate, _ := time.Parse(YMD, dateStr)
	return &tomorrowDate
}

func StrToTime(str string) (*time.Time, error) {
	parse, err := time.ParseInLocation(YMD_HMS, str, timeLocal)
	if err != nil {
		return nil, err
	}
	return &parse, nil
}
