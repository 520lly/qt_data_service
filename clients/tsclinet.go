package clients

import(
   "fmt"

   //"github.com/520lly/qt_data_service/models"
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

func (c *TsClient) GetStockBasic() (sb *[]string){
   params := make(map[string]string)
   params["is_hs"]="N"
   params["list_status"]="L"
   params["exchange"]=""
   fields := []string{ "ts_code", "symbol", "name", "area", "industry", "fullname", "market", "exchange", "curr_type", "list_status", "list_date", "delist_date", "is_hs"}
   data, err := c.c.StockBasic(params, fields)
   if err == nil {
      //var infos []*models.StockBasic
      fmt.Printf("%T %s %d msg:%v fields:%v size:%d\n", data, data.RequestID, data.Code, data.Msg, data.Data.Fields, len(data.Data.Items))
      for _, item := range data.Data.Items {
         //fmt.Printf("%T %v\n", item, item)
         str := fmt.Sprintf("%s", item)
         fmt.Printf("%T %v\n", str, str)
         //var data []string
         //for _, f := range item {
            //fmt.Printf("Type %T %v\n", f, f)
            ////data = append(data, f)
         //}
      }
   }
   return nil
}

