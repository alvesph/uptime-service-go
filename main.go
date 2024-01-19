package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
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
		fmt.Println(err)
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	statusMap := make(map[string]*StatusInfo)
	for _, service := range urlList.URLs {
		statusMap[service.URL] = &StatusInfo{URL: service.URL, LastStatus: http.StatusOK, LastSentTime: time.Now()}
	}

	for {
		now := time.Now()

		done := make(chan bool)

		var services_error []string
		for _, service := range urlList.URLs {
			go func(s readfile.Service) {
				req, err := http.NewRequest("GET", s.URL, nil)
				if err != nil {
					fmt.Printf("Erro ao criar a requisição para o serviço [%s] URL: [%s]\n", s.Name, s.URL)
					done <- true
					return
				}

				get, err := client.Do(req)
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
					services_error = append(services_error, message)
					statusInfo.LastStatus = status
					statusInfo.LastSentTime = time.Now()
				}
				fmt.Printf("Service: [%s] Status: [%d] Charging time: [%f]\n", s.Name, status, elapsedTime)
				done <- true
			}(service)
		}

		for range urlList.URLs {
			<-done
		}

		if len(services_error) > 0 {
			fmt.Println("Mensagem enviada com sucesso para o Teams!")
			err = teams.SendTeamsMessage(strings.Join(services_error, "\n"))
			if err != nil {
				fmt.Printf("Erro ao enviar mensagem para o Teams: %v\n", err)
				return
			}
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
