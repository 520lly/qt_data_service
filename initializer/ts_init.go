package initializer

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/520lly/qt_data_service/clients"
	"github.com/520lly/qt_data_service/collector"
	"github.com/520lly/qt_data_service/config"
	"github.com/520lly/qt_data_service/models"
	"github.com/520lly/qt_data_service/storage"
	mycsv "github.com/520lly/qt_data_service/storage/csv"
	//"github.com/520lly/qt_data_service/storage/influxdb"
	"github.com/520lly/qt_data_service/strategies"
	"github.com/520lly/qt_data_service/utils"
)

type TsContext struct {
	TsClients    *map[string]*clients.TsClient
	TsStorage    *map[string]*storage.Storage
	TsStrategies *map[string]*strategies.StockStrategy
	//TsCollector  *map[string]*collector.TsCollector
}

func (tsc *TsContext) InitStockStrategies() {
	if tsc == nil {
		log.Println("tsc is nil!")
		return
	}
	for key, val := range *tsc.TsStorage {
		log.Println("key-val:", key, *val)
		sub := (*val).GetSubscribe()
		fullPath := (*val).GetFullPath()
		//Check if update basic info. condition reached
		ret, files := strategies.CheckUpdateBasic(sub.Period.StockBasic, sub.StockBasic, fullPath)
		if ret {
			tse := strategies.TsEvent{models.DataFlag_Stock_Basic, nil, nil}
			(*(*tsc.TsStrategies)[key]).TsEvent <- tse
		} else {
			//Load the Basic stock infor and start to download history data
			tsc.CheckUpdateDaily(sub.ExchangeName, files[0])
		}
		ret, files = strategies.CheckUpdateBasic(sub.Period.CompanyBasic, sub.CompanyBasic, fullPath)
		if ret {
			tse := strategies.TsEvent{models.DataFlag_Stock_Company, nil, nil}
			(*(*tsc.TsStrategies)[key]).TsEvent <- tse
		} else {
			//Load the Basic stock infor and start to download history data
			//tsc.CheckUpdateDaily(files[0])
		}
	}
}

func (tsc *TsContext) CheckUpdateDaily(exchangeName string, f string) {
	log.Printf("files:[%v]", f)
	_, fp := utils.OpenCsvFile(f)
	if fp != nil {
		var csvReader *csv.Reader
		csvReader = csv.NewReader(fp)
		records, err := csvReader.ReadAll()
		if err == nil {
			for _, r := range records[1:] {
				log.Printf("%T:%v", r, r)
				sym := r[1]
				start := r[10]
				today, end := utils.GetTodayString("20060102")
				_ = today
				//dir := fmt.Sprintf("%s", (*(*tsc.TsStorage)[exchangeName]).GetFullPath()+"/history/")
				//log.Printf("%s", dir)
				//start, end, csvFp := isDataUpdated(dir, sym, listDate)
				//if csvFp == nil {
				//log.Printf("[%]")
				//continue
				//break
				//} else {
				//}
				historyFileNameTillToday := fmt.Sprintf("%s_%s-%s.csv", (*(*tsc.TsStorage)[exchangeName]).GetFullPath()+"/history/"+sym, start, end)
				isNew, csvFp := utils.OpenCsvFile(historyFileNameTillToday)
				log.Printf("[%s] isNew:%v csvFp: %v", historyFileNameTillToday, isNew, csvFp)
				//if isNew {

				//} else {

				//}
				tse := strategies.TsEvent{
					models.DataFlag_Trace_Daily,
					&map[string]string{"ts_code": r[0], "start_date": start, "end_date": end}, csvFp}

				log.Printf("[%s]TsEvent:%v", exchangeName, tsc)
				(*(*tsc.TsStrategies)[exchangeName]).TsEvent <- tse
			}
		}
	}
}

func isDataUpdated(dir string, sym string, listDate string) (start string, end string, fp *os.File) {
	fp = nil
	reStr := fmt.Sprintf("%s", sym+"_"+listDate+".*.csv")
	today, ts := utils.GetTodayString("20060102")
	_ = today
	start = listDate
	end = ts
	diff := utils.CalcDateDiffByDay(start, end)
	log.Printf("diff = %d", diff)
	err, files := utils.FilteredSearchOfDirectoryTree(reStr, dir)
	if err == nil && len(files) == 1 {
		isNew, fp := utils.OpenCsvFile(files[0])
		log.Printf("[%s] isNew:%v csvFp: %v", files[0], isNew, fp)
		if fp != nil {
			if isNew {
			} else {
				//old file and needs to be update
				//csvr := csv.NewReader(fp)
				//lines, err := utils.LineCounter(csvr)
				if err == nil {
					//log.Printf("len[%v]", lines)
				}
				//utils.SeekToLine(csvr, 1)
				fp.Close()
				fp = nil
			}
		} else {
			log.Printf("open file [%s] faild\n", files[0])
		}
	} else {
		log.Printf("No file or more files found with regex [%s]\n", reStr)
	}
	return start, end, fp
}

func TsInit(ctx *context.Context, cfg *config.Config, home string) *TsContext {
	ts_stgs := make(map[string]*strategies.StockStrategy)
	ts_clis := make(map[string]*clients.TsClient)
	ts_stos := make(map[string]*storage.Storage)

	for _, v := range cfg.Subs {
		var sto storage.Storage
		if cfg.Store.Csv {
			sto = mycsv.NewCsvStorage(*ctx, cfg.Store, v, home)
		}
		if cfg.Store.InfluxDB {
			//sto = influxdb.NewInfluxdb(*ctx, v.ExchangeName, v.CurrencyPair, v.ContractType, cfg.Store.InfluxDbCfg.Url, cfg.Store.InfluxDbCfg.Database, cfg.Store.InfluxDbCfg.Username, cfg.Store.InfluxDbCfg.Password)
		}
		go sto.SaveWorker()
		cl := &clients.TsClient{}
		cl = clients.NewTsClient(v.Market, v.ExchangeName, v.CurrencyPair, cfg.Tokens.TuShare)
		stg := strategies.NewStockStrategy(cfg)
		collector.NewTsCollector(ctx, cl, stg, &sto)
		ts_stgs[v.ExchangeName] = stg
		ts_clis[v.ExchangeName] = cl
		ts_stos[v.ExchangeName] = &sto
	}

	var tsc TsContext
	tsc.TsStrategies = &ts_stgs
	tsc.TsClients = &ts_clis
	tsc.TsStorage = &ts_stos

	tsc.InitStockStrategies()

	return &tsc
}
