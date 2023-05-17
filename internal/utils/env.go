package utils

import (
	"errors"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

func GetStringEnv(key string, fallback string) string {
	if env, ok := os.LookupEnv(key); ok {
		return env
	}
	return fallback
}

func GetBoolEnv(key string, fallback bool) bool {
	parseBool, err := strconv.ParseBool(GetStringEnv(key, "INVALID_BOOL"))
	if err != nil {
		return fallback
	}
	return parseBool
}

func GetIntEnv(key string, fallback int) int {
	parseInt, err := strconv.ParseInt(GetStringEnv(key, "INVALID_INT"), 10, 32)
	if err != nil {
		return fallback
	}
	return int(parseInt)
}

func validateEnv(env string) {
	if GetStringEnv(env, "") == "" {
		log.Fatal().Str("tag", "env").Msgf("You must set your '%s' environmental variable.", env)
	}
}

// ValidateEnvs Required envs
func ValidateEnvs() {
	validateEnv("MONGODB_URI")

	filesPath := GetStringEnv("FILES_PATH", "/files")
	if files, err := os.Stat(filesPath); !errors.Is(err, os.ErrNotExist) {
		if !files.IsDir() {
			log.Fatal().Str("tag", "env").Str("path", filesPath).Msgf("FILES_PATH' is no a folder")
		}
	} else {
		log.Fatal().Str("tag", "env").Str("path", filesPath).Err(err).Msgf("Folder defined in 'FILES_PATH' does not exist")
	}
}
