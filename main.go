package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"uptimeService/teams"
)

type Service struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type URLlist struct {
	URLs []Service `json:"urls"`
}

type StatusInfo struct {
	URL          string
	LastStatus   int
	LastSentTime time.Time
}

func main() {

	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Println("Erro ao carregar o arquivo .env:", err)
	// 	return
	// }

	// urlListJSON := os.Getenv("URL_LIST")

	// var urlList URLlist
	// err = json.Unmarshal([]byte(urlListJSON), &urlList)
	// if err != nil {
	// 	fmt.Println("Erro ao decodificar a lista de URLs JSON:", err)
	// 	return
	// }

	file, err := os.Open(os.Getenv("listaurl.json")) // Substitua pelo caminho real do seu arquivo
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return
	}
	defer file.Close()

	var urlList URLlist
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&urlList)
	if err != nil {
		fmt.Println("Erro ao decodificar a lista de URLs JSON:", err)
		return
	}

	statusMap := make(map[string]*StatusInfo)
	for _, service := range urlList.URLs {
		statusMap[service.URL] = &StatusInfo{URL: service.URL, LastStatus: http.StatusOK, LastSentTime: time.Now()}
	}

	for {
		now := time.Now()
		for _, service := range urlList.URLs {
			get, err := http.Get(service.URL)
			if err != nil {
				fmt.Printf("Ocorreu um erro ao executar o serviço [%s] URL: [%s]\n", service.Name, service.URL)
				continue
			}
			elapsed_time := time.Since(now).Seconds()
			status := get.StatusCode

			statusInfo := statusMap[service.URL]
			if status != statusInfo.LastStatus {
				message := fmt.Sprintf("Service: [%s] Status: [%d] Charging time: [%f]\n", service.Name, status, elapsed_time)
				err = teams.SendTeamsMessage(message)
				fmt.Println("Mensagem enviada com sucesso para o Teams!")
				if err != nil {
					fmt.Printf("Erro ao enviar mensagem para o Teams: %v\n", err)
					return
				}

				statusInfo.LastStatus = status
				statusInfo.LastSentTime = time.Now()

			}
			fmt.Printf("Service: [%s] Status: [%d] Charging time: [%f]\n", service.Name, status, elapsed_time)
		}

		sleepDurationStr := os.Getenv("RUNTIME")
		sleepDuration, err := strconv.Atoi(sleepDurationStr)
		if err != nil {
			fmt.Println("Erro ao converter o tempo de espera para um número:", err)
			return
		}
		time.Sleep(time.Duration(sleepDuration) * time.Second)
	}

}
