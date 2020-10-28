package models

import (
	//"fmt"
)

type DataFlag int

const (
	_ DataFlag = iota
	DataFlag_Stock_Basic
	DataFlag_Stock_Company
	DataFlag_Trace_Daily
)

