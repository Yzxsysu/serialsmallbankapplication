package main

import (
	"fmt"
	"sync"
)

type worker struct {
	ID   int           // 协程 ID
	Job  chan func()   // 存储任务的 channel
	exit chan struct{} // 协程退出信号
}

func NewWorker(id int) *worker {
	return &worker{
		ID:   id,
		Job:  make(chan func()),
		exit: make(chan struct{}),
	}
}

func (w *worker) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case job, ok := <-w.Job:
				if !ok {
					return
				}
				job()
			case <-w.exit:
				return
			}
		}
	}()
}

func (w *worker) Stop() {
	close(w.exit)
}

type Pool struct {
	workers []*worker
	jobChan chan func()
	wg      sync.WaitGroup
}

func NewPool(size int) *Pool {
	p := &Pool{
		workers: make([]*worker, size),
		jobChan: make(chan func()),
	}
	for i := 0; i < size; i++ {
		p.workers[i] = NewWorker(i)
	}
	return p
}

func (p *Pool) Start() {
	for _, w := range p.workers {
		w.Start(&p.wg)
	}
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			select {
			case job, ok := <-p.jobChan:
				if !ok {
					return
				}
				w := p.getAvailableWorker()
				if w == nil {
					job()
				} else {
					w.Job <- job
				}
				//case <-p.:
				//	return
			}
		}
	}()
}

func (p *Pool) getAvailableWorker() *worker {
	for _, w := range p.workers {
		select {
		case <-w.exit:
		default:
			return w
		}
	}
	return nil
}

func (p *Pool) Stop() {
	close(p.jobChan)
	p.wg.Wait()
	for _, w := range p.workers {
		w.Stop()
	}
}

func main() {
	// 初始化协程池，设置最大协程数为3
	p := NewPool(3)

	// 启动协程池，开始接收任务
	p.Start()

	// 添加任务到协程池中
	for i := 0; i < 10; i++ {
		id := i
		p.jobChan <- func() {
			fmt.Println("Worker ", id, " started job")
			// 执行任务
			fmt.Println("Worker ", id, " finished job")
		}
	}

	// 停止协程池
	p.Stop()
}
