package configs

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"worker/pkg/logger"

	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/spf13/viper"
)

var (
	passphrase string
	fang       *viper.Viper
)

func init() {
	storagePath, isExists := os.LookupEnv("STORAGE_DIRECTORY_PATH")
	if !isExists {
		logger.Log().Fatal().Msg("$STORAGE_DIRECTORY_PATH is not set. Please set it before running the application.")
	}

	storageInfo, err := os.Stat(storagePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Log().Fatal().Msg("Storage path does not exist")
		}
		logger.Log().Fatal().Err(err).Msgf("Error while reading storage path: %s", err.Error())
	}

	storageMode := storageInfo.Mode()
	if !storageMode.IsDir() {
		logger.Log().Fatal().Msg("Storage path is not a directory")
	}

	if storageMode&0100 != 0100 {
		logger.Log().Fatal().Msg("Storage path is not writable")
	}

	if storageMode&0400 != 0400 {
		logger.Log().Fatal().Msg("Storage path is not readable")
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
			logger.Log().Fatal().Err(err).Msg("Error while reading passphrase")
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
		logger.Log().Fatal().Msg("$CONFIG_DIRECTORY_PATH is not set. Please set it before running the application.")
	}

	dirToScan, err := os.ReadDir(dirPath)
	if err != nil {
		if selectedEnv == "LOCAL" {
			logger.Log().Debug().Msgf("Error on reading $CONFIG_$DIRECTORY_PATH %s: %s", dirPath, err.Error())
			logger.Log().Warn().Msg("Unable to read $CONFIG_DIRECTORY_PATH. Trying to use environment variable as configuration source.")
			fang.AutomaticEnv()
			return
		}
		logger.Log().Debug().Msgf("Error on reading $CONFIG_$DIRECTORY_PATH %s: %s", dirPath, err.Error())
		logger.Log().Fatal().Msg("Unable to read $CONFIG_DIRECTORY_PATH. Please check the directory path before running the application.")
	}

	configCount := 0
	for _, item := range dirToScan {
		if item.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, item.Name())
		bytesContent, err := os.ReadFile(filePath)
		if err != nil {
			logger.Log().Fatal().Msgf("Unable to read file %s: %s", filePath, err.Error())
		}

		if isEncrypted {
			decryptedString, err := helper.DecryptMessageWithPassword(passBytes, string(bytesContent))
			if err != nil {
				logger.Log().Fatal().Msgf("Unable to decrypt config file : %s", err.Error())
			}
			sReader := strings.NewReader(decryptedString)
			ioReader = io.Reader(sReader)
		} else {
			bReader := bytes.NewReader(bytesContent)
			ioReader = io.Reader(bReader)
		}

		err = fang.MergeConfig(ioReader)
		if err != nil {
			logger.Log().Fatal().Msgf("Unable to merge config file %s: %s", filePath, err.Error())
		}

		configCount++
	}

	if configCount <= 0 {
		logger.Log().Fatal().Msg("No configuration file found in $CONFIG_DIRECTORY_PATH. Please check the directory path before running the application.")
	}

}
