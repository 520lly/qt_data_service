package strategies

import (
	"log"

	"github.com/520lly/qt_data_service/config"
	"github.com/520lly/qt_data_service/models"
)

type TsEvent struct {
	DataFlag  models.DataFlag    `default: models.DataFlag_Stock_Basic`
	ApiParams *map[string]string `default: nil`
}
type StockStrategy struct {
	PeriodStockBasic   int64
	PeriodCompanyBasic int64
	PeriodDaily        int64
	TsEvent            chan TsEvent
}

func NewStockStrategy(cfg *config.Config) (stg *StockStrategy) {
	if cfg == nil {
		log.Println("cfg is nil")
		return
	}
	tse := make(chan TsEvent)
	//tse.DataFlag               = models.DataFlag_Stock_Basic
	//tse.ApiParams              = nil

	return &StockStrategy{
		PeriodStockBasic:   7,
		PeriodCompanyBasic: 7,
		PeriodDaily:        120,
		TsEvent:            tse,
	}
}

func (sstg *StockStrategy) InitBasic() {

}
