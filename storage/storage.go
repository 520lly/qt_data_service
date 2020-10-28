package storage

import "github.com/520lly/qt_data_service/config"

type Storage interface {
	SaveStockBasic(items *[][]interface{})
	SaveCompanyBasic(items *[][]interface{})
	GetFullPath() string
   GetSubscribe() config.Subscribe
	SaveWorker()
	Close()
}
