package main

import (
	"github.com/ereminIvan/fffb/app"
)

func main() {
	a, err := app.Init()
	if err != nil {
		panic(err)
	}

	a.Run()
}
