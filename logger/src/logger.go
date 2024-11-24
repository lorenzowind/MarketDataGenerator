package src

func StartAppLog(a_strPath string) (LogInfoType, error) {
	return createLogFolder(a_strPath)
}

func CreateLog(a_LogInfo LogInfoType, a_strLogName string) (LogInfoType, error) {
	return createLogFile(a_LogInfo, a_strLogName)
}

func Log(a_LogInfo LogInfoType, a_strPath, a_strMethodName, a_strMessage string) error {
	return logFile(getLogsPath(a_LogInfo, a_strPath), a_strMethodName+" : "+a_strMessage)
}

func LogWarning(a_LogInfo LogInfoType, a_strPath, a_strMethodName, a_strMessage string) error {
	return logFile(getLogsPath(a_LogInfo, a_strPath), "***Warning*** : "+a_strMethodName+" : "+a_strMessage)
}

func LogError(a_LogInfo LogInfoType, a_strPath, a_strMethodName, a_strMessage string) error {
	return logFile(getLogsPath(a_LogInfo, a_strPath), "***Error*** : "+a_strMethodName+" : "+a_strMessage)
}

func LogException(a_LogInfo LogInfoType, a_strPath, a_strMethodName, a_strMessage string) error {
	return logFile(getLogsPath(a_LogInfo, a_strPath), "***Exception*** : "+a_strMethodName+" : "+a_strMessage)
}
