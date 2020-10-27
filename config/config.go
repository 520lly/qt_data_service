package config

import "github.com/520lly/qt_data_service/models"

type Subscribe struct {
	Market       string `json:"market" yaml:"market" default:""`
	ExchangeName string `json:"exchange_name" yaml:"exchange_name" default:""`
	CurrencyPair string `json:"currency_pair" yaml:"currency_pair" default:""`
	ContractType string `json:"contract_type,omitempty" yaml:"contract_type" default:""`
	Period       struct {
		StockBasic   int64 `json:"stock_basic" yaml:"stock_basic" default:"30"`
		CompanyBasic int64 `json:"company_basic" yaml:"company_basic" default:"30"`
		TradeDaily   int64 `json:"trade_daily" yaml:"trade_daily" default:"120"`
	}
	Flag         models.DataFlag `json:"flag" yaml:"flag" default:"1"`
	StockBasic   string          `json:"stock_basic" yaml:"stock_basic" default:"stock_basic"`
	CompanyBasic string          `json:"company_basic" yaml:"company_basic" default:"company_basic"`
}

type Storage struct {
	Csv    bool `json:"csv" yaml:"csv" `
	CsvCfg struct {
		Location string `json:"location" yaml:"location" default:"qt_data"`
	} `json:"csv_cfg" yaml:"csv_cfg"`
	InfluxDB    bool `json:"influx_db" yaml:"influx_db" `
	InfluxDbCfg struct {
		Url      string `json:"url" yaml:"url" default:"http://localhost:8086"`
		Database string `json:"database" yaml:"database" default:"market_data"`
		Username string `json:"username" yaml:"username" default:""`
		Password string `json:"password" yaml:"password" default:""`
	} `json:"influx_db_cfg" yaml:"influx_db_cfg"`
	// TBD
}

type Tokens struct {
	TuShare string `json:"tushare" yaml:"tushare" default:""`
}

type Config struct {
	Subs             []Subscribe `json:"subs" yaml:"subs" default:"subs"`
	Store            Storage     `json:"store" yaml:"store" default:""`
	Tokens           Tokens      `json:"tokens" yaml:"tokens" default:""`
	WithMarketCenter bool        `json:"with_market_center" yaml:"with_market_center" `
	MarketCenterPath string      `json:"market_center_path" yaml:"market_center_path" default:"/tmp/goex.market.center"`
}
