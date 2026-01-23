package workers

import (
	"sync"
)

type Result[T any] struct {
	Ok  T
	Err error
}

type Pool[I any, O any] struct {
	workers int
	process func(I) (O, error)
}

func NewPool[I any, O any](workers int, processingFunction func(I) (O, error)) Pool[I, O] {
	return Pool[I, O]{
		workers: workers,
		process: processingFunction,
	}
}

func (p Pool[I, O]) Process(jobs []I) ([]O, error) {
	in := make(chan I, len(jobs))
	out := make(chan Result[O], p.workers)

	cancel := make(chan struct{})
	var wg sync.WaitGroup
	defer wg.Wait()

	for _, job := range jobs {
		in <- job
	}
	close(in)

	for range p.workers {
		wg.Go(func() {
			for job := range in {
				select {
				case <-cancel:
					return
				default:
				}

				result, err := p.process(job)
				out <- Result[O]{Ok: result, Err: err}

				if err != nil {
					return
				}
			}
		})
	}

	processResult := make([]O, 0, len(jobs))
	for range len(jobs) {
		jobResult := <-out
		if jobResult.Err != nil {
			close(cancel)
			return nil, jobResult.Err
		}

		processResult = append(processResult, jobResult.Ok)
	}

	return processResult, nil
}
