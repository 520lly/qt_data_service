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
   "github.com/520lly/qt_data_service/utils"
   "github.com/520lly/qt_data_service/strategies"
)

type TsContext struct {
   TsClients    *map[string]*clients.TsClient
   TsStorage    *map[string]*storage.Storage
   TsEvents     *map[string]*chan models.TsEvent
   TsCron       chan models.DataFlag
   //TsStrategy   *strategies.StockStrategy
   //TsCollector  *map[string]*collector.TsCollector
}

func (tsc *TsContext) CheckUpdateBasics() {
   if tsc == nil {
      log.Println("FATAL: tsc is nil!")
      return
   }
   for key, val := range *tsc.TsStorage {
      //log.Printf(" ========= key:%s", key)
      sub := (*val).GetSubscribe()
      fi, ok := (*val).GetBasicFile(sub.StockBasic) //BasicCsvFileInfo
      if ok {
         log.Printf("BasicCsvFileInfo:%v exist check need update or not", fi)
         if utils.IsDateAfterToday(utils.DateFormat2, fi.NextUpdate) {
            log.Println("Yes! update it!")
            _, fp := utils.OpenCsvFile(fi.FullPath)
            if fp != nil {
               //utils.EmptyCsvFileContentWithHeader(fp, models.StockFieldSymbol)
               //fp.Truncate(0)
               params := make(map[string]string)
               params["is_hs"] = ""
               params["list_status"] = ""
               params["exchange"] = sub.ExchangeName
               tse := models.TsEvent{models.DataFlag_Stock_Basic, &params, fp}
               (*(*tsc.TsEvents)[key])<- tse
            } else {
               log.Printf("File [%s] is empty!", fi.FullPath)
            }
         } else {
            log.Println("No! No need to update it!")
         }
      } else {
         log.Printf("Get [%s] failed!", sub.StockBasic)
      }

      fi, ok = (*val).GetBasicFile(sub.CompanyBasic) //BasicCsvFileInfo
      if ok {
         log.Printf("BasicCsvFileInfo:%v exist check need update or not", fi)
         if utils.IsDateAfterToday(utils.DateFormat2, fi.NextUpdate) {
            log.Println("Yes! update it!")
            _, fp := utils.OpenCsvFile(fi.FullPath)
            if fp != nil {
               //utils.EmptyCsvFileContentWithHeader(fp, models.CompanyFieldSymbol)
               //fp.Truncate(0)
               params := make(map[string]string)
               params["ts_code"] = ""
               params["exchange"] = sub.ExchangeName
               tse := models.TsEvent{models.DataFlag_Stock_Company, &params, fp}
               (*(*tsc.TsEvents)[key])<- tse
            } else {
               log.Printf("File [%s] is empty!", fi.FullPath)
            }
         } else {
            log.Println("No! No need to update it!")
         }
      } else {
         log.Printf("Get [%s] failed!", sub.CompanyBasic)
      }
   }
}

func (tsc *TsContext) CheckUpdateTradeCal() {
   if tsc == nil {
      log.Println("FATAL: tsc is nil!")
      return
   }
   for key, val := range *tsc.TsStorage {
      sub := (*val).GetSubscribe()
      fi, ok := (*val).GetBasicFile(sub.TradeCal) //BasicCsvFileInfo
      if ok {
         log.Printf("BasicCsvFileInfo:%v exist check need update or not", fi)
         if utils.IsDateAfterToday(utils.DateFormat2, fi.NextUpdate) {
            log.Println("Yes! update it!")
            _, fp := utils.OpenCsvFile(fi.FullPath)
            if fp != nil {
               _, ts := utils.GetTodayString(utils.DateFormat1)
               params := make(map[string]string)
               params["exchange"] = sub.ExchangeName
               params["start_date"] = "19890000"
               params["end_date"] = ts
               csvr := csv.NewReader(fp)
               records, err := csvr.ReadAll()
               if err == nil {
                  size := len(records)
                  if size < 2 {
                     log.Printf("%s records is empty")
                  } else {
                     log.Printf("last line[%v]", records[size -1])
                     lastUpdate := records[size - 1][1]
                     diff := utils.CalcDateDiffByDay(utils.DateFormat1, lastUpdate, ts)
                     log.Printf("diff = %d", diff)
                     if diff > 1 { //at least is 2
                        start_new := utils.AddDays2Date(utils.DateFormat1, lastUpdate, 0, 0, 1)
                        params["start_date"] = start_new
                     } else {
                        fp.Close()
                        fp = nil
                     }
                  }
               }
               params["is_open"] = "" //N=No, Y=Yes
               tse := models.TsEvent{models.DataFlag_Trade_Calendar, &params, fp}
               (*(*tsc.TsEvents)[key])<- tse
            } else {
               log.Printf("File [%s] is empty!", fi.FullPath)
            }
         } else {
            log.Println("No! No need to update it!")
         }
      }
   }
}

