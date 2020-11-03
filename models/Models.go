package models

import (
   "os"
)

type DataFlag int

const (
	_ DataFlag = iota
	DataFlag_Stock_Basic
	DataFlag_Stock_Company
	DataFlag_Trace_Daily
	DataFlag_Trace_Calendar
)


type TsEvent struct {
   DataFlag  DataFlag           `default: models.DataFlag_Stock_Basic`
   ApiParams *map[string]string `default: nil`
   CsvFile   *os.File           `default: nil`
}

