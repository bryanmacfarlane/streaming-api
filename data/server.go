package data

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Foo struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type Bar struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
	BazId   int    `json:"bazid"`
}

type Baz struct {
	Id      int    `json:"id"`
	Address string `json:"address"`
}

const DEF_RAND_SLEEP = 10
const DEF_ERROR_DENOM = 4

func randSleep(ms int) error {
	rnd := rand.Intn(ms)
	if rnd >= ms-ms/DEF_ERROR_DENOM {
		return errors.New("timeout")
	}

	// dur := time.Duration(rnd)
	time.Sleep(time.Duration(rnd) * time.Millisecond)
	return nil
}

func GetFoo(id int) (*Foo, error) {
	err := randSleep(DEF_RAND_SLEEP)
	if err != nil {
		return nil, err
	}
	return &Foo{Id: id, Title: fmt.Sprintf("foo title for %v", id)}, nil
}

func GetBar(id int) (*Bar, error) {
	err := randSleep(DEF_RAND_SLEEP)
	if err != nil {
		return nil, err
	}
	return &Bar{Id: id, Message: fmt.Sprintf("bar message for %v", id), BazId: 3}, nil
}

func GetBaz(id int) (*Baz, error) {
	if id < 1 {
		return nil, errors.New("bad request")
	}

	err := randSleep(DEF_RAND_SLEEP)
	if err != nil {
		return nil, err
	}

	return &Baz{Id: id, Address: fmt.Sprintf("baz address for %v", id)}, nil
}
