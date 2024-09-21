package src

const (
	c_strDataFolder  = "/data"
	c_strLogsFolder  = "/logs"
	c_strInputFolder = "/input"
	c_strDateFormat  = "%d-%d-%d"                        // dd-mm-yyyy
	c_strTradeFile   = c_strDateFormat + "_%s_BUY.csv"   // dd-mm-yyyy_<TICKER>_TRADE.csv
	c_strBuyFile     = c_strDateFormat + "_%s_SELL.csv"  // dd-mm-yyyy_<TICKER>_BUY.csv
	c_strSellFile    = c_strDateFormat + "_%s_TRADE.csv" // dd-mm-yyyy_<TICKER>_SELL.csv
)
