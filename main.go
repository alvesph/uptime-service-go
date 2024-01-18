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
	file, err := os.Open("listaurl.json")
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

		// Usar um canal para sincronizar goroutines
		done := make(chan bool)

		for _, service := range urlList.URLs {
			go func(s Service) {
				get, err := http.Get(s.URL)
				if err != nil {
					fmt.Printf("Ocorreu um erro ao executar o serviço [%s] URL: [%s]\n", s.Name, s.URL)
					done <- true
					return
				}

				elapsedTime := time.Since(now).Seconds()
				status := get.StatusCode

				statusInfo := statusMap[s.URL]
				if status != statusInfo.LastStatus {
					message := fmt.Sprintf("Service: [%s] Status: [%d] Charging time: [%f]\n", s.Name, status, elapsedTime)
					err = teams.SendTeamsMessage(message)
					fmt.Println("Mensagem enviada com sucesso para o Teams!")
					if err != nil {
						fmt.Printf("Erro ao enviar mensagem para o Teams: %v\n", err)
						return
					}

					statusInfo.LastStatus = status
					statusInfo.LastSentTime = time.Now()
				}

				fmt.Printf("Service: [%s] Status: [%d] Charging time: [%f]\n", s.Name, status, elapsedTime)
				done <- true
			}(service)
		}

		// Aguardar todas as goroutines concluírem
		for range urlList.URLs {
			<-done
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
