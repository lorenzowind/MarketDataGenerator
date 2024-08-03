package src

import (
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	c_strLogExtension = ".log"
)

func TryCreateLogFile(a_strPath, a_strFileName string) bool {
	var (
		err         error
		file        *os.File
		strFileName string
		strFullPath string
		bRenameOld  bool
	)

	// Verifica se parametros sao validos
	if a_strPath == "" || a_strFileName == "" {
		return false
	}

	strFileName = a_strFileName + c_strLogExtension

	strFullPath = filepath.Join(a_strPath, strFileName)

	// Verifica de arquivo de log ja existe, entao renomeia adicionado um _N
	_, err = os.Stat(strFullPath)
	if err == nil {
		bRenameOld = tryRenameOldLogFile(a_strPath, a_strFileName, strFullPath)
		if !bRenameOld {
			return false
		}
	}

	// Forca a criacao do diretorio caso nao tenha sido criado
	err = os.MkdirAll(a_strPath, os.ModePerm)
	if err != nil {
		return false
	}

	// Cria arquivo de log
	file, err = os.OpenFile(strFullPath, os.O_CREATE, 0644)
	if err != nil {
		return false
	}

	defer file.Close()

	return true
}

func tryRenameOldLogFile(a_strPath, a_strFileName, a_strFullPath string) bool {
	var (
		err            error
		dir            fs.DirEntry
		arrDir         []fs.DirEntry
		strFile        string
		strFileName    string
		strNewFileName string
		nCount         int
	)

	arrDir, err = os.ReadDir(a_strPath)

	if err != nil {
		return false
	}

	// Itera sobre arquivos de logs "iguais" para encontrar a quantidade
	nCount = 0
	for _, dir = range arrDir {
		strFile = filepath.Base(dir.Name())

		strFileName = strings.TrimSuffix(strFile, filepath.Ext(strFile))
		strFileName = strings.Split(strFileName, "_")[0]

		if !dir.IsDir() && strFileName == a_strFileName {
			nCount++
		}
	}

	strNewFileName = a_strFileName + "_" + strconv.Itoa(nCount) + c_strLogExtension

	// Renomeia arquivo de log antigo
	err = os.Rename(a_strFullPath, filepath.Join(a_strPath, strNewFileName))

	return err == nil
}
