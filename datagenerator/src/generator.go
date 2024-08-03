package src

import (
	logger "marketmanipulationdetector/logger/src"
)

func Start() {
	logger.CreateLog(getLogsPath(), "Main")
}
