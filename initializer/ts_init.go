package initializer

import (
	"context"
	"log"
	"encoding/csv"
   "fmt"
	//"strings"

	"github.com/520lly/qt_data_service/clients"
	"github.com/520lly/qt_data_service/collector"
	"github.com/520lly/qt_data_service/config"
	"github.com/520lly/qt_data_service/models"
	"github.com/520lly/qt_data_service/storage"
   mycsv "github.com/520lly/qt_data_service/storage/csv"
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
      log.Println("key-val:", key, *val)
      sub := (*val).GetSubscribe()
      fullPath := (*val).GetFullPath()
      //Check if update basic info. condition reached
      ret, files := strategies.CheckUpdateBasic(sub.Period.StockBasic, sub.StockBasic, fullPath)
      if ret {
         tse := strategies.TsEvent{models.DataFlag_Stock_Basic, nil}
         (*(*tsc.TsStrategies)[key]).TsEvent <- tse
      } else {
         //Load the Basic stock infor and start to download history data
         tsc.CheckUpdateDaily(files[0])
      }
      ret, files = strategies.CheckUpdateBasic(sub.Period.CompanyBasic, sub.CompanyBasic, fullPath)
      if ret {
         tse := strategies.TsEvent{models.DataFlag_Stock_Company, nil}
         (*(*tsc.TsStrategies)[key]).TsEvent <- tse
      } else {
         //Load the Basic stock infor and start to download history data
         //tsc.CheckUpdateDaily(files[0])
      }
   }
}

func (tsc *TsContext) CheckUpdateDaily(f string) {
   log.Printf("files:[%v]",f)
   _, fp := utils.OpenCsvFile(f)
   if fp != nil {
      var csvReader *csv.Reader
      csvReader = csv.NewReader(fp)
      records, err:= csvReader.ReadAll()
      if err == nil {
         for _, r := range records[1:] {
            log.Printf("%T:%v", r, r)
            sym := r[1]
            listDate := r[10]
            today, ts := utils.GetTodayString("20060102")
            _ = today
            historyFileNameTillToday := fmt.Sprintf("%s_%s-%s.csv", sym, listDate, ts)
            isNew, csvFp := utils.OpenCsvFile(historyFileNameTillToday)
            log.Printf("isNew:%v csvFp: %v", isNew, csvFp)
            tse := strategies.TsEvent{models.DataFlag_Trace_Daily, &map[string]string{"ts_code":"000001.SZ", "start_date":listDate, "end_date":ts}}
            (*(*tsc.TsStrategies)["sh"]).TsEvent <- tse
            break
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
			sto = mycsv.NewCsvStorage(*ctx, cfg.Store, v, home)
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
