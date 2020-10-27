package storage

//import "github.com/520lly/qt_data_service/models"

type Storage interface {
	SaveStockBasic(items *[][]interface{})
	SaveCompanyBasic(items *[][]interface{})
	GetFullPath() string
	GetStockBasic() string
	GetCompanyBasic() string
	SaveWorker()
	Close()
}
