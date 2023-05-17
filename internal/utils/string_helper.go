package utils

func StartWith(str, charsets string) bool {
	return str[0:len(charsets)] == charsets
}
