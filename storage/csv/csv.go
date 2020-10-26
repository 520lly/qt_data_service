package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/goex-top/market_center"
	jsoniter "github.com/json-iterator/go"
   "github.com/520lly/qt_data_service/utils"
   "github.com/520lly/qt_data_service/models"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type CsvStorage struct {
   market         string
	exchangeName   string
	pair           string
	contractType   string
	flag           market_center.DataFlag
	prefix         string
	outputPath     string
   fullPath       string
   saveStockBasic chan models.StockBasic 
	fileTimestamp  time.Time
	ctx            context.Context
	basicFile      *os.File
	basicCsv       *csv.Writer
}

func NewCsvStorage(
	ctx context.Context,
   market         string,
	exchangeName string,
	pair string,
	contractType string,
	flag market_center.DataFlag,
	prefix string,
	outputPath string,
) *CsvStorage {
	var saveStockBasic chan models.StockBasic
	var basicFile *os.File
	var basicCsv *csv.Writer

	fileTimestamp := time.Now()
	ts := fileTimestamp.Format("2006-01-02")
	isNew := false

   path := prefix + "/" + outputPath + "/" + exchangeName + "/" + market
   ret, _ := utils.EnsurePathExist(path)
   if ret == false {
      panic("path was not exist!")
   }
   log.Println("%s exist!", path)

   isNew, basicFile = OpenCsvFile(fmt.Sprintf("%s/basic_%s.csv", path,ts))
   basicCsv = csv.NewWriter(basicFile)
   if isNew {
      data := models.FieldSymbol
      basicCsv.Write(data)
      basicCsv.Flush()
   }
   saveStockBasic = make(chan models.StockBasic)

	return &CsvStorage{
		ctx:            ctx,
      market:         market,
		exchangeName:   exchangeName,
		pair:           pair,
		contractType:   contractType,
		flag:           flag,
		prefix:         prefix,
		outputPath:     outputPath,
      fullPath:       path,
		fileTimestamp:  fileTimestamp,
      saveStockBasic: saveStockBasic,
		basicFile:      basicFile,
		basicCsv:       basicCsv,
	}
}

//func (s *CsvStorage) SaveStockBasic(sb *models.StockBasic) {
func (s *CsvStorage) SaveStockBasic(sb *[]string) {
	if s.saveStockBasic == nil {
		return
	}
	s.saveStockBasic <- *sb
}

func (s *CsvStorage) Close() {
	if s.basicCsv != nil {
		s.basicCsv.Flush()
		s.basicFile.Close()
	}
	//close(s.saveDepthChan)
	//close(s.saveTickerChan)
	//close(s.saveKlineChan)
}

func (s *CsvStorage) compress(fileTimestamp time.Time) {
	ts := fileTimestamp.Format("2006-01-02")
	//src := fmt.Sprintf("%s_%s_%s.csv", s.exchangeName, s.pair, ts)
	filters := []string{s.exchangeName, s.pair, ts, ".csv"}
	dst := fmt.Sprintf("%s/%s_%s_%s.tar.gz", s.outputPath, s.exchangeName, s.pair, ts)

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

func (s *CsvStorage) reNewFile() {
	now := time.Now()
	if now.Day() == s.fileTimestamp.Day() {
		return
	}
	s.Close()
	log.Printf("now day:%d, file timestamp day:%d", now.Day(), s.fileTimestamp.Day())
	go s.compress(s.fileTimestamp)

	s.fileTimestamp = now

	ts := s.fileTimestamp.Format("2006-01-02")
   isNew := false

   isNew, s.basicFile = OpenCsvFile(fmt.Sprintf("%s/basic_%s.csv", s.fullPath,ts))
   s.basicCsv = csv.NewWriter(s.basicFile)
   if isNew {
      data := models.FieldSymbol
      s.basicCsv.Write(data)
      s.basicCsv.Flush()
   }
}

func (s *CsvStorage) SaveWorker() {

	tick := time.NewTicker(time.Second)

	for {
		select {
		case <-tick.C:
			s.reNewFile()
		case o := <-s.saveStockBasic:
			data := []string{
				fmt.Sprint(o.TsCode),
				fmt.Sprint(o.Symbol),
				fmt.Sprint(o.Name ),
				fmt.Sprint(o.Area),
				fmt.Sprint(o.Industry),
				fmt.Sprint(o.Market),
				fmt.Sprint(o.Exchange),
				fmt.Sprint(o.CurrType),
				fmt.Sprint(o.ListStatus),
				fmt.Sprint(o.ListDate),
				fmt.Sprint(o.DelistDate),
				fmt.Sprint(o.IsHs)}
			s.basicCsv.Write(data)
			s.basicCsv.Flush()

		case <-s.ctx.Done():
			s.Close()
			log.Printf("(%s) %s saveWorker exit\n", s.exchangeName, s.pair)
			return
		}
	}
}

func OpenCsvFile(fileName string) (bool, *os.File) {
	var file *os.File
	var err1 error
	var isNew = false
	checkFileIsExist := func(fileName string) bool {
		var exist = true
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			exist = false
		}
		return exist
	}
	if checkFileIsExist(fileName) {
		file, err1 = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 666)
	} else {
		file, err1 = os.Create(fileName)
		isNew = true
	}
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "unable to write file on filehook %v", err1)
		panic(err1)
	}
	return isNew, file
}