func (tsc *TsContext) CheckUpdateNameChangeHistory() {
   if tsc == nil {
      log.Println("FATAL: tsc is nil!")
      return
   }
   for key, val := range *tsc.TsStorage {
      sub := (*val).GetSubscribe()
      fi, ok := (*val).GetBasicFile(sub.NameChangeHistory) //BasicCsvFileInfo
      if ok {
         log.Printf("BasicCsvFileInfo:%v exist check need update or not", fi)
         if utils.IsDateAfterToday(utils.DateFormat2, fi.NextUpdate) {
            log.Println("Yes! update it!")
            _, fp := utils.OpenCsvFile(fi.FullPath)
            if fp != nil {
               //utils.EmptyCsvFileContentWithHeader(fp, models.NameChangeFieldSymbol)
               //fp.Truncate(0)
               params := make(map[string]string)
               params["ts_code"] = ""
               params["start_date"] = ""
               params["end_date"] = ""
               tse := models.TsEvent{models.DataFlag_NameChange_Histtory, &params, fp}
               (*(*tsc.TsEvents)[key])<- tse
            } else {
               log.Printf("File [%s] is empty!", fi.FullPath)
            }
         } else {
            log.Println("No! No need to update it!")
         }
      }
   }
}

//func (tsc *TsContext) CheckUpdateHsConst() {
   //if tsc == nil {
      //log.Println("FATAL: tsc is nil!")
      //return
   //}
   //for key, val := range *tsc.TsStorage {
      //sub := (*val).GetSubscribe()
      //fi, ok := (*val).GetBasicFile(sub.) //BasicCsvFileInfo
      //if ok {
         //log.Printf("BasicCsvFileInfo:%v exist check need update or not", fi)
         //if utils.IsDateAfterToday(utils.DateFormat2, fi.NextUpdate) {
            //log.Println("Yes! update it!")
            //_, fp := utils.OpenCsvFile(fi.FullPath)
            //if fp != nil {
               //utils.EmptyCsvFileContentWithHeader(fp, models.NameChangeFieldSymbol)
               //params := make(map[string]string)
               //params["ts_code"] = ""
               //params["start_date"] = ""
               //params["end_date"] = ""
               //tse := models.TsEvent{models.DataFlag_NameChange_Histtory, &params, fp}
               //(*(*tsc.TsEvents)[key])<- tse
            //} else {
               //log.Printf("File [%s] is empty!", fi.FullPath)
            //}
         //} else {
            //log.Println("No! No need to update it!")
         //}
      //}
   //}
//}

func (tsc *TsContext) CheckUpdateDaily() {
   if tsc == nil {
      log.Println("FATAL: tsc is nil!")
      return
   }
   for key, val := range *tsc.TsStorage {
      sub := (*val).GetSubscribe()
      fi, ok := (*val).GetBasicFile(sub.StockBasic) //BasicCsvFileInfo
      if ok {
         log.Printf("CheckUpdateDaily files:[%v]", fi)
         _, fp := utils.OpenCsvFile(fi.FullPath)
         if fp != nil {
            var csvReader *csv.Reader
            csvReader = csv.NewReader(fp)
            records, err := csvReader.ReadAll()
            if err == nil && len(records) > 1 {
               for _, r := range records[1:] { //skip the firt line of symbols
                  //log.Printf("%T:%v", r, r)
                  sym := r[0]
                  start := r[10] //listDate
                  today, end := utils.GetTodayString(utils.DateFormat1)
                  _ = today
                  dir := fmt.Sprintf("%s", (*(*tsc.TsStorage)[key]).GetFullPath()+"/history/")
                  log.Printf("%s", dir)
                  updated, start, end, csvFp := isDataUpdated(dir, sym, start)
                  if updated {
                     log.Printf("Info: file [%s] is update to date", sym)
                     continue
                  } else {
                     if csvFp == nil {
                        log.Printf("FATAL ERROR: file [%s] open filed", sym)
                        continue
                     }
                     tse := models.TsEvent { 
                        models.DataFlag_Trade_Daily,
                        &map[string]string{"ts_code": r[0], "start_date": start, "end_date": end},
                        csvFp}  
                     log.Printf("[%s]TsEvent:%v", key, tsc)
                     (*(*tsc.TsEvents)[key])<- tse
                  }
               }
            } else {
               log.Println("StockBasic is empty or is not existing")
            }
         } else {
            log.Println("Not ok")
         }
      } else {
         log.Println("Not ok")
      }
   }
}

