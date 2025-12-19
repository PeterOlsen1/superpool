package superpool

// Returns the channel of errors for the user to listen
func (p *Pool[T]) Errors() <-chan error {
	return p.errors
}

// Allows the user to pass in an error-handling function to apply to all errors
func (p *Pool[T]) HandleErrors(handler func(error)) {
	p.wg.Add(1)

	go func() {
		for err := range p.errors {
			handler(err)
		}
		p.wg.Done()
	}()
}
