package util

import "fmt"

const (
	Minute = 60
	Hour   = 60 * Minute
	Day    = 24 * Hour
)

func SecondToTime(second int64) string {
	if second < Minute {
		return fmt.Sprintf("%v秒", second)
	} else if second < Hour {
		return fmt.Sprintf("%v分钟", second /Minute)
	} else if second < Day {
		return fmt.Sprintf("%v小时", second / Hour)
	} else {
		//return fmt.Sprintf("%v天", second / Day)
		return "> 1天"
	}
}