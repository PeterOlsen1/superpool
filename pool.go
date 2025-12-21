package superpool

// Blocks if the eventChan is full
func (p *Pool[T]) Add(e T) {
	p.eventChan <- e
}

func (p *ReturnPool[T, R]) Add(e T) <-chan R {
	task := Task[T, R]{
		Input:  e,
		Result: make(chan R, 1),
	}
	p.eventChan <- task
	return task.Result
}

func (p *Pool[T]) Shutdown() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for range p.numWorkers {
		p.quitChan <- struct{}{}
	}
	close(p.errors)
	p.wg.Wait()

	close(p.quitChan)
	close(p.eventChan)
}

func (p *Pool[T]) UpdateEventHandler(f EventHandler[T]) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.eventHandler = f
}

func (p *ReturnPool[T, R]) UpdateEventHandler(f ReturnEventHandler[T, R]) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.eventHandler = f
}

// Kill n workers
func (p *Pool[T]) KillN(n uint16) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.unsafeKillN(n)
}

// Kill n workers. No locking
func (p *Pool[T]) unsafeKillN(n uint16) {
	if n > p.numWorkers {
		n = p.numWorkers
	}

	p.numWorkers -= n
	for range n {
		p.quitChan <- struct{}{}
	}
}

// Resizes the current worker pool. If newSize < curSize, kill (cur - new) threads.
// Else, start (new - cur) threads.
func (p *Pool[T]) Resize(newSize uint16) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if newSize < p.numWorkers {
		p.unsafeKillN(p.numWorkers - newSize)
	} else {
		p.numWorkers = newSize
		for range newSize - p.numWorkers {
			p.startWorker()
		}
	}
}
