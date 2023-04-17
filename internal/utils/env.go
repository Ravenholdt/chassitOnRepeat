package utils

import (
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

func validateEnv(env string) {
	if GetStringEnv(env, "") == "" {
		log.Fatal().Str("tag", "env").Msgf("You must set your '%s' environmental variable.", env)
	}
}

// ValidateEnvs Required envs
func ValidateEnvs() {
	validateEnv("MONGODB_URI")
}
