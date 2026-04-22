package executor

import (
	"log"
	"time"

	"github.com/l0neman/go-kit/executor"
)

func testLoopTask1() {
	scheduler := executor.NewLoopTask(time.Second)
	go func() {
		scheduler.Start(func() {
			log.Printf("run")
		})
	}()
	go func() {
		time.Sleep(5 * time.Second)
		scheduler.Stop()
	}()
	select {}
}

func testLoopTask2() {
	task := executor.GoLoopTask(func() {
		log.Printf("run")
	}, time.Second)

	time.Sleep(5 * time.Second)
	task.Stop()
}

func Test() {
	testLoopTask2()
}
