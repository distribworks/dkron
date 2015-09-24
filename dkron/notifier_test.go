package dkron

import (
	"io"
	"net/http"
	"net/http/httptest"
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
	c := &Config{
		MailHost:     "mailtrap.io",
		MailPort:     2525,
		MailUsername: "45326e3b115066bbb",
		MailPassword: "7f496ed2b06688",
		MailFrom:     "dkron@dkron.io",
	}

	n := Notification(c, &Execution{
		Job: &Job{
			OwnerEmail: "victorcoder@gmail.com",
		},
	})

	n.Send()
}
