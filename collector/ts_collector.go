package collector

import(
   "log"
   "time"
   "context"

   "github.com/520lly/qt_data_service/clients"
   "github.com/520lly/qt_data_service/storage"
)


func NewTsCollector(ctx context.Context, c *clients.TsClient, period int64, store storage.Storage) {
	log.Printf("(%s) %s new collector with period[%d]\n", c.ExchangeName, c.CurrencyPair, period)
   go func() {
      tick := time.NewTicker(time.Millisecond * time.Duration(period))
      for {
         select {
         case <-ctx.Done():
            log.Printf("(%s) %s collector exit\n", c.Market, c.ExchangeName)
            return
         case <-tick.C:
            log.Printf("(%v) collector tick.C\n", tick.C)
            data := c.GetStockBasic()
            if data != nil {
               store.SaveStockBasic(data)
            }
         }
      }
   }()
}
