package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"uptimeService/readfile"
	"uptimeService/teams"
)

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

	filePath := "listaurl.json"
	urlList, err := readfile.ReadURLList(filePath)
	if err != nil {
		fmt.Println("Erro ao carregar a lista de urls:", err)
		return
	}

	statusMap := make(map[string]*StatusInfo)
	for _, service := range urlList.URLs {
		statusMap[service.URL] = &StatusInfo{URL: service.URL, LastStatus: http.StatusOK, LastSentTime: time.Now()}
	}

	for {
		now := time.Now()

		done := make(chan bool)

		for _, service := range urlList.URLs {
			go func(s readfile.Service) {
				get, err := http.Get(s.URL)
				if err != nil {
					fmt.Printf("Ocorreu um erro ao executar o serviço [%s] URL: [%s]\n", s.Name, s.URL)
					done <- true
					return
				}

				elapsedTime := time.Since(now).Seconds()
				status := get.StatusCode

				statusInfo := statusMap[s.URL]
				fmt.Println(statusInfo.LastStatus)
				if status != statusInfo.LastStatus {
					message := fmt.Sprintf("Service: [%s] Status: [%d] Charging time: [%f]\n", s.Name, status, elapsedTime)
					fmt.Println("Mensagem enviada com sucesso para o Teams!")
					err = teams.SendTeamsMessage(message)
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
