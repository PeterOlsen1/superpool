// package superpool
// Superpool is designed to make worker pools easy and simple.
//
// By allocating a set of goroutines to run compute / I/O bound tasks,
// a user can know that their tasks will be executed asynchronously.
//
// Features:
//   - Easy creation and management
//   - Asynchronous task execution
//   - Graceful shutdown
//
// Superpool is ideal for scenarios where you need to process many independent tasks concurrently,
// such as batch processing, web scraping, or parallel I/O operations.
package superpool
