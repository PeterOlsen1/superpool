package superpool

import "fmt"

func (p *Pool[T]) Add(e T) error {
	if len(p.eventChan) >= cap(p.eventChan) {
		return fmt.Errorf("event channel is at capacity")
	}

	p.eventChan <- e
	return nil
}

func (p *Pool[T]) Close() {
	close(p.eventChan)
	p.wg.Wait()
}

func (p *Pool[T]) UpdateEventHandler(f EventHandler[T]) {
	p.eventHandler = f
}
