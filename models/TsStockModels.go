package models

import (
	"fmt"
)

type ExchangeNameBasic int
const (
	_ ExchangeNameBasic = iota
   ExchangeNameBasic_SSE
   ExchangeNameBasic_SZSE
   ExchangeNameBasic_CFFEX
   ExchangeNameBasic_SHFE
   ExchangeNameBasic_CZCE
   ExchangeNameBasic_DCE
   ExchangeNameBasic_INE
   ExchangeNameBasic_IB
   ExchangeNameBasic_XHKG
)


type FieldStockBasic int
const (
	_ FieldStockBasic = iota
	FieldStockBasic_Ts_Code
	FieldStockBasic_Symbol
	FieldStockBasic_Name
	FieldStockBasic_Area
	FieldStockBasic_Industry
	FieldStockBasic_Fullname
	FieldStockBasic_Market
	FieldStockBasic_Exchange
	FieldStockBasic_CurrType
	FieldStockBasic_ListStatus
	FieldStockBasic_ListDate
	FieldStockBasic_DelistDate
	FieldStockBasic_IsHs
)

//type StockBasics truct {
	//TsCode     string
	//Symbol     string
	//Name       string
	//Area       string
	//Industry   string
	//Fullname   string
	//Market     string
	//Exchange   string
	//CurrType   string
	//ListStatus string
	//ListDate   string
	//DelistDate string
	//IsHs       string
//}

const (
   StockBasic string = "stock_basic"
   CompanyBasic string = "company_basic"
   TradeCalendar string = "trade_cal"
)

var BasicFileNames = []string{StockBasic, CompanyBasic, TradeCalendar}

func (s DataFlag) String() (sym string) {
	if s > 0 && int(s) <= len(BasicFileNames) {
		return BasicFileNames[s-1]
	}
	return fmt.Sprintf("UNKNOWN_METHOD_TYPE (%d)", s)
}


var StockFieldSymbol = []string{"ts_code", "symbol", "name", "area", "industry", "fullname", "market", "exchange", "curr_type", "list_status", "list_date", "delist_date", "is_hs"}

//func (s FieldStockBasic) String() (sym string) {
	//if s > 0 && int(s) <= len(FieldSymbol) {
		//return FieldSymbol[s-1]
	//}
	//return fmt.Sprintf("UNKNOWN_METHOD_TYPE (%d)", s)
//}

var CompanyFieldSymbol = []string{"ts_code", "exchange", "chairman", "manager", "secretary", "reg_capital", "setup_date", "province", "city", "introduction", "website", "email", "office", "main_business", "employees", "business_scope"}

var TradeDailyFieldSymbol = []string{"ts_code", "trade_date", "open", "high", "low", "close", "pre_close", "change", "pct_chg", "vol", "amount"}
var TradeCalendarFieldSymbol = []string{"exchange", "cal_date", "is_open", "pretrade_date"}
