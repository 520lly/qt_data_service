package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
   //"strings"
	"time"
   "strconv"

	"github.com/520lly/qt_data_service/config"
	"github.com/520lly/qt_data_service/models"
	"github.com/520lly/qt_data_service/utils"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type StoCsvInstance struct {
	DoubleArrayData *[][]interface{}
	SingleArrayData *[]interface{}
	CsvFile         *os.File
}

type BasicCsvFileInfo struct {
   FullPath string
   LastModified string
   NextUpdate string
}

type CsvStorage struct {
	sci    chan StoCsvInstance
	sub    config.Subscribe
	store  config.Storage
	prefix string
	fullPath         string
	//saveStockBasic   chan [][]interface{}
	//saveCompanyBasic chan [][]interface{}
	//fileTimestamp    time.Time
   files            *map[string]BasicCsvFileInfo
   ctx              context.Context
	//basicFile        *os.File
	//basicCsv         *csv.Writer
	//companyFile      *os.File
	//companyCsv       *csv.Writer
}

func NewCsvStorage(
	ctx context.Context,
	store config.Storage,
	sub config.Subscribe,
	prefix string,
) *CsvStorage {

   root := prefix + "/" + store.CsvCfg.Location + "/" + sub.Market + "/" + sub.ExchangeName
   sci := make(chan StoCsvInstance)

   csvSto := &CsvStorage{
      ctx:              ctx,
      sci:              sci,
      sub:              sub,
      store:            store,
      prefix:           prefix,
      fullPath:         root,
      files:            nil,
   }
   csvSto.InitCsvWorkspace(root, sub)

   return csvSto
}

