package src

const (
	c_strDataFolder  = "/data"
	c_strLogsFolder  = "/logs"
	c_strInputFolder = "/input"

	c_strDateFormat = "%.4d-%.2d-%.2d"                  // yyyy-mm-dd
	c_strTradeFile  = c_strDateFormat + "_%s_TRADE.csv" // yyyy-mm-dd_<TICKER>_TRADE.csv
	c_strBuyFile    = c_strDateFormat + "_%s_BUY.csv"   // yyyy-mm-dd_<TICKER>_BUY.csv
	c_strSellFile   = c_strDateFormat + "_%s_SELL.csv"  // yyyy-mm-dd_<TICKER>_SELL.csv

	c_strCustomTimestampLayout  = "2006-01-02T15:04:05.999"
	c_strCustomTimestampLayout2 = "2006-01-02 15:04:05.999"
)
