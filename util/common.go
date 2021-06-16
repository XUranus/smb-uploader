package util

func StringOmit(str string, length int) string {
	if len(str) > length {
		return str[0: length] + "..."
	} else {
		return str
	}
}

func MakeSpace(length int) string {
	ret := ""
	for i:= 0; i < length; i++ {
		ret += " "
	}
	return ret
}
