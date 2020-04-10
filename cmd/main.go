package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardaguclu/pexcutor"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		for {
			select {
			case <-sigint:
				cancel()
			}
		}
	}()

	p := pexcutor.New(ctx, 3, "ls", "-alh")
	err := p.Start()
	if err != nil {
		log.Fatal("start error ", err)
	}

	sOut, sErr, err := p.GetResult()
	if err != nil {
		log.Fatal("get result error ", err)
	}

	log.Println(sOut)
	log.Println(sErr)
}
