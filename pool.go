package superpool

// Blocks if the eventChan is full
func (p *Pool[T]) Add(e T) {
	p.wg.Add(1)
	p.eventChan <- e
}

func (p *ReturnPool[T, R]) Add(e T) <-chan R {
	p.wg.Add(1)
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

// Needs a separate shutdown method becuase eventChan is different from default pool
func (p *ReturnPool[T, R]) Shutdown() {
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

// Waits until all tasks in the pool are completed
func (p *Pool[T]) Wait() {
	p.wg.Wait()
}

func (p *Pool[T]) PendingTasks() int {
	return len(p.eventChan)
}

func (p *ReturnPool[T, R]) PendingTasks() int {
	return len(p.eventChan)
}
