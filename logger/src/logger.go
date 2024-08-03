package src

func CreateLog(a_strPath, a_strLogName string) (string, error) {
	var (
		err          error
		strFolderLog string
	)

	strFolderLog, err = createLogFolder(a_strPath)
	if err != nil {
		return strFolderLog, err
	}

	return createLogFile(strFolderLog, a_strLogName)
}
