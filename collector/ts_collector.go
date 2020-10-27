package collector

import (
	"log"
	//"time"
	"context"

	"github.com/520lly/qt_data_service/clients"
	"github.com/520lly/qt_data_service/models"
	"github.com/520lly/qt_data_service/storage"
	"github.com/520lly/qt_data_service/strategies"
	//"github.com/520lly/qt_data_service/utils"
)

func NewTsCollector(ctx *context.Context, c *clients.TsClient, stg *strategies.StockStrategy, store *storage.Storage) {
	log.Printf("(%s) %s new collector with stg[%v]\n", c.ExchangeName, c.CurrencyPair, stg)
	if ctx == nil || c == nil || stg == nil || store == nil {
		panic("ctx c stg store are(is) nil!!!")
		return
	}
	go func() {
		//tickStockBasic := time.NewTicker(time.Millisecond * time.Duration(period))
		//tickCompanyBasic := time.NewTicker(time.Millisecond * time.Duration(period))
		for {
			select {
			case <-(*ctx).Done():
				log.Printf("(%s) %s collector exit\n", c.Market, c.ExchangeName)
				return
			//case <-tickStockBasic.C:
			//defer tickStockBasic.Stop()
			//log.Printf("(%v) collector tickStockBasic.C\n", tickStockBasic.C)
			//data := c.GetStockBasic()
			//if data != nil {
			//store.SaveStockBasic(data)
			//}
			case o := <-stg.TsEvent:
				switch o.DataFlag {
				case models.DataFlag_Stock_Basic:
					data := c.GetStockBasic()
					if data != nil {
						(*store).SaveStockBasic(data)
					}
				case models.DataFlag_Stock_Company:
					data := c.GetCompanyBasic()
					if data != nil {
						(*store).SaveCompanyBasic(data)
					}
				case models.DataFlag_Trace_Daily:
					data := c.GetTradeDaily(o.ApiParams)
					if data != nil {
						(*store).SaveCompanyBasic(data)
					}
				default:
				}
			}
		}
	}()
}
