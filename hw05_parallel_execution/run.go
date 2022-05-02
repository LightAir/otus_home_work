package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	maximumErrors := int32(m)
	errorsCounter := int32(0)

	queue := make(chan Task)

	if len(tasks) < n {
		n = len(tasks)
	}

	wg := &sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go consumer(queue, wg, &errorsCounter)
	}

	for _, task := range tasks {
		if isMaxErr(atomic.LoadInt32(&errorsCounter), maximumErrors) {
			break
		}
		queue <- task
	}

	close(queue)
	wg.Wait()

	if isMaxErr(errorsCounter, maximumErrors) {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func isMaxErr(errorsCounter int32, maximumErrors int32) bool {
	return errorsCounter >= maximumErrors && maximumErrors > 0
}

func consumer(queue chan Task, wg *sync.WaitGroup, errorsCounter *int32) {
	defer wg.Done()

	for task := range queue {
		err := task()
		if err != nil {
			atomic.AddInt32(errorsCounter, 1)
		}
	}
}
