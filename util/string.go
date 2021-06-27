package util

import (
	"fmt"
	"strconv"
)

func StringOmit(str string, length int) string {
	if len(str) > length {
		return str[0: length] + "..."
	} else {
		return str
	}
}

/**
	return a string made up with space
 */
func MakeSpace(length int) string {
	ret := ""
	for i:= 0; i < length; i++ {
		ret += " "
	}
	return ret
}

func NumberWithComma(n int64) string {
	if n >= 1000 {
		return fmt.Sprintf("%v,%v",NumberWithComma(n / 1000) , n % 1000)
	} else {
		return strconv.FormatInt(n, 10)
	}
}