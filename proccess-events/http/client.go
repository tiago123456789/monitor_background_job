package http

import (
	"bytes"
	"net/http"
	"process-events/models"
)

func Post(alert models.Alert) {
	req, err := http.NewRequest("POST", alert.Url, bytes.NewBuffer([]byte(alert.Payload)))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
