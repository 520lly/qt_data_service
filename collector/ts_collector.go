package collector

import (
	"log"
   "time"
	"context"

	"github.com/520lly/qt_data_service/clients"
	"github.com/520lly/qt_data_service/models"
   "github.com/520lly/qt_data_service/storage"
   mycsv"github.com/520lly/qt_data_service/storage/csv"
)

func NewTsCollector(ctx *context.Context, c *clients.TsClient, tse chan models.TsEvent, store *storage.Storage) {
   log.Printf("(%s) %s new collector with tse[%v]\n", c.ExchangeName, c.CurrencyPair, tse)
   if ctx == nil || c == nil || tse == nil || store == nil {
      panic("ctx c stg store are(is) nil!!!")
      return
   }
	tick := time.NewTicker(time.Second)
   go func() {
      sci := &mycsv.StoCsvInstance{nil, nil, nil}
      for {
         select {
         case <-tick.C:
            log.Printf("Collector [%s] is in IDLE!", c.ExchangeName)
         case <-(*ctx).Done():
            log.Printf("(%s) %s collector exit\n", c.Market, c.ExchangeName)
            return
         case o := <-tse:
            switch o.DataFlag {
            case models.DataFlag_Stock_Basic:
               if o.ApiParams != nil {
                  log.Printf("ApiParams: %v\n", o.ApiParams)
                  data := c.GetStockBasic(o.ApiParams)
                  if data != nil {
                     sci.DoubleArrayData = data
                     sci.CsvFile = o.CsvFile
                     (*store).SaveData(sci)
                  }
               } else {
                  log.Println("FATAL: receive nil csv file writer!")
               }
            case models.DataFlag_Stock_Company:
               if o.ApiParams != nil {
                  log.Printf("ApiParams: %v\n", o.ApiParams)
                  data := c.GetCompanyBasic(o.ApiParams)
                  if data != nil {
                     sci.DoubleArrayData = data
                     sci.CsvFile = o.CsvFile
                     (*store).SaveData(sci)
                  }
               } else {
                  log.Println("FATAL: receive nil csv file writer!")
               }
            case models.DataFlag_Trade_Daily:
               if o.CsvFile != nil {
                  log.Printf("ApiParams: %v\n", o.ApiParams)
                  data := c.GetTradeDaily(o.ApiParams, &models.TradeDailyFieldSymbol)
                  if data != nil {
                     sci.DoubleArrayData = data
                     sci.CsvFile = o.CsvFile
                     (*store).SaveData(sci)
                  }
               } else {
                  log.Println("FATAL: receive nil csv file writer!")
               }
            case models.DataFlag_Trade_Calendar:
               if o.CsvFile != nil {
                  log.Printf("ApiParams: %v\n", o.ApiParams)
                  data := c.GetTradeCalender(o.ApiParams)
                  if data != nil {
                     sci.DoubleArrayData = data
                     sci.CsvFile = o.CsvFile
                     (*store).SaveData(sci)
                  }
               } else {
                  log.Println("FATAL: receive nil csv file writer!")
               }
            case models.DataFlag_NameChange_Histtory:
               if o.CsvFile != nil {
                  log.Printf("ApiParams: %v\n", o.ApiParams)
                  data := c.GetNameChangeHistory(o.ApiParams)
                  if data != nil {
                     sci.DoubleArrayData = data
                     sci.CsvFile = o.CsvFile
                     (*store).SaveData(sci)
                  }
               } else {
                  log.Println("FATAL: receive nil csv file writer!")
               }
            default:
            }
         }
      }
   }()
}
