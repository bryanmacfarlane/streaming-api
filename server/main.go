package main

import (
	"concgroup"
	"context"
	"data"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const MAX_CONCURRENCY = 5

func main() {
	r := chi.NewRouter()
	r.Use(Chunked)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/foobars", func(w http.ResponseWriter, r *http.Request) {
		cw := NewChunkedWriter(w)

		cg, _ := concgroup.WithOptions(context.Background(), MAX_CONCURRENCY, func(identifier string, obj interface{}, err error) {
			cw.Send(NewChunk(identifier, obj, err))
		})

		cg.Go("foo", func() (interface{}, error) {
			return data.GetFoo(1)
		})

		bazId := 0
		cg.Go("bar", func() (interface{}, error) {
			bar, err := data.GetBar(2)
			if err != nil {
				return nil, err
			}
			bazId = bar.BazId
			return bar, err
		})

		cg.Wait()

		baz, err := data.GetBaz(bazId)
		cw.Send(NewChunk("baz", baz, err))
		// cw.Done()
	})
	http.ListenAndServe(":3000", r)
}
