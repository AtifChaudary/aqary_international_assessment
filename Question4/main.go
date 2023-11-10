package main

import (
	"sync"
)

func main() {
	M := 8
	N := 2

	// Shared buffer
	buffer := make([]byte, 0)

	// Channels for communication between readers and writers
	writeCh := make(chan struct{}, 1)
	readCh := make(chan struct{}, 1)

	// Use a mutex for protecting the shared buffer
	var mutex sync.Mutex

	// Start N writer goroutines
	for i := 0; i < N; i++ {
		go func() {
			for {
				// Signal that writing is about to start
				writeCh <- struct{}{}

				// Acquire the lock
				mutex.Lock()

				// Write to the shared buffer
				// (Simulated write operation)
				buffer = append(buffer, 'A')

				// Release the lock
				mutex.Unlock()

				// Signal that writing is done
				<-writeCh
			}
		}()
	}

	// Start M reader goroutines
	for i := 0; i < M; i++ {
		go func() {
			for {
				// Signal that reading is about to start
				readCh <- struct{}{}

				// Acquire the lock
				mutex.Lock()

				// Read from the shared buffer
				// (Simulated read operation)
				_ = buffer

				// Release the lock
				mutex.Unlock()

				// Signal that reading is done
				<-readCh
			}
		}()
	}

	// Keep the program running
	select {}
}
