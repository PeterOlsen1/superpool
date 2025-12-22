package superpool

import "sync"

type EventHandler[T any] func(T) error
type ReturnEventHandler[T, R any] func(T) (R, error)

type PoolInterface[T any] interface {
	Shutdown()

	Add(e T)
	UpdateEventHandler(f EventHandler[T])

	Resize(newSize uint16)
	KillN(n uint16)

	Errors()
	HandleErrors(handler func(error))
}

type Pool[T any] struct {
	// The number of workers to be initialized in the pool.
	//
	// Each of these numWorkers corresponds to 1 goroutine,
	// so keep that in mind when choosing what to put for this value
	numWorkers uint16

	// The channel on which events are sent
	//
	// Capacity of this chan is static, 1000.
	// Keep this number in mind when deciding how many worker threads to spawn
	eventChan chan T

	// The function to run when events are dequeued from the eventChan
	//
	// Should log its own errors
	eventHandler EventHandler[T]

	// Capacity of the event channel
	cap uint32

	// wait group for polling if pool is empty
	wg sync.WaitGroup

	// Mutex for modifying struct properties
	mu sync.Mutex

	// Send errors to the user
	errors chan error

	// Channel to quit a certain # of goroutines.
	// Used in dynamic pool resizing
	quitChan chan struct{}
}

type Task[T, R any] struct {
	Input  T
	Result chan R
}

// Default pool implemented with return values
//
// Need to overrite any functions that use eventHandler/Chan
type ReturnPool[T, R any] struct {
	Pool[T]
	eventHandler ReturnEventHandler[T, R]
	eventChan    chan Task[T, R]
}
