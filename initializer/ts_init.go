package initializer

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/520lly/qt_data_service/clients"
	"github.com/520lly/qt_data_service/collector"
	"github.com/520lly/qt_data_service/config"
	"github.com/520lly/qt_data_service/models"
	"github.com/520lly/qt_data_service/storage"
	mycsv "github.com/520lly/qt_data_service/storage/csv"
	//"github.com/520lly/qt_data_service/storage/influxdb"
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
		//log.Println("key-val:", key, *val)
		sub := (*val).GetSubscribe()
		fullPath := (*val).GetFullPath()
		//Check if update basic info. condition reached
		ret, files := strategies.CheckUpdateBasic(sub.Period.StockBasic, sub.StockBasic, fullPath)
		if ret {
			tse := strategies.TsEvent{models.DataFlag_Stock_Basic, nil, nil}
			(*(*tsc.TsStrategies)[key]).TsEvent <- tse
		} else {
			//Load the Basic stock infor and start to download history data
			tsc.CheckUpdateDaily(sub.ExchangeName, files[0])
		}
		ret, files = strategies.CheckUpdateBasic(sub.Period.CompanyBasic, sub.CompanyBasic, fullPath)
		if ret {
			tse := strategies.TsEvent{models.DataFlag_Stock_Company, nil, nil}
			(*(*tsc.TsStrategies)[key]).TsEvent <- tse
		} else {
			//Load the Basic stock infor and start to download history data
			//tsc.CheckUpdateDaily(files[0])
		}
	}
}

func (tsc *TsContext) CheckUpdateDaily(exchangeName string, f string) {
   log.Printf("files:[%v]", f)
   _, fp := utils.OpenCsvFile(f)
   if fp != nil {
      var csvReader *csv.Reader
      csvReader = csv.NewReader(fp)
      records, err := csvReader.ReadAll()
      if err == nil {
         for _, r := range records[1:] { //skip the firt line of symbols
            //log.Printf("%T:%v", r, r)
            sym := r[0]
            start := r[10] //listDate
            today, end := utils.GetTodayString("20060102")
            _ = today
            dir := fmt.Sprintf("%s", (*(*tsc.TsStorage)[exchangeName]).GetFullPath()+"/history/")
            log.Printf("%s", dir)
            start, end, csvFp := isDataUpdated(dir, sym, start)
            if csvFp == nil {
               log.Printf("FATAL ERROR: file [%s] open filed", sym)
               continue
            } else {
               tse := strategies.TsEvent{
                  models.DataFlag_Trace_Daily,
                  &map[string]string{"ts_code": r[0], "start_date": start, "end_date": end}, csvFp}

                  log.Printf("[%s]TsEvent:%v", exchangeName, tsc)
                  (*(*tsc.TsStrategies)[exchangeName]).TsEvent <- tse
               }
            }
         }
      }
   }

func isDataUpdated(dir string, sym string, listDate string) (start string, end string, fp *os.File) {
	fp = nil
	hisFile := fmt.Sprintf("%s.csv", dir+sym)
	today, ts := utils.GetTodayString("20060102")
	_ = today
   start = listDate
   end = ts
   isNew, fp := utils.OpenCsvFile(hisFile)
   log.Printf("[%s] isNew:%v csvFp: %v", hisFile, isNew, fp)
   if fp != nil {
      if isNew {
         utils.WriteData2CsvFile(fp, models.TradeDailyFieldSymbol)
         return start, end, fp
      } else {
         //old file and needs to be update
         csvr := csv.NewReader(fp)
         records, err := csvr.ReadAll()
         if err == nil {
            size := len(records)
            log.Printf("last line[%v]", records[size -1])
            lastUpdate := records[size - 1][1]
            diff := utils.CalcDateDiffByDay(lastUpdate, ts)
            log.Printf("diff = %d", diff)
            if diff > 1 { //at least is 2
               start_new := utils.AddDays2Date(lastUpdate, 0, 0, 1)
               start = start_new
            } else {
               fp.Close()
               fp = nil
            }
         }
      }
   } else {
      log.Printf("open file [%s] faild\n", hisFile)
   }

	return start, end, fp
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
			//sto = influxdb.NewInfluxdb(*ctx, v.ExchangeName, v.CurrencyPair, v.ContractType, cfg.Store.InfluxDbCfg.Url, cfg.Store.InfluxDbCfg.Database, cfg.Store.InfluxDbCfg.Username, cfg.Store.InfluxDbCfg.Password)
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
