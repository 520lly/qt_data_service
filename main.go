package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	//"github.com/520lly/qt_data_service/clients"
	//"github.com/520lly/qt_data_service/collector"
	"github.com/520lly/qt_data_service/config"
	//"github.com/520lly/qt_data_service/models"
	//"github.com/520lly/qt_data_service/storage"
	//"github.com/520lly/qt_data_service/storage/csv"
	//"github.com/520lly/qt_data_service/storage/influxdb"
	//"github.com/520lly/qt_data_service/strategies"
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

	//stgs := make(map[string]strategies.StockStrategy)
	//for _, v := range cfg.Subs {
	//var sto storage.Storage
	//if cfg.Store.Csv {
	//sto = csv.NewCsvStorage(ctx, v.Market, v.ExchangeName, v.CurrencyPair, v.ContractType, v, home, cfg.Store.CsvCfg.Location)
	//}
	//if cfg.Store.InfluxDB {
	//sto = influxdb.NewInfluxdb(ctx, v.ExchangeName, v.CurrencyPair, v.ContractType, cfg.Store.InfluxDbCfg.Url, cfg.Store.InfluxDbCfg.Database, cfg.Store.InfluxDbCfg.Username, cfg.Store.InfluxDbCfg.Password)
	//}
	//go sto.SaveWorker()
	//cl := &clients.TsClient{}
	//cl = clients.NewTsClient(v.Market, v.ExchangeName, v.CurrencyPair, cfg.Tokens.TuShare)
	//stg := strategies.NewStockStrategy(&cfg)
	//stgs[v.ExchangeName] = stg
	//// TODO Initializing collector worker
	//collector.NewTsCollector(ctx, cl, stg, sto)
	//}
	//stgs["sh"].TsEvent <- models.DataFlag_Stock_Basic
	//stgs["sz"].TsEvent <- models.DataFlag_Stock_Basic
	//stgs["sh"].TsEvent <- models.DataFlag_Stock_Company
	//stgs["sz"].TsEvent <- models.DataFlag_Stock_Company

	exitSignal := make(chan os.Signal, 1)
	sigs := []os.Signal{os.Interrupt, syscall.SIGILL, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGTERM}
	signal.Notify(exitSignal, sigs...)
	<-exitSignal
	cancel()
	time.Sleep(time.Second)
	log.Println("quantitative data collector exit")
}
