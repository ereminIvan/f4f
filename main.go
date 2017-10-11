package main

import (
	//"os"
	//"os/signal"
	//"syscall"

	"github.com/ereminIvan/fffb/app"
)

func main() {
	a, err := app.Init()
	if err != nil {
		panic(err)
	}
	//
	//sigs := make(chan os.Signal, 1)
	//done := make(chan bool, 1)
	//
	//signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//
	//go func() {
	//	<-sigs
	//	done <- true
	//}()

	a.Run()

	//<-done

	//a.Finish()
}
