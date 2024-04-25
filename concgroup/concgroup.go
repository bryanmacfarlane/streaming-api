package concgroup

import (
	"context"
	"sync"
)

// similar to ErrGroup: https://cs.opensource.google/go/x/sync/+/refs/tags/v0.7.0:errgroup/errgroup.go
type token struct{}

type ConcGroup struct {
	wg   sync.WaitGroup
	sem  chan token
	comp func(identifier string, data interface{}, err error)
	mu   sync.Mutex
}

// WithOptions takes a completion callback which will be called if not nil
// on completion of each concurrent Go call
// the comp callback is synchronized so that only one at a time can run.
//
// limit limits the number of active goroutines in this group to at most n.
// A negative value indicates no limit.
//
// Any subsequent call to the Go method will block until it can add an active
// goroutine without exceeding the configured limit.
func WithOptions(ctx context.Context, limit int, comp func(identifier string, data interface{}, err error)) (*ConcGroup, context.Context) {
	// ctx, cancel := withCancelCause(ctx)
	// return &Group{cancel: cancel}, ctx
	// TODO: wrap cancel
	cg := &ConcGroup{}
	cg.sem = nil
	if limit > 0 {
		cg.sem = make(chan token, limit)
	}
	cg.comp = comp

	return cg, ctx
}

func (cg *ConcGroup) Wait() {
	cg.wg.Wait()
}

func (cg *ConcGroup) done() {
	if cg.sem != nil {
		<-cg.sem
	}
	cg.wg.Done()
}

// Go calls the given function in a new goroutine.
// It blocks until the new goroutine can be added without the number of
// active goroutines in the group exceeding the configured limit.
//
// if data is returned as an interface{} and completion callback is not nil
// the completion callback will be called synchronously
func (cg *ConcGroup) Go(identifier string, f func() (interface{}, error)) {
	if cg.sem != nil {
		cg.sem <- token{}
	}

	cg.wg.Add(1)
	go func() {
		defer cg.done()

		data, err := f()
		cg.mu.Lock()
		defer cg.mu.Unlock()
		if cg.comp != nil {
			cg.comp(identifier, data, err)
		}
	}()
}
