package main

import "net/http"

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From string `json:"from"`
		To string `json:"To"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	if err := app.readJSON(w, r, &requestPayload); err != nil {
		app.errorJSON(w, err)
		return 
	}

	msg := Message{
		From: requestPayload.From,
		To: requestPayload.To,
		Subject: requestPayload.Subject,
		Data: requestPayload.Message,
	}

	if err := app.Mailer.SendSMTPMessage(msg); err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error: false,
		Message: "Sent to "+ requestPayload.To,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}