package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

type CsvStorage struct {
	//market           string
	//exchangeName     string
	//pair             string
	//contractType     string
	sci    chan StoCsvInstance
	sub    config.Subscribe
	store  config.Storage
	prefix string
	//outputPath       string
	fullPath         string
	saveStockBasic   chan [][]interface{}
	saveCompanyBasic chan [][]interface{}
	fileTimestamp    time.Time
	ctx              context.Context
	basicFile        *os.File
	basicCsv         *csv.Writer
	companyFile      *os.File
	companyCsv       *csv.Writer
}

func NewCsvStorage(
	ctx context.Context,
	store config.Storage,
	//market string,
	//exchangeName string,
	//pair string,
	//contractType string,
	sub config.Subscribe,
	prefix string,
	//outputPath string,
) *CsvStorage {
	var saveStockBasic chan [][]interface{}
	var saveCompanyBasic chan [][]interface{}
	var basicFile *os.File
	var basicCsv *csv.Writer
	var companyFile *os.File
	var companyCsv *csv.Writer

	fileTimestamp := time.Now()
	ts := fileTimestamp.Format("2006-01-02")
	isNew := false

	root := prefix + "/" + store.CsvCfg.Location + "/" + sub.Market + "/" + sub.ExchangeName
	ret, _ := utils.EnsurePathExist(root)
	if ret == false {
		panic("path was not exist!")
	}
	log.Printf("%s exist!", root)

	history := root + "/" + store.CsvCfg.History
	ret, _ = utils.EnsurePathExist(history)
	if ret == false {
		panic("path was not exist!")
	}
	log.Printf("%s exist!", history)

	isNew, basicFile = utils.OpenCsvFile(fmt.Sprintf("%s/%s_%s.csv", root, sub.StockBasic, ts))
	basicCsv = csv.NewWriter(basicFile)
	if isNew {
		utils.WriteData2CsvFile(basicCsv, models.FieldSymbol)
	}
	saveStockBasic = make(chan [][]interface{})

	isNew, companyFile = utils.OpenCsvFile(fmt.Sprintf("%s/%s_%s.csv", root, sub.CompanyBasic, ts))
	companyCsv = csv.NewWriter(companyFile)
	if isNew {
		utils.WriteData2CsvFile(companyCsv, models.CompanyFieldSymbol)
	}
	saveCompanyBasic = make(chan [][]interface{})

	sci := make(chan StoCsvInstance)

	return &CsvStorage{
		ctx: ctx,
		//market:           sub.Market,
		//exchangeName:     sub.ExchangeName,
		//pair:             sub.CurrencyPair,
		//contractType:     sub.ContractType,
		sci:    sci,
		sub:    sub,
		store:  store,
		prefix: prefix,
		//outputPath:       outputPath,
		fullPath:         root,
		fileTimestamp:    fileTimestamp,
		saveStockBasic:   saveStockBasic,
		basicFile:        basicFile,
		basicCsv:         basicCsv,
		saveCompanyBasic: saveCompanyBasic,
		companyFile:      companyFile,
		companyCsv:       companyCsv,
	}
}

func (s *CsvStorage) SaveStockBasic(items *[][]interface{}) {
	if s.saveStockBasic == nil {
		return
	}
	s.saveStockBasic <- *items
}

func (s *CsvStorage) SaveCompanyBasic(items *[][]interface{}) {
	if s.saveCompanyBasic == nil {
		return
	}
	s.saveCompanyBasic <- *items
}

func (s *CsvStorage) SaveData(sci *StoCsvInstance) {
	if sci == nil {
		return
	}
	s.sci <- *sci
	log.Printf("save data [%v]\n", sci)
}

//func (s *CsvStorage) GetStoCsvInstance() *StoCsvInstance {
//return s.sci
//}

func (s *CsvStorage) UpdateStoCsvInstance(da_data *[][]interface{}, sa_data *[]interface{}, csvw *os.File) {
	//s.sci.DoubleArrayData = da_data
	//s.sci.SingleArrayData = sa_data
	//s.sci.CsvFile =  csvw

}

func (s *CsvStorage) GetFullPath() string {
	return s.fullPath
}
func (s *CsvStorage) GetSubscribe() config.Subscribe {
	return s.sub
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

	isNew, s.basicFile = utils.OpenCsvFile(fmt.Sprintf("%s/basic_%s.csv", s.fullPath, ts))
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
			//s.reNewFile()
			log.Println("tick.C")
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
		case o := <-s.saveStockBasic:
			//empty old file
			//WriteData2CsvFile(companyCsv, models.CompanyFieldSymbol)
			for _, item := range o {
				ss := strings.Split(fmt.Sprintf("%s", item[0]), ".")
				if strings.ToUpper(s.sub.ExchangeName) == ss[1] {
					var data []string
					for _, f := range item {
						data = append(data, fmt.Sprintf("%v", f))
					}
					s.basicCsv.Write(data)
				}
			}
			s.basicCsv.Flush()

		case o := <-s.saveCompanyBasic:
			//empty old file
			for _, item := range o {
				ss := strings.Split(fmt.Sprintf("%s", item[0]), ".")
				if strings.ToUpper(s.sub.ExchangeName) == ss[1] {
					var data []string
					for _, f := range item {
						data = append(data, fmt.Sprintf("%v", f))
					}
					s.companyCsv.Write(data)
				}
			}
			s.companyCsv.Flush()

		case <-s.ctx.Done():
			s.Close()
			log.Printf("(%s) %s saveWorker exit\n", s.sub.ExchangeName, s.sub.CurrencyPair)
			return
		}
	}
}
