package main

import (
	"encoding/json"
	"net/http"
)

type Chunk struct {
	Identifier string      `json:"identifier"`
	Obj        interface{} `json:"data"`
	Error      string      `json:"error"`
}

func NewChunk(identifier string, obj interface{}, err error) Chunk {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	return Chunk{Identifier: identifier, Obj: obj, Error: errMsg}
}

func Chunked(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Transfer-Encoding", "chunked")
		next.ServeHTTP(w, r)
	})
}

var lineDelimiter = []byte("\n\n")

type ChunkedWriter struct {
	w http.ResponseWriter
}

func NewChunkedWriter(w http.ResponseWriter) *ChunkedWriter {
	return &ChunkedWriter{w: w}
}

func (cw *ChunkedWriter) Send(obj interface{}) {
	b, _ := json.MarshalIndent(obj, "", "  ")
	_, _ = cw.w.Write(b)
	_, _ = cw.w.Write(lineDelimiter)
	cw.w.(http.Flusher).Flush()
}
