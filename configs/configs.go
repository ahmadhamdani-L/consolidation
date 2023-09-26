package configs

import (
	"bytes"
	"io"
	"mcash-finance-console-core/pkg/logger"
	"os"
	"path/filepath"
	"strings"

	gpg "github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/spf13/viper"
)

var (
	passphrase string
	fang       *viper.Viper
)

func init() {
	storagePath, isExists := os.LookupEnv("STORAGE_DIRECTORY_PATH")
	if !isExists {
		panic("$STORAGE_DIRECTORY_PATH is not set. Please set it before running the application.")
	}

	storageInfo, err := os.Stat(storagePath)
	if err != nil {
		if !os.IsNotExist(err) {
			logger.Log().Err(err).Msg("Error while reading storage path")
			panic(err)
		}
	}

	storageMode := storageInfo.Mode()
	if !storageMode.IsDir() {
		logger.Log().Err(err).Msg("Storage path is not a directory")
		panic(err)
	}

	if storageMode&0100 != 0100 {
		logger.Log().Err(err).Msg("Storage path is not writable")
		panic(err)
	}

	if storageMode&0400 != 0400 {
		logger.Log().Err(err).Msg("Storage path is not readable")
		panic(err)
	}

	var passBytes []byte
	var ioReader io.Reader

	selectedEnv := strings.ToUpper(os.Getenv("ENV"))

	isEncrypted := false
	if os.Getenv("CONFIG_ENCRYPTED") == "true" {
		strReader := strings.NewReader(passphrase)
		passBytes = make([]byte, strReader.Size())
		_, err := strReader.Read(passBytes)
		if err != nil {
			logger.Log().Err(err).Msg("Error while reading passphrase")
			panic(err)
		}
		isEncrypted = true
	}

	fang = viper.New()
	fang.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	fang.SetConfigType("yaml")

	dirPath, isExists := os.LookupEnv("CONFIG_DIRECTORY_PATH")
	if !isExists {
		if selectedEnv == "LOCAL" {
			logger.Log().Warn().Msg("$CONFIG_DIRECTORY_PATH is not set. Trying to use environment variable as configuration source.")
			fang.AutomaticEnv()
			return
		}
		panic("$CONFIG_DIRECTORY_PATH is not set. Please set it before running the application.")
	}

	dirToScan, err := os.ReadDir(dirPath)
	if err != nil {
		if selectedEnv == "LOCAL" {
			logger.Log().Debug().Msgf("Error on reading $CONFIG_$DIRECTORY_PATH %s: %s", dirPath, err.Error())
			logger.Log().Warn().Msg("Unable to read $CONFIG_DIRECTORY_PATH. Trying to use environment variable as configuration source.")
			fang.AutomaticEnv()
			return
		}
		logger.Log().Err(err).Msg("Unable to read $CONFIG_DIRECTORY_PATH. Please check the directory path before running the application.")
		panic(err)
	}

	configCount := 0
	for _, item := range dirToScan {
		if item.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, item.Name())
		bytesContent, err := os.ReadFile(filePath)
		if err != nil {
			logger.Log().Err(err).Msg("Unable to read configuration file")
			panic(err)
		}

		if isEncrypted {
			decryptedString, err := gpg.DecryptMessageWithPassword(passBytes, string(bytesContent))
			if err != nil {
				logger.Log().Err(err).Msg("Unable to decrypt config file")
				panic(err)
			}
			sReader := strings.NewReader(decryptedString)
			ioReader = io.Reader(sReader)
		} else {
			bReader := bytes.NewReader(bytesContent)
			ioReader = io.Reader(bReader)
		}

		err = fang.MergeConfig(ioReader)
		if err != nil {
			logger.Log().Err(err).Msg("Unable to merge config file")
			panic(err)
		}

		configCount++
	}

	if configCount <= 0 {
		panic("No configuration file found in $CONFIG_DIRECTORY_PATH. Please check the directory path before running the application.")
	}
}
