package stringx

// IsEmpty 判断字符串是否为空
func IsEmpty(str string) bool {
	return len(str) == 0
}

// HasEmpty 是否存在空字符串
func HasEmpty(strings ...string) bool {
	for _, str := range strings {
		if IsEmpty(str) {
			return true
		}
	}

	return false
}
