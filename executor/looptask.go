package executor

import (
	"context"
	"sync"
	"time"
)

// LoopTask 管理定时任务的执行
type LoopTask struct {
	ticker *time.Ticker
	ctx    context.Context
	cancel context.CancelFunc
	lock   sync.Mutex
	stop   bool
}

// NewLoopTask 创建新的循环任务
func NewLoopTask(interval time.Duration) *LoopTask {
	ctx, cancel := context.WithCancel(context.Background())
	return &LoopTask{
		stop:   false,
		ticker: time.NewTicker(interval),
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start 启动循环任务
func (t *LoopTask) Start(task func()) {
	t.lock.Lock()
	if t.stop {
		t.lock.Unlock()
		return
	}

	t.lock.Unlock()

	task()

	for {
		select {
		case <-t.ticker.C:
			task()
		case <-t.ctx.Done():
			return
		}
	}
}

// Stop 停止任务，停止后将不能再次启动
func (t *LoopTask) Stop() {
	t.lock.Lock()
	defer t.lock.Unlock()

	if t.stop {
		return
	}

	t.stop = true
	t.cancel()
	t.ticker.Stop()
	return
}

func (t *LoopTask) Close() error {
	t.Stop()
	return nil
}

// Reset 重置任务执行间隔
func (t *LoopTask) Reset(interval time.Duration) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.ticker.Reset(interval)
}

func GoLoopTask(task func(), interval time.Duration) *LoopTask {
	loopTask := NewLoopTask(interval)
	go func() {
		loopTask.Start(task)
	}()
	return loopTask
}
