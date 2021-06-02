package poolworker

import (
	"log"
	"sync"
	"time"
)

type Dispatcher struct {
	workers     []*Worker
	pool        chan chan *Task // A pool of workers channels that are registered
	queue       chan *Task      // tasks are sent in this channel
	concurrency int
	quit        chan bool
	stopped     bool
	wg          sync.WaitGroup
	stopChan    chan bool
}

func New(concurrency, queueSize int) *Dispatcher {
	disp := Dispatcher{
		workers:     make([]*Worker, concurrency),
		pool:        make(chan chan *Task, concurrency),
		queue:       make(chan *Task, queueSize),
		concurrency: concurrency,
		quit:        make(chan bool),
		wg:          sync.WaitGroup{},
		stopChan:    make(chan bool),
	}
	return &disp
}

func (d *Dispatcher) Start() *Dispatcher {
	log.Printf("Starting workers")
	for i := 0; i < d.concurrency; i++ {
		d.workers[i] = NewWorker(i+1, d.pool, 0)
		d.workers[i].Start(&d.wg)
	}
	go d.dispatch()
	return d
}

func (d *Dispatcher) Stop() {
	log.Printf("Stopping workers...")
	d.stopped = true
	close(d.queue)
	close(d.quit)
	<-d.stopChan
	for _, w := range d.workers {
		w.Stop()
	}
	d.wg.Wait()
}

func (d *Dispatcher) Queue(task *Task, timeout time.Duration) bool {
	if d.stopped {
		return false
	}
	d.queue <- task
	log.Printf("%s queued\n", task.ID)
	return true
}

func (d *Dispatcher) dispatch() {
	defer func() { close(d.stopChan) }()
	for {
		select {
		case task := <-d.queue:
			channel := <-d.pool
			channel <- task
		case <-d.quit:
			if len(d.queue) > 0 {
				d.processWaitingTasks()
			}
			return
		}
	}
}

func (d *Dispatcher) processWaitingTasks() {
	log.Printf("wait pending tasks to complete")
	for {
		task, done := <-d.queue
		if !done {
			return
		}
		channel := <-d.pool
		channel <- task
	}
}
