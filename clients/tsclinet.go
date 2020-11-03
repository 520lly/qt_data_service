package clients

import (
	"log"

	"github.com/520lly/qt_data_service/models"
	"github.com/520lly/qt_data_service/utils"
	tsg "github.com/520lly/tushare-go"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type TsClient struct {
	Market       string
	ExchangeName string
	CurrencyPair string
	c            *tsg.TuShare
}

func NewTsClient(market, exchangeName, currencyPair, token string) *TsClient {
	c := tsg.New(token)
	return &TsClient{
		Market:       market,
		ExchangeName: exchangeName,
		CurrencyPair: currencyPair,
		c:            c}
}

func (c *TsClient) GetTsStockBasics(params *map[string]string, fields *[]string) (items *[][]interface{}) {
	if fields == nil {
		log.Printf("ERROR: fields must be valid!")
		return
	}
   //TODO: How to ensure data to be collected if stock items are larger than MAX(so far 5000)
	data, err := c.c.StockBasic(*params, *fields)
	if err == nil {
		log.Printf("%T %s %d msg:%v fields:%v size:%d\n", data, data.RequestID, data.Code, data.Msg, data.Data.Fields, len(data.Data.Items))
		return &data.Data.Items
	}
	return nil
}

func (c *TsClient) GetStockBasic(params *map[string]string) (items *[][]interface{}) {
	data, err := c.c.StockBasic(*params, models.StockFieldSymbol)
	if err == nil {
		log.Printf("%T %s %d msg:%v fields:%v size:%d\n", data, data.RequestID, data.Code, data.Msg, data.Data.Fields, len(data.Data.Items))
		return &data.Data.Items
	}
	return nil
}

func (c *TsClient) GetCompanyBasic(params *map[string]string) (items *[][]interface{}) {
	data, err := c.c.StockCompany(*params, models.CompanyFieldSymbol)
	if err == nil {
		log.Printf("%T %s %d msg:%v fields:%v size:%d\n", data, data.RequestID, data.Code, data.Msg, data.Data.Fields, len(data.Data.Items))
		return &data.Data.Items
	}
	return nil
}

func (c *TsClient) GetTradeDaily(params *map[string]string, fields *[]string) (items *[][]interface{}) {
	if params == nil || fields == nil {
		log.Printf("ERROR: params and fields must be valid!")
		return
	}
   var data [][]interface{}
	//data = append(data, models.TradeDailyFieldSymbol)
   start_date := (*params)["start_date"]
   end_date := (*params)["end_date"]
   end_last := end_date
   start_last := start_date

	diff := utils.CalcDateDiffByDay(utils.DateFormat1, start_date, end_date)
   rounds := utils.CalcRoundsCeil(diff)
   var i int = 1
   log.Printf("end:%s - start:%s = %d, needs: %d rounds\n", end_date, start_date, diff, rounds)
   for ; i <= rounds; i++ {
      end_tmp := utils.AddDays2Date(utils.DateFormat1, start_last, 0, 0, i * utils.MAX_ITEMS_DAILY)
      if utils.IsDateAfter(utils.DateFormat1, end_tmp, end_date) {
         end_last = end_date
      } else {
         end_last = end_tmp
      }
      (*params)["end_date"] = end_last

      if utils.IsDateAfter(utils.DateFormat1, end_date, start_last) {
         (*params)["start_date"] = start_last
      }
      log.Printf("round[%d] -- start_date:%s - end_date:%s\n", i, (*params)["start_date"], (*params)["end_date"])
      res, err := c.c.Daily(*params, *fields)
      if err == nil {
         log.Printf("%T-id:%s code:%d msg:%v\n", res, res.RequestID, res.Code, res.Msg)
         //means success
         if res.Code == 0 {
            size := len(res.Data.Items)
            log.Printf("fields:%v size:%d\n", res.Data.Fields,size)
            if size > 0 {
               end_real := utils.ConvertInterface2String(res.Data.Items[0][1])
               start_real := utils.ConvertInterface2String(res.Data.Items[size - 1][1])
               start_last = utils.AddDays2Date(utils.DateFormat1, utils.ConvertInterface2String(res.Data.Items[0][1]), 0, 0, 1)
               utils.Reverse2DArray(res.Data.Items)
               data = append(data, res.Data.Items...)
               log.Printf("start_real:[%s] end_real:[%s] --> new start_last :[%s]", start_real, end_real, start_last)
            } else {
               log.Printf("FATAL: empty data")
            }
         } else {
            log.Printf("c.c.Daily respond failed!!!")
         }
      }
   }
   return &data
}

func (c *TsClient) GetTradeDailyReverse(params *map[string]string) (items *[][]interface{}) {
	if params == nil {
		log.Printf("ERROR: params must be valid!")
		return
	}
   var data [][]interface{}
	//data = append(data, models.TradeDailyFieldSymbol)
   start_date := (*params)["start_date"]
   end_date := (*params)["end_date"]
   end_last := end_date
   start_last := start_date

	diff := utils.CalcDateDiffByDay(utils.DateFormat1, start_date, end_date)
   rounds := utils.CalcRoundsCeil(diff)
   var i int = 1
   log.Printf("end:%s - start:%s = %d, needs: %d rounds\n", end_date, start_date, diff, rounds)
   for ; i <= rounds; i++ {
      //start_tmp := utils.SubDays2Date(end_last, i * utils.MAX_ITEMS_DAILY)
      start_tmp := utils.AddDays2Date(utils.DateFormat1, end_last, 0, 0, -i * utils.MAX_ITEMS_DAILY)
      if utils.IsDateAfter(utils.DateFormat1, start_tmp, start_date) {
         start_last = start_tmp
      } else {
         start_last = start_date
      }
      (*params)["start_date"] = start_last
      (*params)["end_date"] = end_last

      log.Printf("round[%d] -- start_date:%s - end_date:%s\n", i, (*params)["start_date"], (*params)["end_date"])
      res, err := c.c.Daily(*params, models.TradeDailyFieldSymbol)
      if err == nil {
         log.Printf("%T %s %d msg:%v fields:%v size:%d\n", res, res.RequestID, res.Code, res.Msg, res.Data.Fields, len(res.Data.Items))
         //means success
         if res.Code == 0 {
            utils.Reverse2DArray(data)
            data = append(data, res.Data.Items...)
            end_last = utils.ConvertInterface2String(res.Data.Items[len(res.Data.Items) -1 ][1])
            log.Printf("new end_last :[%s]", end_last)
         } else {
            log.Printf("c.c.Daily respond failed!!!")
         }
      }
   }
   return &data
}

func (c *TsClient) GetTradeCalender(params *map[string]string) (items *[][]interface{}) {
   if params == nil {
		log.Printf("ERROR: params must be valid!")
		return
   }
   return nil
}

