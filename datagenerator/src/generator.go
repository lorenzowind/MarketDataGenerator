package src

import (
	logger "marketmanipulationdetector/logger/src"
	"path/filepath"
	"runtime"
)

func Start() {
	logger.TryCreateLogFile(getLogsPath(), "SystemLogs")
}

func getLogsPath() string {
	const (
		c_strDataFolder = "/data"
		c_strLogsFolder = "/logs"
	)
	var (
		strFile     string
		strRootPath string
	)

	_, strFile, _, _ = runtime.Caller(0)

	strRootPath = filepath.Join(filepath.Dir(strFile), "..")

	return strRootPath + c_strDataFolder + c_strLogsFolder
}
