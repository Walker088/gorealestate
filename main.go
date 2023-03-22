package main

import (
	"fmt"
	"os"

	"github.com/Walker088/gorealestate/config"
	"github.com/Walker088/gorealestate/crawler/plvr"
	"github.com/Walker088/gorealestate/database"
	"github.com/Walker088/gorealestate/logger"
)

const (
	envFile = ".env.development"
)

func main() {
	p, _ := os.Getwd()
	c, err := config.New(fmt.Sprintf("%s/%s", p, envFile))
	if err != nil {
		fmt.Println(err.ToString())
		os.Exit(1)
	}
	l := logger.New(c.GetLoggerConfig())
	defer l.Sync()

	pool, err := database.New(c.GetPgConfig(), l)
	if err != nil {
		l.DPanicf("Init DB Conn Pool error: %w", err)
	}
	defer pool.ShutDownPool()

	l.Info("Welcome to GoRealEstate")

	crawler := plvr.New(p, l, pool)
	crawler.Run()
	crawler.Stop()
}
