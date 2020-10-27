package strategies

import (
   "github.com/520lly/qt_data_service/models"
   "github.com/520lly/qt_data_service/config"
)

type TsEvent struct {
   DataFlag                models.DataFlag `default: models.DataFlag_Stock_Basic`
   ApiParams               *map[string]string `default: nil`
}
type StockStrategy struct {
   PeriodStockBasic        int64
   PeriodCompanyBasic      int64
   PeriodDaily             int64
   TsEvent                 chan TsEvent
}

func (sstg *StockStrategy) LoadStockStrntegy(cfg *config.Config) (stg *StockStrategy){
   //sstg.PeriodStockBasic = 7
   //sstg.PeriodCompanyBasic = 7
   //sstg.PeriodDaily = 120
   tse := make(chan TsEvent)
   //tse.DataFlag               = models.DataFlag_Stock_Basic
   //tse.ApiParams              = nil

   return &StockStrategy{
      PeriodStockBasic:             7,
      PeriodCompanyBasic:           7,
      PeriodDaily:                  120,
      TsEvent:                      tse,
   }
}

func (sstg *StockStrategy) InitBasic() {
   
}
