package strategies

import (
   "log"

   "github.com/520lly/qt_data_service/config"
   "github.com/520lly/qt_data_service/models"
   "github.com/jasonlvhit/gocron"
)

const (
   EMPTY_FILE_SIZE_MIN int = 1000
)

var (
   TsEvents chan models.DataFlag
)

type StockStrategy struct {
   PeriodStockBasic   int
   PeriodCompanyBasic int
   PeriodDaily        int
}

func tsTaskEvent(e models.DataFlag) {
   log.Printf("tsTaskEvent %d", e)
   TsEvents<-e
}

func NewStockStrategy(cfg *config.Config, tse chan models.DataFlag) (stg *StockStrategy) {
   if cfg == nil {
      log.Println("cfg is nil")
      return
   }
   log.Printf("TODO finish cfg:%v", *cfg)
   TsEvents = tse

   sts := &StockStrategy{
      PeriodStockBasic:   7,
      PeriodCompanyBasic: 7,
      PeriodDaily:        120,
   }
   s := gocron.NewScheduler()
   s.Every(7).Days().Do(tsTaskEvent, models.DataFlag_Stock_Basic)
   s.Every(1).Day().At("15:00").Do(tsTaskEvent, models.DataFlag_Trace_Daily)

   // Start all the pending jobs
   log.Println("+++++++++++++++++++ gocron.Start ")
   go func () {
      <- s.Start()
   }()
   return sts
}

