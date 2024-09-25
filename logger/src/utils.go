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

func createLogFile(a_strFolderPath, a_strLogName string) (string, error) {
	var (
		err         error
		file        *os.File
		strFullPath string
		strFileName string
	)

	// Verifica se parametros sao validos
	if a_strFolderPath == "" || a_strLogName == "" {
		return "", errors.New("parameter is empty")
	}

	strFileName = a_strLogName + c_strLogExtension
	strFullPath = filepath.Join(a_strFolderPath, strFileName)

	// Verifica se arquivo de log ja existe
	_, err = os.Stat(strFullPath)
	if err == nil {
		return "", errors.New("log already exists")
	}

	// Cria arquivo de log
	file, err = os.OpenFile(strFullPath, os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}

	defer file.Close()

	return strFullPath, nil
}

func createLogFolder(a_strPath string) (string, error) {
	var (
		err         error
		strFolder   string
		strFullPath string
		dtNow       time.Time
	)

	// Verifica se parametros sao validos
	if a_strPath == "" {
		return "", errors.New("parameter is empty")
	}

	dtNow = time.Now()

	strFolder = c_strSystemLogs + fmt.Sprintf("%02d%02d%d", dtNow.Day(), dtNow.Month(), dtNow.Year())

	strFullPath = filepath.Join(a_strPath, strFolder)

	// Verifica se pasta de log ja existe, entao renomeia adicionando um _N
	_, err = os.Stat(strFullPath)
	if err == nil {
		err = renameOldLogFolder(a_strPath, strFolder, strFullPath)
		if err != nil {
			return "", err
		}
	}

	// Forca a criacao do diretorio caso nao tenha sido criado
	err = os.MkdirAll(strFullPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return strFullPath, nil
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
