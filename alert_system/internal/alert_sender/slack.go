package alertsender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"alert_system/internal/models"
)


type SlackAlertSender struct{
	webhookURL string

}

func NewSlackAlertSender(webhookURL string)AlertSender{
	return &SlackAlertSender{webhookURL:webhookURL}

}

// sendSlackAlert sends an error message to the specified Slack webhook URL
func  (sa SlackAlertSender)SendAlert(message string) error {
	payload := models.SlackPayload{Text: message}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(sa.webhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-ok response from Slack: %s", resp.Status)
	}

	return nil
}

