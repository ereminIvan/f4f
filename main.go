package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ereminIvan/fffb/app"
)

var configPath, fbDump string

func init() {
	flag.StringVar(&fbDump,"fb_dump", "./fb.dump", "Path to dump file of old facebook message")
	flag.StringVar(&configPath,"config", "./config.json", "Path to config")
	flag.Parse()
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	a, err := app.Init(configPath, fbDump)
	if err != nil {
		panic(err)
	}

	go func() {
		log.Printf("Recived interupt signal %v", <-signals)
		done <- true
	}()

	a.Run()

	<-done

	a.Stop()

	close(signals)
	close(done)
}
