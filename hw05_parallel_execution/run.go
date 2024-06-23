package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	taskCh := make(chan Task)

	go func() {
		for _, task := range tasks {
			taskCh <- task
		}
		close(taskCh)
	}()

	errorTaskCount := 0
	// Place your code here.
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range taskCh {
				if errorTaskCount != m || len(taskCh) != 0 {
					err := task()
					if err != nil {
						errorTaskCount++
					}
				}
			}
		}()
	}

	wg.Wait()
	switch {
	case errorTaskCount != 0:
		return ErrErrorsLimitExceeded
	default:
		return nil
	}
}
