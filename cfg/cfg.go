package cfg

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var envs map[string]string

// Getenv get setting from env or .env
func Getenv(k string) string {
	if len(envs) == 0 {
		initEnvs()
	}
	return envs[k]
}
func initEnvs() {

	envs = make(map[string]string)
	err := godotenv.Load()
	if err != nil {
		fmt.Println("cant load .env")
	}
	envs["CERT2019_PORT"] = os.Getenv("CERT2019_PORT")
	if envs["CERT2019_PORT"] == "" {
		envs["CERT2019_PORT"] = "3000"
	}
	envs["CERT2019_MONGO_URL"] = os.Getenv("CERT2019_MONGO_URL")
	if envs["CERT2019_MONGO_URL"] == "" {
		envs["CERT2019_MONGO_URL"] = "mongodb://localhost:27017/certimg"
	}
	envs["CERT2019_MONGO_DB_NAME"] = os.Getenv("CERT2019_MONGO_DB_NAME")
	if envs["CERT2019_MONGO_DB_NAME"] == "" {
		envs["CERT2019_MONGO_DB_NAME"] = "certimg"
	}
}
