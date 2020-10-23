package main

import(
	"context"
	"flag"
	"fmt"
	"log"
	"os"

   "github.com/520lly/qt_data_service/clients"
	"github.com/520lly/qt_data_service/config"
	"github.com/520lly/qt_data_service/storage"
	"github.com/520lly/qt_data_service/storage/csv"
	"github.com/520lly/qt_data_service/storage/influxdb"
)

var (
	cfg config.Config
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

	ctx, cancel := context.WithCancel(context.Background())
	for _, v := range cfg.Subs {
		var sto storage.Storage
		if cfg.Store.Csv {
			sto = csv.NewCsvStorage(ctx, v.ExchangeName, v.CurrencyPair, v.ContractType, v.Flag, "output/csv", "output/tar")
		}
		if cfg.Store.InfluxDB {
			sto = influxdb.NewInfluxdb(ctx, v.ExchangeName, v.CurrencyPair, v.ContractType, cfg.Store.InfluxDbCfg.Url, cfg.Store.InfluxDbCfg.Database, cfg.Store.InfluxDbCfg.Username, cfg.Store.InfluxDbCfg.Password)
		}
		go sto.SaveWorker()
		cl := &client.McClient{}
   }

	exitSignal := make(chan os.Signal, 1)
	sigs := []os.Signal{os.Interrupt, syscall.SIGILL, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGTERM}
	signal.Notify(exitSignal, sigs...)
	<-exitSignal
	cancel()
	time.Sleep(time.Second)
	log.Println("market data collector exit")
}
