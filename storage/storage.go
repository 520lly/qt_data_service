package storage

import (
   "github.com/520lly/qt_data_service/config"
   mycsv "github.com/520lly/qt_data_service/storage/csv"
)

type Storage interface {
	SaveStockBasic(items *[][]interface{})
	SaveCompanyBasic(items *[][]interface{})
	GetFullPath() string
   GetSubscribe() config.Subscribe
	SaveWorker()
	SaveData(sci *mycsv.StoCsvInstance)
	Close()
}
