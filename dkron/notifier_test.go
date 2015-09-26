package dkron

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

	n := Notification(c, &Execution{}, []*Execution{})

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

	job := &Job{
		OwnerEmail: "victorcoder@gmail.com",
	}

	ex1 := &Execution{
		JobName:    "test",
		StartedAt:  time.Now(),
		FinishedAt: time.Now(),
		Success:    true,
		Job:        job,
		NodeName:   "test-node",
		Output:     []byte("test-output"),
	}

	exg := []*Execution{
		&Execution{
			JobName:   "test",
			StartedAt: time.Now(),
			Job:       job,
			NodeName:  "test-node2",
			Output:    []byte("test-output"),
		},
		ex1,
	}

	n := Notification(c, ex1, exg)

	n.Send()
}
