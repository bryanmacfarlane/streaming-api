package concgroup

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConcGroup(t *testing.T) {
	results := make([]string, 0)

	cg := WithOptions(context.Background(), 1, func(identifier string, obj string, err error) {
		results = append(results, identifier)
	})

	cg.Go("foo", func() (string, error) { return "foo", nil })
	cg.Go("bar", func() (string, error) { return "bar", nil })

	cg.Wait()

	require.Contains(t, results, "foo")
}

func TestConcurrencyLimit(t *testing.T) {
	results := make([]string, 0)

	cg := WithOptions(context.Background(), 1, func(identifier string, obj string, err error) {
		results = append(results, identifier)
	})

	cg.Go("foo", func() (string, error) {
		time.Sleep(1 * time.Second)
		return "foo", nil
	})
	// Go blocks on concurrency limit, so lets put it in a goroutine
	// so we can validate that the limit is respected
	go cg.Go("bar", func() (string, error) { return "bar", nil })
	time.Sleep(10 * time.Millisecond)
	require.Empty(t, results)
}