func (c *CsvStorage) InitCsvWorkspace(root string, sub config.Subscribe) {
   ret, _ := utils.EnsurePathExist(root)
   if ret == false {
      panic("path was not exist and create failed!")
   }
   log.Printf("%s exist!", root)

   history := root + "/history"
   ret, _ = utils.EnsurePathExist(history)
   if ret == false {
      panic("path was not exist and create failed!")
   }
   log.Printf("%s exist!", history)

   _, ts := utils.GetTodayString(utils.DateFormat2)
   var basicPat string = `basic_[a-z]{2,}.*\.csv`
   fileMap := make(map[string]BasicCsvFileInfo)
   err, files := utils.FilteredSearchOfDirectoryTree(basicPat, root) 
   if err == nil && len(files) >=1 {
      log.Printf("InitCsvWorkspace files:%v", files)
      for _, f := range files {
         // Check StockBasic
         if utils.IsFileSameFromFullPath(sub.StockBasic, f) {
            //log.Printf("------- pat:%s matchs f:%s", sub.StockBasic, f)
            updatedate := utils.GetLastModifyTime(f)
            csvInfo := BasicCsvFileInfo{f,updatedate,""}
            i, err := strconv.Atoi(sub.Period.StockBasic) 
            if err == nil {
               nextupdatedate := utils.AddDays2Date(utils.DateFormat2, updatedate, 0, 0 , i)
               csvInfo.NextUpdate = nextupdatedate
            }
            fileMap[sub.StockBasic] = csvInfo
         } else {
            isNew, fp := utils.OpenCsvFile(fmt.Sprintf("%s/%s.csv", root, sub.StockBasic))
            if isNew {
               utils.WriteData2CsvFile(fp, models.StockFieldSymbol)
               fileMap[sub.StockBasic] = BasicCsvFileInfo{fmt.Sprintf("%s/%s.csv", root, sub.StockBasic), ts, ts}
               fp.Close()
            }
         }

         // Check CompanyBasic
         if utils.IsFileSameFromFullPath(sub.CompanyBasic, f) {
            //log.Printf("------- pat:%s matchs f:%s", sub.CompanyBasic, f)
            updatedate := utils.GetLastModifyTime(f)
            csvInfo := BasicCsvFileInfo{f,updatedate,""}
            i, err := strconv.Atoi(sub.Period.CompanyBasic) 
            if err == nil {
               nextupdatedate := utils.AddDays2Date(utils.DateFormat2, updatedate, 0, 0 , i)
               csvInfo.NextUpdate = nextupdatedate
            }
            fileMap[sub.CompanyBasic] = csvInfo
         } else {
            //log.Printf("pat:%s not matched", sub.CompanyBasic)
            isNew, fp := utils.OpenCsvFile(fmt.Sprintf("%s/%s.csv", root, sub.CompanyBasic))
            if isNew {
               utils.WriteData2CsvFile(fp, models.CompanyFieldSymbol)
               fileMap[sub.CompanyBasic] = BasicCsvFileInfo{fmt.Sprintf("%s/%s.csv", root, sub.CompanyBasic), ts, ts}
               fp.Close()
            }
         }

         // Check TradeCal
         if utils.IsFileSameFromFullPath(sub.TradeCal, f) {
            //log.Printf("------- pat:%s matchs f:%s", sub.TradeCal, f)
            updatedate := utils.GetLastModifyTime(f)
            csvInfo := BasicCsvFileInfo{f,updatedate,""}
            i, err := strconv.Atoi(sub.Period.TradeDaily)  //Daily as default
            if err == nil {
               nextupdatedate := utils.AddDays2Date(utils.DateFormat2, updatedate, 0, 0 , i)
               csvInfo.NextUpdate = nextupdatedate
            }
            fileMap[sub.TradeCal] = csvInfo
         } else {
            //log.Printf("pat:%s not matched", sub.TradeCal)
            nfile := fmt.Sprintf("%s/%s.csv", root, sub.TradeCal)
            isNew, fp := utils.OpenCsvFile(nfile)
            if isNew {
               utils.WriteData2CsvFile(fp, models.TradeCalendarFieldSymbol)
               fileMap[sub.TradeCal] = BasicCsvFileInfo{nfile, ts, ts}
               fp.Close()
            }
         }

         // Check Name Change History
         if utils.IsFileSameFromFullPath(sub.NameChangeHistory, f) {
            //log.Printf("------- pat:%s matchs f:%s", sub.NameChangeHistory, f)
            updatedate := utils.GetLastModifyTime(f)
            csvInfo := BasicCsvFileInfo{f,updatedate,""}
            i, err := strconv.Atoi(sub.Period.StockBasic)  //Weekly as default same as Stock list
            if err == nil {
               nextupdatedate := utils.AddDays2Date(utils.DateFormat2, updatedate, 0, 0 , i)
               csvInfo.NextUpdate = nextupdatedate
            }
            fileMap[sub.NameChangeHistory] = csvInfo
         } else {
            nfile := fmt.Sprintf("%s/%s.csv", root, sub.NameChangeHistory)
            isNew, fp := utils.OpenCsvFile(nfile)
            if isNew {
               utils.WriteData2CsvFile(fp, models.NameChangeFieldSymbol)
               fileMap[sub.NameChangeHistory] = BasicCsvFileInfo{nfile, ts, ts}
               fp.Close()
            }
         }

         //More to be defined
      }
   } else {
      isNew, fp := utils.OpenCsvFile(fmt.Sprintf("%s/%s.csv", root, sub.StockBasic))
      if isNew {
         utils.WriteData2CsvFile(fp, models.StockFieldSymbol)
         fileMap[sub.StockBasic] = BasicCsvFileInfo{fmt.Sprintf("%s/%s.csv", root, sub.StockBasic), ts, ts}
         fp.Close()
         fp = nil
      }

      isNew, fp = utils.OpenCsvFile(fmt.Sprintf("%s/%s.csv", root, sub.CompanyBasic))
      if isNew {
         utils.WriteData2CsvFile(fp, models.CompanyFieldSymbol)
         fileMap[sub.CompanyBasic] = BasicCsvFileInfo{fmt.Sprintf("%s/%s.csv", root, sub.CompanyBasic), ts, ts}
         fp.Close()
         fp = nil
      }

      isNew, fp = utils.OpenCsvFile(fmt.Sprintf("%s/%s.csv", root, sub.TradeCal))
      if isNew {
         utils.WriteData2CsvFile(fp, models.TradeCalendarFieldSymbol)
         fileMap[sub.TradeCal] = BasicCsvFileInfo{fmt.Sprintf("%s/%s.csv", root, sub.TradeCal), ts, ts}
         fp.Close()
         fp = nil
      }

      isNew, fp = utils.OpenCsvFile(fmt.Sprintf("%s/%s.csv", root, sub.NameChangeHistory))
      if isNew {
         utils.WriteData2CsvFile(fp, models.NameChangeFieldSymbol)
         fileMap[sub.NameChangeHistory] = BasicCsvFileInfo{fmt.Sprintf("%s/%s.csv", root, sub.NameChangeHistory), ts, ts}
         fp.Close()
         fp = nil
      }
      //More to be defined

   }
   c.files = &fileMap
   log.Println(fileMap)
}

