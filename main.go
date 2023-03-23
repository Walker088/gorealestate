package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Walker088/gorealestate/config"
	"github.com/Walker088/gorealestate/crawler/plvr"
	"github.com/Walker088/gorealestate/database"
	"github.com/Walker088/gorealestate/logger"
	"github.com/Walker088/gorealestate/migrations"
)

const (
	envFile = ".env.development"
)

func main() {
	deadlineChannel := make(chan os.Signal, 1)
	signal.Notify(deadlineChannel, os.Interrupt)

	rootDir, _ := os.Getwd()
	c, err := config.New(fmt.Sprintf("%s/%s", rootDir, envFile))
	if err != nil {
		fmt.Println(err.ToString())
		os.Exit(1)
	}
	l := logger.New(c.GetLoggerConfig())
	defer l.Sync()

	sm, err := migrations.New(rootDir, c.GetPgConfig(), l)
	if err != nil {
		l.DPanic(err.ToString())
	}
	defer sm.Stop()
	if err := sm.Migrate(); err != nil {
		l.DPanic(err.ToString())
	}

	pool, err := database.New(c.GetPgConfig(), l)
	if err != nil {
		l.DPanicf("init db conn pool error: %w", err)
	}
	defer pool.ShutDownPool()

	l.Info("welcome to gorealestate")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	crawler := plvr.New(ctx, cancel, rootDir, l, pool.GetPool())
	go crawler.Start()
	for {
		select {
		case <-deadlineChannel:
			l.Info("Interrupt signal received")
			crawler.Stop()
			return
		case <-ctx.Done():
			l.Infof("Download finished")
			return
		case downloaded := <-crawler.ResultsCh:
			l.Infof("[main] Downloaded: %s", downloaded)
		}
	}

}
