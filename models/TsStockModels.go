package models

import (
	"fmt"
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

type StockBasic struct {
	TsCode     string
	Symbol     string
	Name       string
	Area       string
	Industry   string
	Fullname   string
	Market     string
	Exchange   string
	CurrType   string
	ListStatus string
	ListDate   string
	DelistDate string
	IsHs       string
}

var FieldSymbol = []string{"ts_code", "symbol", "name", "area", "industry", "fullname", "market", "exchange", "curr_type", "list_status", "list_date", "delist_date", "is_hs"}

func (s FieldStockBasic) String() (sym string) {
	if s > 0 && int(s) <= len(FieldSymbol) {
		return FieldSymbol[s-1]
	}
	return fmt.Sprintf("UNKNOWN_METHOD_TYPE (%d)", s)
}

var CompanyFieldSymbol = []string{"ts_code", "exchange", "chairman", "manager", "secretary", "reg_capital", "setup_date", "province", "city", "introduction", "website", "email", "office", "main_business", "employees", "business_scope"}

var TradeDailyFieldSymbol = []string{"ts_code", "trade_date", "open", "high", "low", "close", "pre_close", "change", "pct_chg", "vol", "amount"}