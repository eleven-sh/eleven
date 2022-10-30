package queues

import "sync"

type Infrastructure interface{}

// InfrastructureQueue represents the queue used by cloud providers
// to modify an infrastructure in an ordered manner.
//
// A queue is defined as a two dimensional slice where each element
// is a slice of steps that will run at the same time.
//
// Example: [[funcA1, funcA2, funcA3], [funcB1, funcB2]]
// will run funcA1, funcA2, and funcA3 at the same time,
// BEFORE executing funcB1 and funcB2 also at the same time.
type InfrastructureQueue[T Infrastructure] []InfrastructureQueueSteps[T]

type InfrastructureQueueSteps[T Infrastructure] []InfrastructureQueueStep[T]

type InfrastructureQueueStep[T Infrastructure] func(infrastructure T) error

func (queue InfrastructureQueue[T]) Run(infrastructure T) error {
	for _, steps := range queue {
		if len(steps) == 0 {
			continue
		}

		stepErrorsChan := make(chan error, len(steps))
		var stepsWaiter sync.WaitGroup

		stepsWaiter.Add(len(steps))

		for _, step := range steps {
			go func(step InfrastructureQueueStep[T]) {
				defer stepsWaiter.Done()
				stepErrorsChan <- step(infrastructure)
			}(step)
		}

		stepsWaiter.Wait()
		close(stepErrorsChan)

		for stepError := range stepErrorsChan {
			if stepError != nil {
				return stepError
			}
		}
	}

	return nil
}
