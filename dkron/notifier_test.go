package dkron

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func TestNotifier_callExecutionWebhook(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, r.Body)
	}))
	defer ts.Close()

	c := &Config{
		WebhookURL:     ts.URL,
		WebhookPayload: `payload={"text": "{{.Report}}"}`,
		WebhookHeaders: []string{"Content-Type: application/x-www-form-urlencoded"},
	}

	n := Notification(c, &Execution{})

	n.Send()
}

func TestNotifier_sendExecutionEmail(t *testing.T) {
	port, _ := strconv.ParseInt(os.Getenv("DKRON_MAIL_PORT"), 10, 16)

	c := &Config{
		MailHost:     os.Getenv("DKRON_MAIL_HOST"),
		MailPort:     uint16(port),
		MailUsername: os.Getenv("DKRON_MAIL_USERNAME"),
		MailPassword: os.Getenv("DKRON_MAIL_PASSWORD"),
		MailFrom:     os.Getenv("DKRON_MAIL_FROM"),
	}

	n := Notification(c, &Execution{
		Job: &Job{
			OwnerEmail: "victorcoder@gmail.com",
		},
	})

	n.Send()
}
