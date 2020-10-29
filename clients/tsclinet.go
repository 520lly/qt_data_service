package clients

import (
	"fmt"

	"github.com/520lly/qt_data_service/models"
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

func (c *TsClient) GetStockBasic() (items *[][]interface{}) {
	params := make(map[string]string)
	params["is_hs"] = "N"
	params["list_status"] = "L"
	params["exchange"] = ""
	data, err := c.c.StockBasic(params, models.FieldSymbol)
	if err == nil {
		fmt.Printf("%T %s %d msg:%v fields:%v size:%d\n", data, data.RequestID, data.Code, data.Msg, data.Data.Fields, len(data.Data.Items))
		return &data.Data.Items
	}
	return nil
}

func (c *TsClient) GetCompanyBasic() (items *[][]interface{}) {
	params := make(map[string]string)
	data, err := c.c.StockCompany(params, models.CompanyFieldSymbol)
	if err == nil {
		fmt.Printf("%T %s %d msg:%v fields:%v size:%d\n", data, data.RequestID, data.Code, data.Msg, data.Data.Fields, len(data.Data.Items))
		return &data.Data.Items
	}
	return nil
}

func (c *TsClient) GetTradeDaily(params *map[string]string) (items *[][]interface{}) {
	if params == nil {
		fmt.Printf("ERROR: params must be valid!")
		return
	}
	var data [][]string
	data = append(data, models.TradeDailyFieldSymbol)
	fmt.Printf("start:%s end:%s\n", (*params)["start_date"], (*params)["end_date"])
	res, err := c.c.Daily(*params, models.TradeDailyFieldSymbol)
	if err == nil {
		fmt.Printf("%T %s %d msg:%v fields:%v size:%d\n", res, res.RequestID, res.Code, res.Msg, res.Data.Fields, len(res.Data.Items))
		//means success
		if res.Code == 0 {

		}
		return &res.Data.Items
	}
	return nil
}
