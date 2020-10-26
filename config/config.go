package config

import "github.com/goex-top/market_center"

type Subscribe struct {
	Market       string                 `json:"market" yaml:"market" default:""`
	ExchangeName string                 `json:"exchange_name" yaml:"exchange_name" default:""`
	CurrencyPair string                 `json:"currency_pair" yaml:"currency_pair" default:""`
	ContractType string                 `json:"contract_type,omitempty" yaml:"contract_type" default:""`
	Period       int64                  `json:"period" yaml:"period" default:"100"`
	Flag         market_center.DataFlag `json:"flag" yaml:"flag" default:"1"`
}

type Storage struct {
	Csv         bool `json:"csv" yaml:"csv" `
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
   TuShare     string `json:"tushare" yaml:"tushare" default:""`
}
type Config struct {
	Subs             []Subscribe `json:"subs" yaml:"subs" default:"subs"`
   Store            Storage     `json:"store" yaml:"store" default:""`
   Tokens           Tokens      `json:"tokens" yaml:"tokens" default:""`
	WithMarketCenter bool        `json:"with_market_center" yaml:"with_market_center" `
	MarketCenterPath string      `json:"market_center_path" yaml:"market_center_path" default:"/tmp/goex.market.center"`
}
