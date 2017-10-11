package main

import (
	"github.com/ereminIvan/fffb/app"
	"time"
)

func main() {
	a, err := app.Init()
	if err != nil {
		panic(err)
	}

	for {
		a.Run()
		time.Sleep(1 * time.Minute)
	}
}
