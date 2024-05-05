package main

import (
	"context"
	"encoding/json"
	"fmt"

	"data"

	"concgroup"
)

type Response struct {
	Identifier string      `json:"identifier"`
	Data       interface{} `json:"data"`
	Error      string      `json:"error"`
}

func send(identifier string, data interface{}, err error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	res := Response{Identifier: identifier, Data: data, Error: errMsg}
	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))
}
func main() {
	limit := 3
	bazId := 0
	cg := concgroup.WithOptions(context.Background(), limit, func(identifier string, data interface{}, err error) {
		send(identifier, data, err)
	})

	cg.Go("foo", func() (interface{}, error) {
		return data.GetFoo(1)
	})

	cg.Go("bar", func() (interface{}, error) {
		bar, err := data.GetBar(2)
		if err != nil {
			return nil, err
		}

		bazId = bar.BazId

		return bar, nil
	})

	cg.Wait()

	baz, err := data.GetBaz(bazId)
	send("baz", baz, err)
}