func (s *CsvStorage) SaveData(sci *StoCsvInstance) {
   if sci == nil {
      return
   }
   s.sci <- *sci
}

func (s *CsvStorage) GetFullPath() string {
	return s.fullPath
}

func (s *CsvStorage) GetSubscribe() config.Subscribe {
	return s.sub
}

func (s *CsvStorage) GetBasicFile(pat string) (BasicCsvFileInfo, bool) {
   val, ok := (*s.files)[pat]
   log.Printf("files:%v", (*s.files))
   log.Printf("pat:%s val:%v, ok:%t", pat, val, ok)
   return val, ok
}

func (s *CsvStorage) Close() {
   log.Println("empty")
	//if s.basicFile != nil {
		//s.basicFile.Close()
	//}
	//close(s.saveDepthChan)
	//close(s.saveTickerChan)
	//close(s.saveKlineChan)
}

func (s *CsvStorage) compress(fileTimestamp time.Time) {
	ts := fileTimestamp.Format("2006-01-02")
	//src := fmt.Sprintf("%s_%s_%s.csv", s.exchangeName, s.pair, ts)
	filters := []string{s.sub.ExchangeName, s.sub.CurrencyPair, ts, ".csv"}
	dst := fmt.Sprintf("%s/%s_%s_%s.tar.gz", s.store.CsvCfg.Location, s.sub.ExchangeName, s.sub.CurrencyPair, ts)

	csvs := GetSrcFileName(s.prefix, filters)
	log.Println("start to compress *.csv to *.tar.gz, ts:", ts)
	err := CompressFile(s.prefix, csvs, dst)
	if err != nil {
		log.Println(err)
		return
	}
	for _, v := range csvs {
		err := os.Remove(s.prefix + "/" + v)
		if err != nil {
			log.Printf("remove file %s fail:%s\n", s.prefix+"/"+v, err.Error())
		} else {
			log.Printf("remove file %s success\n", s.prefix+"/"+v)
		}
	}
}

func (s *CsvStorage) SaveWorker() {
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-tick.C:
         log.Printf("SaveWorker [%s] is in IDLE!", s.sub.ExchangeName)
		case o := <-s.sci:
			log.Printf("o.CsvFile: %v", o.CsvFile)
			if o.CsvFile != nil {
				csvw := csv.NewWriter(o.CsvFile)
				if o.DoubleArrayData != nil {
					for _, item := range *o.DoubleArrayData {
						var data []string
						for _, f := range item {
							data = append(data, fmt.Sprintf("%v", f))
						}
						csvw.Write(data)
					}
				}
				csvw.Flush()
				o.CsvFile.Close()
			}
		case <-s.ctx.Done():
			s.Close()
			log.Printf("(%s) %s saveWorker exit\n", s.sub.ExchangeName, s.sub.CurrencyPair)
			return
		}
	}
}
