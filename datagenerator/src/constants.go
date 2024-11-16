package src

const (
	c_strDataFolder      = "/data"
	c_strLogsFolder      = "/logs"
	c_strInputFolder     = "/input"
	c_strReferenceFolder = "/reference"

	c_strReferenceDateFormat = "%.2d%.2d%.4d"                              // ddmmyyyy
	c_strReferenceBuyFile    = c_strReferenceDateFormat + "_%s_compra.csv" // ddmmyyyy_<TICKER>_compra.csv
	c_strReferenceSellFile   = c_strReferenceDateFormat + "_%s_venda.csv"  // ddmmyyyy_<TICKER>_venda.csv

	c_strDateFormat = "%.4d-%.2d-%.2d"                 // yyyy-mm-dd
	c_strBuyFile    = c_strDateFormat + "_%s_BUY.csv"  // yyyy-mm-dd_<TICKER>_BUY.csv
	c_strSellFile   = c_strDateFormat + "_%s_SELL.csv" // yyyy-mm-dd_<TICKER>_SELL.csv

	c_strBenchmarksFile = "BENCHMARKS.csv"

	c_strCustomDateFormat      = "%.4d-%.2d-%.2d"
	c_strCustomTimestampFormat = c_strCustomDateFormat + "T%.2d:%.2d:%.2d.%.3d"

	c_strCustomTimestampLayout  = "2006-01-02T15:04:05.999"
	c_strCustomTimestampLayout2 = "2006-01-02 15:04:05.999"
	c_strCustomTimestampLayout3 = "15:04:05.999"
)
