package src

func StartAppLog(a_strPath string) (string, error) {
	return createLogFolder(a_strPath)
}

func CreateLog(a_strFolderPath, a_strLogName string) (string, error) {
	return createLogFile(a_strFolderPath, a_strLogName)
}

func Log(a_strPath, a_strMethodName, a_strMessage string) error {
	return logFile(a_strPath, a_strMethodName+" : "+a_strMessage)
}

func LogWarning(a_strPath, a_strMethodName, a_strMessage string) error {
	return logFile(a_strPath, "***Warning*** : "+a_strMethodName+" : "+a_strMessage)
}

func LogError(a_strPath, a_strMethodName, a_strMessage string) error {
	return logFile(a_strPath, "***Error*** : "+a_strMethodName+" : "+a_strMessage)
}

func LogException(a_strPath, a_strMethodName, a_strMessage string) error {
	return logFile(a_strPath, "***Exception*** : "+a_strMethodName+" : "+a_strMessage)
}
