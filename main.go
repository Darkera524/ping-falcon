package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func main(){
	cfg := flag.String("c", "cfg.json", "configuration file")

	ParseConfig(*cfg)

	InitLog()

	go CronConfig(60, *cfg)

	Init()

	GetHost()
	go CronHost()
	go Ping()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		DB.Close()
		os.Exit(0)
	}()

	select{}
}
