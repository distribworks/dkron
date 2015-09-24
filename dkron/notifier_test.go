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
