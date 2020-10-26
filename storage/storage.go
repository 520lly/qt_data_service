package storage

//import "github.com/520lly/qt_data_service/models"

type Storage interface {
	SaveStockBasic(sb *[]string)
	SaveWorker()
	Close()
}
