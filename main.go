package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/520lly/qt_data_service/config"
	"github.com/520lly/qt_data_service/initializer"
	"github.com/jinzhu/configor"

	"os/signal"
	"syscall"
	"time"
)

var (
	cfg  config.Config
	ctx  context.Context
	tsc  *initializer.TsContext
	home string
)

func usage() {
	fmt.Fprintf(os.Stderr, `quantitative trading data service version: v0.0.1
Usage: market_data_collector [-h] [-c config.json]
Options:
`)
	flag.PrintDefaults()
}

func main() {
	var c string
	var help bool
	flag.StringVar(&c, "c", "config.yml", "set configuration `json/yml file`")
	flag.BoolVar(&help, "h", false, "this help")
	flag.Usage = usage

	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	err := configor.Load(&cfg, c)

	if err != nil {
		panic(err)
	}

	if cfg.Store.Csv == cfg.Store.InfluxDB {
		panic("currently only support csv, please check your configure")
	}
	home = os.Getenv("HOME")

	ctx, cancel := context.WithCancel(context.Background())
	tsc := initializer.TsInit(&ctx, &cfg, home)
	log.Println("tsc %v\n", tsc)

	exitSignal := make(chan os.Signal, 1)
	sigs := []os.Signal{os.Interrupt, syscall.SIGILL, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGTERM}
	signal.Notify(exitSignal, sigs...)
	<-exitSignal
	cancel()
	time.Sleep(time.Second)
	log.Println("quantitative data collector exit")
}
