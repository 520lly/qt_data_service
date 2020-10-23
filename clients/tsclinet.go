package clients

import(
	//"fmt"

  tsg "github.com/520lly/tushare-go"
)

type TsClient struct {
	Market       string
	ExchangeName string
	CurrencyPair string
   c            *tsg.TuShare
}

func NewTsClient(market, exchangeName, currencyPair, token string) *TsClient {
   c := tsg.New(token)
   //params := make(map[string]string)
   //params["is_hs"]="N"
   //params["list_status"]="L"
   //params["exchange"]=""
   //fields := []string{ "ts_code", "symbol", "name", "area", "industry", "fullname", "market", "exchange", "curr_type", "list_status", "list_date", "delist_date", "is_hs"}
   //data, _ := c.StockBasic(params, fields)

   return &TsClient{
      Market:       market,
      ExchangeName: exchangeName,
		CurrencyPair: currencyPair,
		c:            c}
}

