package env

import "os"

func ReadJson(path string) []byte {
	var file []byte
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic(err)
	}
	file, _ = os.ReadFile(path)
	return file
}