# superpool

A generic, high-performance worker pool for Go, supporting both simple and return-value tasks. Designed for concurrent task processing with configurable worker count, error handling, and extensibility.

## Installation

```sh
go get github.com/PeterOlsen1/superpool
```

## Usage

### 1. Basic Pool (No Return Value)

```go
package main

import (
    "fmt"
    "github.com/PeterOlsen1/superpool"
)

func myTask(n int) error {
    fmt.Println("Processing:", n)
    return nil
}

func main() {
    pool, err := superpool.NewPool[int](100, 10, myTask)
    if err != nil {
        panic(err)
    }
    defer pool.Shutdown()

    for i := 0; i < 50; i++ {
        pool.Add(i)
    }

    pool.Wait() // Wait for all tasks to finish
}
```

### 2. Return Pool (With Return Value)

```go
package main

import (
    "fmt"
    "github.com/PeterOlsen1/superpool"
)

func square(n int) (int, error) {
    return n * n, nil
}

func main() {
    pool, err := superpool.NewReturnPool[int, int](100, 10, square)
    if err != nil {
        panic(err)
    }
    defer pool.Shutdown()

    resultChan := pool.Add(5)
    for result := range resultChan {
        fmt.Println("Result:", result) // Output: Result: 25
    }
}
```

### 3. Error Handling

Both pool types provide an `Errors()` channel to receive errors from tasks:

```go
go func() {
    for err := range pool.Errors() {
        fmt.Println("Task error:", err)
    }
}()
```

### 4. Dynamic Resizing

You can resize the worker pool at runtime:

```go
err := pool.Resize(20) // Change to 20 workers
if err != nil {
    fmt.Println("Resize error:", err)
}
```

## Benchmarks

See `benchmarks/benchmarks_test.go` for performance comparisons with goroutines and sequential execution.
