package teams

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func SendTeamsMessage(message string) error {
	teamsWebhookURL := os.Getenv("TEAMS_WEBHOOK_URL")
	if teamsWebhookURL == "" {
		return nil
	}

	formattedMessage := fmt.Sprintf("```css\n%s\n```", message)

	payload := map[string]interface{}{
		"text": formattedMessage,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(teamsWebhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro ao enviar mensagem para o Teams. Código de status: %d", resp.StatusCode)
	}

	return nil
}
