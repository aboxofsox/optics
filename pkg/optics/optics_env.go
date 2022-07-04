package optics

import "github.com/joho/godotenv"

func readEnv() map[string]string {
	var envMap map[string]string
	if envExists() {
		envMap, _ = godotenv.Read()
	}
	return envMap
}
