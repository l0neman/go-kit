package debug

import (
	"log"

	"github.com/l0neman/go-kit/debug"
)

func Test() {
	for i := 0; i < 1024; i++ {
		debug.PrintThrottled(func(count int) {
			log.Println("log 1 打印次数", count)
		})
	}

	for i := 0; i < 4096; i++ {
		debug.PrintThrottled(func(count int) {
			log.Println("log 2 打印次数", count)
		})
	}
}
