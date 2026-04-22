package errorx

import (
	"fmt"
)

// Wrap 根据现有 err 包装一个新 err
// 提供了便捷的拼接函数，不用手写模板代码
// 相当于 fmt.Errorf("%s > %w", title, err)
func Wrap(err error, desc string) error {
	return fmt.Errorf("%s > %w", desc, err)
}

// Wrapf 根据现有 err 包装一个新 err
// 提供了便捷的拼接函数，不用手写模板代码
// 相当于 fmt.Errorf("%s > %w", fmt.Sprintf(format, a...), err)
func Wrapf(err error, format string, a ...any) error {
	return fmt.Errorf("%s > %w", fmt.Sprintf(format, a...), err)
}

// Wraps 返回 Wrap 版本的错误信息字符串
func Wraps(err error, desc string) string {
	return fmt.Sprintf("%s > %v", desc, err)
}

// Wrapfs 返回 Wrapf 版本的错误信息字符串
func Wrapfs(err error, format string, a ...any) string {
	return fmt.Sprintf("%s > %v", fmt.Sprintf(format, a...), err)
}
