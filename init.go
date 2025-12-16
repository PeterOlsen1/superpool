package superpool

func NewPool[T any](cap uint32, numWorkers uint16, eventHandler EventHandler[T]) *Pool[T] {
	pool := Pool[T]{
		numWorkers:   numWorkers,
		eventHandler: eventHandler,
		cap:          cap,
	}

	pool.startPool()

	return &pool
}

func (p *Pool[T]) startPool() {
	p.eventChan = make(chan T, p.cap)

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
