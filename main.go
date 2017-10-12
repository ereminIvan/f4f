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
	//catch system signals for future graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	//chanel for graceful shutdown
	done := make(chan struct{})

	application, err := app.Init(configPath, fbDumpPath, tgDumpPath)
	if err != nil {
		panic(err)
	}

	go func() {
		sig := <-signals
		log.Printf("Interruption signal has received: %v", sig)
		application.Stop()
		done <- struct{}{}
	}()

	application.Run(done)

	close(signals)
	close(done)
}
