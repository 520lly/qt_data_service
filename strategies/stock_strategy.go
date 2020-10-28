package strategies

import (
	"log"
	"regexp"
	"time"

	"github.com/520lly/qt_data_service/config"
	"github.com/520lly/qt_data_service/models"
	"github.com/520lly/qt_data_service/utils"
)

const (
   EMPTY_FILE_SIZE_MIN int64 = 1000
)

type TsEvent struct {
	DataFlag  models.DataFlag    `default: models.DataFlag_Stock_Basic`
	ApiParams *map[string]string `default: nil`
}
type StockStrategy struct {
	PeriodStockBasic   int64
	PeriodCompanyBasic int64
	PeriodDaily        int64
	TsEvent            chan TsEvent
}

func NewStockStrategy(cfg *config.Config) (stg *StockStrategy) {
	if cfg == nil {
		log.Println("cfg is nil")
		return
	}
	tse := make(chan TsEvent)
	//tse.DataFlag               = models.DataFlag_Stock_Basic
	//tse.ApiParams              = nil

	return &StockStrategy{
		PeriodStockBasic:   7,
		PeriodCompanyBasic: 7,
		PeriodDaily:        120,
		TsEvent:            tse,
	}
}

func CheckUpdateBasic(p string, market string, path string) (bool, []string) {
   var update bool = false
   err, files := utils.FilteredSearchOfDirectoryTree(regexp.MustCompile(market), path)
   if err == nil {
      //StockBasic is not exist, need to create it first
      if len(files) <= 0 {
         update = true
      } else { //Means StockBasic file exist, need to Check whether update need or not
         re := regexp.MustCompile("^.*" + market + "_(.*)\\.csv")
         log.Println("re = ", re)
         for _, f := range files {
            match := re.FindStringSubmatch(f)
            log.Println("match[1] = ", match[1])
            if len(match) > 1 {
               today, ts := utils.GetTodayString("2006-01-02")
               if ts == match[1] { //the date when the existed file was created was in today, it was just a iniitial file, need to be update
                  log.Println("ts", ts, match[1])
                  size, err := utils.GetFileSize(match[0])
                  log.Printf("err: %v size:%d", err, size)
                  if err == nil {
                     if size < EMPTY_FILE_SIZE_MIN {
                        update = true
                     }
                  }
               } else {
                  duration, _ := time.ParseDuration(p)
                  created, _ := time.Parse("2006-01-02", match[1])
                  expired := created.Add(duration)
                  if today.After(expired) {
                     log.Printf("created date %v expired date %v is reached", created, expired)
                     update = true
                  } else {
                     log.Printf("created date %v expired date %v is NOT reached", created, expired)
                  }
               }
            }
            break
         }
      }
   }
   return update, files
}


func (sstg *StockStrategy) CheckUpdateDailySingle(code string, start string, end string) {

}
