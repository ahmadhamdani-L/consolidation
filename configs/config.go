package configs

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"notification/pkg/logger"

	gpg "github.com/ProtonMail/gopenpgp/v2/helper"

	"github.com/spf13/viper"
)

var (
	passphrase string
	fang       *viper.Viper
)

func init() {
	var passBytes []byte
	var ioReader io.Reader

	selectedEnv := strings.ToUpper(os.Getenv("ENV"))
	isEncrypted := false
	if os.Getenv("CONFIG_ENCRYPTED") == "true" {
		strReader := strings.NewReader(passphrase)
		passBytes = make([]byte, strReader.Size())
		_, err := strReader.Read(passBytes)
		if err != nil {
			logger.Log().Error().Err(err).Msg("Failed to read passphrase")
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
