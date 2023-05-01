package sync

import "golang.org/x/sync/errgroup"

// CreateWorkers spins up the specified amount of workers. If one of the workers
// emits an error, then that error will be returned.
func CreateWorkers(f func() error, n int) error {
	g := new(errgroup.Group)

	for i := 0; i < n; i++ {
		g.Go(f)
	}

	return g.Wait()
}
