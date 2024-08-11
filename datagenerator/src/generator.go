package src

import (
	logger "marketmanipulationdetector/logger/src"
)

func Start() {
	const (
		c_strMethodName = "generator.Start"
	)
	var (
		strLogFolder string
		strLogFile   string
		err          error
	)

	strLogFolder, err = logger.StartAppLog(getLogsPath())
	if err != nil {
		panic("log folder can not be created")
	}

	strLogFile, err = logger.CreateLog(strLogFolder, "Main")
	if err != nil {
		panic("log file can not be created")
	}

	logger.Log(strLogFile, c_strMethodName, "Test")
}
