package service

import (
	"context"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

type TaskKey struct {
	Source  string
	ComicId string
}

type Task struct {
	key  TaskKey
	run  func() error
	done chan struct{}
	err  error
}

type Queue struct {
	mu      sync.Mutex
	cond    *sync.Cond
	tasks   map[TaskKey]*Task
	pending []*Task
}

// NewQueue 创建队列并启动 worker goroutine。
func NewQueue() *Queue {
	q := &Queue{
		tasks: make(map[TaskKey]*Task),
	}
	q.cond = sync.NewCond(&q.mu)
	go q.worker()
	return q
}

func (q *Queue) Submit(ctx context.Context, key TaskKey, run func() error) error {
	q.mu.Lock()
	t, ok := q.tasks[key]
	if !ok {
		t = &Task{
			key:  key,
			run:  run,
			done: make(chan struct{}),
		}
		q.tasks[key] = t
		q.pending = append(q.pending, t)
		log.Infof("Task enqueued: %s/%s", key.Source, key.ComicId)
		q.cond.Signal()
	}
	q.mu.Unlock()

	select {
	case <-t.done:
		return t.err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *Queue) worker() {
	for {
		q.mu.Lock()
		for len(q.pending) == 0 {
			q.cond.Wait()
		}
		// 取出任务
		t := q.pending[0]
		q.pending = q.pending[1:]
		q.mu.Unlock()

		log.Infof("Task started: %s/%s", t.key.Source, t.key.ComicId)
		start := time.Now()
		t.err = t.run()
		close(t.done)

		q.mu.Lock()
		delete(q.tasks, t.key)
		q.mu.Unlock()

		if t.err != nil {
			log.Errorf("Task failed: %s/%s (%v)", t.key.Source, t.key.ComicId, t.err)
		} else {
			log.Infof("Task done: %s/%s (%s)", t.key.Source, t.key.ComicId, time.Since(start).Round(time.Millisecond))
		}
	}
}
