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

	// unbuffered, events will be received
	p.quitChan = make(chan struct{})

	// initialize worker threads
	for range p.numWorkers {
		p.startWorker()
	}
}

func (p *Pool[T]) startWorker() {
	p.wg.Add(1)

	go func() {
		for {
			select {
			case e := <-p.eventChan:
				err := p.eventHandler(e)
				if err != nil {
					p.errors <- err
				}
			case <-p.quitChan:
				p.wg.Done()
				return
			}
		}
	}()
}

// Duplicate methods for a pool with return values

func NewReturnPool[T, R any](cap uint32, numWorkers uint16, eventHandler ReturnEventHandler[T, R]) (*ReturnPool[T, R], error) {
	if cap == 0 || numWorkers == 0 {
		return nil, fmt.Errorf("parameters must be nonzero")
	}

	pool := ReturnPool[T, R]{
		Pool: Pool[T]{
			numWorkers: numWorkers,
			cap:        cap,
		},
		eventHandler: eventHandler,
	}

	pool.startPool()
	return &pool, nil
}

func (p *ReturnPool[T, R]) startPool() {
	p.eventChan = make(chan Task[T, R], p.cap)
	p.errors = make(chan error)

	// unbuffered, events will be received
	p.quitChan = make(chan struct{})

	// initialize worker threads
	for range p.numWorkers {
		p.startWorker()
	}
}

func (p *ReturnPool[T, R]) startWorker() {
	p.wg.Add(1)

	go func() {
		for {
			select {
			case e := <-p.eventChan:
				ret, err := p.eventHandler(e.Input)
				if err != nil {
					p.errors <- err
					close(e.Result)
				} else {
					e.Result <- ret
					close(e.Result)
				}
			case <-p.quitChan:
				p.wg.Done()
				return
			}
		}
	}()
}
