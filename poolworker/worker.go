package poolworker

import (
	"fmt"
	"log"
	"sync"
	"time"
	"workerpool/random"
)

type Task struct {
	ID     string
	Action func(arg ...interface{}) (interface{}, error)
	Args   interface{}
}

func NewTask(name string, action func(...interface{}) (interface{}, error), args interface{}) *Task {
	return &Task{
		ID:     fmt.Sprintf("task-%s-%s", random.RandString(10), name),
		Action: action,
		Args:   args,
	}
}

type Worker struct {
	ID      int
	pool    chan chan *Task
	channel chan *Task
	timeout time.Duration
}

func NewWorker(id int, pool chan chan *Task, timeout time.Duration) *Worker {
	return &Worker{
		ID:      id,
		pool:    pool,
		channel: make(chan *Task),
		timeout: timeout,
	}
}

func (w *Worker) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case w.pool <- w.channel:
			case task, ok := <-w.channel:
				if !ok {
					log.Printf("worker-%d exited\n", w.ID)
					return
				}
				if task == nil {
					continue
				}
				log.Printf("worker-%d: %s received, args=(%d)\n", w.ID, task.ID, task.Args)
				_, err := task.Action(task.Args)
				if err != nil {
					log.Printf("worker-%d: %s failed, error='%v'\n", w.ID, task.ID, err)
				} else {
					log.Printf("worker-%d: %s succeed\n", w.ID, task.ID)
				}
			}
		}
	}()
	log.Printf("worker-%d started\n", w.ID)
}

func (w Worker) Stop() {
	close(w.channel)
}
