package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ereminIvan/fffb/app"
)

var configPath, fbDumpPath, tgDumpPath string

func init() {
	flag.StringVar(&configPath,"config", "./config.json", "Path to config")
	flag.StringVar(&fbDumpPath, "fb_dump", "./fb.dump", "Path to dump file of old facebook message")
	flag.StringVar(&tgDumpPath, "tg_dump", "./tg.dump", "Path to dump file telegram subscriber chats")

	flag.Parse()
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})

	a, err := app.Init(configPath, fbDumpPath, tgDumpPath)
	if err != nil {
		panic(err)
	}

	go func() {
		sig := <-signals
		log.Printf("Interupt signal %v", sig)
		a.Stop()
		done <- struct{}{}
	}()

	a.Run(done)

	close(signals)
	close(done)
}
