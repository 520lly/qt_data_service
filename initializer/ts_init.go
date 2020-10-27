package initializer

import (
	"context"
	"log"
	"regexp"

	"github.com/520lly/qt_data_service/clients"
	"github.com/520lly/qt_data_service/collector"
	"github.com/520lly/qt_data_service/config"
	"github.com/520lly/qt_data_service/models"
	"github.com/520lly/qt_data_service/storage"
	"github.com/520lly/qt_data_service/storage/csv"
	"github.com/520lly/qt_data_service/storage/influxdb"
	"github.com/520lly/qt_data_service/strategies"
	"github.com/520lly/qt_data_service/utils"
)

type TsContext struct {
	TsClients    *map[string]*clients.TsClient
	TsStorage    *map[string]*storage.Storage
	TsStrategies *map[string]*strategies.StockStrategy
	//TsCollector  *map[string]*collector.TsCollector
}

func (tsc *TsContext) InitStockStrategies() {
	if tsc == nil {
		log.Println("tsc is nil!")
		return
	}
	for key, val := range *tsc.TsStorage {
		log.Println("%v:%v", key, *val)
		//Check StockBasic csv file exist or not
		err, files := utils.FilteredSearchOfDirectoryTree(regexp.MustCompile((*val).GetStockBasic()), (*val).GetFullPath())
		if err == nil {
			log.Println("tsc.TsStrategies %v\n", (*(*tsc.TsStrategies)[key]).TsEvent)
			if len(files) <= 0 {
				tse := strategies.TsEvent{models.DataFlag_Stock_Basic, nil}
				(*(*tsc.TsStrategies)[key]).TsEvent <- tse
				//stgs["sz"].TsEvent <- models.DataFlag_Stock_Basic
				//stgs["sh"].TsEvent <- models.DataFlag_Stock_Company
				//stgs["sz"].TsEvent <- models.DataFlag_Stock_Company
			}
		}
	}
}

func TsInit(ctx *context.Context, cfg *config.Config, home string) *TsContext {
	ts_stgs := make(map[string]*strategies.StockStrategy)
	ts_clis := make(map[string]*clients.TsClient)
	ts_stos := make(map[string]*storage.Storage)

	for _, v := range cfg.Subs {
		var sto storage.Storage
		if cfg.Store.Csv {
			sto = csv.NewCsvStorage(*ctx, v.Market, v.ExchangeName, v.CurrencyPair, v.ContractType, v, home, cfg.Store.CsvCfg.Location)
		}
		if cfg.Store.InfluxDB {
			sto = influxdb.NewInfluxdb(*ctx, v.ExchangeName, v.CurrencyPair, v.ContractType, cfg.Store.InfluxDbCfg.Url, cfg.Store.InfluxDbCfg.Database, cfg.Store.InfluxDbCfg.Username, cfg.Store.InfluxDbCfg.Password)
		}
		go sto.SaveWorker()
		cl := &clients.TsClient{}
		cl = clients.NewTsClient(v.Market, v.ExchangeName, v.CurrencyPair, cfg.Tokens.TuShare)
		stg := strategies.NewStockStrategy(cfg)
		collector.NewTsCollector(ctx, cl, stg, &sto)
		ts_stgs[v.ExchangeName] = stg
		ts_clis[v.ExchangeName] = cl
		ts_stos[v.ExchangeName] = &sto
	}

	var tsc TsContext
	tsc.TsStrategies = &ts_stgs
	tsc.TsClients = &ts_clis
	tsc.TsStorage = &ts_stos

	tsc.InitStockStrategies()

	return &tsc
}
