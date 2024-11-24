package src

import (
	"errors"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	c_strSystemLogs   = "SystemLogs"
	c_strLogExtension = ".log"
)

func createLogFile(a_LogInfo LogInfoType, a_strLogName string) (LogInfoType, error) {
	var (
		err         error
		file        *os.File
		strFullPath string
		strFileName string
	)
	// Verifica se parametros sao validos
	if a_LogInfo.strLogFolder == "" || a_strLogName == "" {
		return a_LogInfo, errors.New("parameter is empty")
	}

	strFileName = a_strLogName + c_strLogExtension
	strFullPath = filepath.Join(a_LogInfo.strLogFolder, strFileName)

	// Verifica se arquivo de log ja existe
	_, err = os.Stat(strFullPath)
	if err == nil {
		return a_LogInfo, errors.New("log already exists")
	}

	// Cria arquivo de log
	file, err = os.OpenFile(strFullPath, os.O_CREATE, 0644)
	if err != nil {
		return a_LogInfo, err
	}

	defer file.Close()

	a_LogInfo.lstLogFiles = append(a_LogInfo.lstLogFiles, strFullPath)

	return a_LogInfo, nil
}

func createLogFolder(a_strPath string) (LogInfoType, error) {
	var (
		err         error
		strFolder   string
		strFullPath string
		dtNow       time.Time
	)
	// Verifica se parametros sao validos
	if a_strPath == "" {
		return LogInfoType{}, errors.New("parameter is empty")
	}

	dtNow = time.Now()

	strFolder = c_strSystemLogs + fmt.Sprintf("%02d%02d%d", dtNow.Day(), dtNow.Month(), dtNow.Year())

	strFullPath = filepath.Join(a_strPath, strFolder)

	// Verifica se pasta de log ja existe, entao renomeia adicionando um _N
	_, err = os.Stat(strFullPath)
	if err == nil {
		err = renameOldLogFolder(a_strPath, strFolder, strFullPath)
		if err != nil {
			return LogInfoType{}, err
		}
	}

	// Forca a criacao do diretorio caso nao tenha sido criado
	err = os.MkdirAll(strFullPath, os.ModePerm)
	if err != nil {
		return LogInfoType{}, err
	}

	return LogInfoType{
		strLogFolder: strFullPath,
		lstLogFiles:  make([]string, 0),
	}, nil
}

func renameOldLogFolder(a_strPath, a_strFolder, a_strFullPath string) error {
	var (
		err          error
		dir          fs.DirEntry
		arrDir       []fs.DirEntry
		strFolder    string
		strNewFolder string
		nCount       int
	)
	arrDir, err = os.ReadDir(a_strPath)

	if err != nil {
		return err
	}

	// Itera sobre pastas de logs "iguais" para encontrar a quantidade
	nCount = 0
	for _, dir = range arrDir {
		strFolder = filepath.Base(dir.Name())
		strFolder = strings.Split(strFolder, "_")[0]

		if dir.IsDir() && strFolder == a_strFolder {
			nCount++
		}
	}

	strNewFolder = strFolder + "_" + strconv.Itoa(nCount)

	// Renomeia pasta de log antiga
	err = os.Rename(a_strFullPath, filepath.Join(a_strPath, strNewFolder))
	if err != nil {
		return err
	}

	return nil
}

func getLogsPath(a_LogInfo LogInfoType, a_strPath string) string {
	var (
		strPath     string
		strFileName string
	)
	for _, strPath = range a_LogInfo.lstLogFiles {
		// Obtem nome do arquivo de log
		strFileName = filepath.Base(strPath)
		// Verifica se nome eh igual ao parametro e retorna
		if strFileName == a_strPath+c_strLogExtension {
			return strPath
		}
	}
	// Retorna vazio caso nao tenha encontrado o log
	return ""
}

func logFile(a_strPath, a_strMessage string) error {
	var (
		err            error
		file           *os.File
		strFullMessage string
	)
	// Abre arquivo de log
	file, err = os.OpenFile(a_strPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	strFullMessage = getFullTimestamp() + " : " + a_strMessage + "\n"

	_, err = file.WriteString(strFullMessage)

	defer file.Close()

	return err
}

func getFullTimestamp() string {
	var (
		dtNow time.Time
		sMs   float64
	)
	dtNow = time.Now()

	// Calcula milissegundos
	sMs = float64(dtNow.UnixMilli()) / 1e3
	sMs = (sMs - math.Floor(sMs)) * 1e3

	return fmt.Sprintf("%02d-%02d-%d %02d:%02d:%02d.%.0f", dtNow.Day(), dtNow.Month(), dtNow.Year(), dtNow.Hour(), dtNow.Minute(), dtNow.Second(), sMs)
}