func isDataUpdated(dir string, sym string, listDate string) (updated bool, start string, end string, fp *os.File) {
   fp = nil
   updated = false
   hisFile := fmt.Sprintf("%s.csv", dir+sym)
   _, ts := utils.GetTodayString(utils.DateFormat1)
   start = listDate
   end = ts
   isNew, fp := utils.OpenCsvFile(hisFile)
   log.Printf("[%s] isNew:%v csvFp: %v", hisFile, isNew, fp)
   if fp != nil {
      if isNew {
         utils.WriteData2CsvFile(fp, models.TradeDailyFieldSymbol)
         return updated, start, end, fp
      } else {
         csvr := csv.NewReader(fp)
         records, err := csvr.ReadAll()
         if err == nil {
            size := len(records)
            if size <=1 {
               log.Println("records is empty")
            } else {
               log.Printf("last line[%v]", records[size -1])
               lastUpdate := records[size - 1][1]
               diff := utils.CalcDateDiffByDay(utils.DateFormat1, lastUpdate, ts)
               log.Printf("diff = %d", diff)
               if diff > 1 { //at least is 2
                  start_new := utils.AddDays2Date(utils.DateFormat1, lastUpdate, 0, 0, 1)
                  start = start_new
               } else {
                  fp.Close()
                  fp = nil
                  updated = true
               }
            }
         }
      }
   } else {
      log.Printf("open file [%s] faild\n", hisFile)
   }

   return updated, start, end, fp
}

func TsInit(ctx *context.Context, cfg *config.Config, home string) *TsContext {
   log.Println(" ==========  TsInit Start ==========")
   ts_evts := make(map[string]*chan models.TsEvent)
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
      tse := make(chan models.TsEvent, 100)
      //stg := strategies.NewStockStrategy(cfg)
      collector.NewTsCollector(ctx, cl, tse, &sto)
      ts_evts[v.ExchangeName] = &tse
      ts_clis[v.ExchangeName] = cl
      ts_stos[v.ExchangeName] = &sto
   }

   var tsc TsContext
   tsc.TsEvents = &ts_evts
   tsc.TsClients = &ts_clis
   tsc.TsStorage = &ts_stos
   tsc.TsCron = make(chan models.DataFlag, 100)

   strategies.NewStockStrategy(cfg, tsc.TsCron)
   tsc.CheckUpdateBasics()
   tsc.CheckUpdateTradeCal()
   tsc.CheckUpdateNameChangeHistory()

   tsc.CheckUpdateDaily()
   tsc.TsCronSchedule(ctx)

   log.Println(" ==========  TsInit End ==========")
   return &tsc
}

func (tsc *TsContext) TsCronSchedule(ctx *context.Context) {
   if tsc == nil {
      log.Println("FATAL: tsc is nil!")
      return
   }
   go func() {
      for {
         select {
         case <-(*ctx).Done():
            log.Println("TsCronSchedule exit\n")
            return
         case o := <-tsc.TsCron:
            log.Printf("TsCronSchedule job [%d]\n", o)
            switch o {
            case models.DataFlag_Stock_Basic:
            case models.DataFlag_Stock_Company:
               tsc.CheckUpdateBasics()
            case models.DataFlag_Trade_Daily:
               tsc.CheckUpdateDaily()
            case models.DataFlag_Trade_Calendar:
               tsc.CheckUpdateTradeCal()
            case models.DataFlag_NameChange_Histtory:
               tsc.CheckUpdateNameChangeHistory()
            }
         }
      }
   }()
}
