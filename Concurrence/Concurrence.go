package Concurrence

import (
	"log"
)

// --------------------------- Job ---------------------
type Job interface {
	Do()
}

// --------------------------- Worker ---------------------
type Worker struct {
	JobQueue chan Job
}

func NewWorker() Worker {
	return Worker{JobQueue: make(chan Job)}
}
func (w Worker) Run(wq chan chan Job) {
	go func() {
		for {
			wq <- w.JobQueue
			select {
			case job := <-w.JobQueue:
				job.Do()
			}
		}
	}()
}

// --------------------------- WorkerPool ---------------------
type WorkerPool struct {
	workerlen   int
	JobQueue    chan Job
	WorkerQueue chan chan Job
}

func NewWorkerPool(workerlen int) *WorkerPool {
	return &WorkerPool{
		workerlen:   workerlen,
		JobQueue:    make(chan Job),
		WorkerQueue: make(chan chan Job, workerlen),
	}
}
func (wp *WorkerPool) Run() {
	log.Println("初始化worker")
	//初始化worker
	for i := 0; i < wp.workerlen; i++ {
		worker := NewWorker()
		worker.Run(wp.WorkerQueue)
	}
	// 循环获取可用的worker,往worker中写job
	go func() {
		for {
			select {
			case job := <- wp.JobQueue:
				worker := <- wp.WorkerQueue
				worker <- job
			}
		}
	}()
}

