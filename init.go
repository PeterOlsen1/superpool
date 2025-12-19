package superpool

import "fmt"

func NewPool[T any](cap uint32, numWorkers uint16, eventHandler EventHandler[T]) (*Pool[T], error) {
	if cap == 0 || numWorkers == 0 {
		return nil, fmt.Errorf("parameters must be nonzero")
	}

	pool := Pool[T]{
		numWorkers:   numWorkers,
		eventHandler: eventHandler,
		cap:          cap,
	}

	pool.startPool()

	return &pool, nil
}

func (p *Pool[T]) startPool() {
	p.eventChan = make(chan T, p.cap)
	p.errors = make(chan error)

	// initialize worker threads
	for range p.numWorkers {
		p.startWorker()
	}
}

func (p *Pool[T]) startWorker() {
	p.wg.Add(1)

	go func() {
		for e := range p.eventChan {
			p.eventHandler(e)
		}
		p.wg.Done()
	}()
}
