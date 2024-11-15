package src

const (
	c_strDataFolder      = "/data"
	c_strLogsFolder      = "/logs"
	c_strInputFolder     = "/input"
	c_strReferenceFolder = "/reference"

	c_strDateFormat = "%.2d%.2d%.4d"                     // ddmmyyyy
	c_strBuyFile    = c_strDateFormat + "_%s_compra.csv" // ddmmyyyy_<TICKER>_compra.csv
	c_strSellFile   = c_strDateFormat + "_%s_venda.csv"  // ddmmyyyy_<TICKER>_venda.csv

	c_strBenchmarksFile = "BENCHMARKS.csv"
)
