package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const avatarURL = "https://www.pci-group.com/wp-content/uploads/CineplexLogoSq.jpg"
const contentType = "application/json"

type payload struct {
	AvatarURL string `json:"avatar_url"`
	Content   string `json:"content"`
}

func Send(url string, message string) error {
	payload, err := json.Marshal(payload{Content: message, AvatarURL: avatarURL})
	if err != nil {
		return fmt.Errorf("error marshalling json: %w", err)
	}

	resp, err := http.Post(url, contentType, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("error sending webhook: %w", err)
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("bad response from webhook: %s", resp.Status)
	}

	return nil
}
