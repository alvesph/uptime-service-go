package readfile

import (
	"encoding/json"
	"fmt"
	"os"
)

type Service struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type URLlist struct {
	URLs []Service `json:"urls"`
}

func ReadURLList(filePath string) (URLlist, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return URLlist{}, fmt.Errorf("erro ao abrir o arquivo: %v", err)
	}
	defer file.Close()

	var urlList URLlist
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&urlList)
	if err != nil {
		return URLlist{}, fmt.Errorf("erro ao decodificar a lista de URLs JSON: %v", err)
	}

	return urlList, nil
}
