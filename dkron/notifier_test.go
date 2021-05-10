package dkron

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

	n := Notification(c, &Execution{}, []*Execution{}, &Job{})
	log := getTestLogger()
	assert.NoError(t, n.Send(log))
}

func TestNotifier_sendExecutionEmail(t *testing.T) {
	c := &Config{
		MailHost:          "smtp.mailtrap.io",
		MailPort:          2525,
		MailUsername:      "45326e3b115066bbb",
		MailPassword:      "7f496ed2b06688",
		MailFrom:          "dkron@dkron.io",
		MailSubjectPrefix: "[Test]",
	}

	job := &Job{
		OwnerEmail: "cron@job.com",
	}

	ex1 := &Execution{
		JobName:    "test",
		StartedAt:  time.Now(),
		FinishedAt: time.Now(),
		Success:    true,
		NodeName:   "test-node",
		Output:     "test-output",
	}

	exg := []*Execution{
		{
			JobName:   "test",
			StartedAt: time.Now(),
			NodeName:  "test-node2",
			Output:    "test-output",
		},
		ex1,
	}

	log := getTestLogger()
	n := Notification(c, ex1, exg, job)
	assert.NoError(t, n.Send(log))
}

func Test_auth(t *testing.T) {
	n1 := &Notifier{
		Config: &Config{
			MailHost:     "localhost",
			MailPort:     25,
			MailUsername: "username",
			MailPassword: "password",
		}}

	a1 := n1.auth()
	assert.NotNil(t, a1)

	n2 := &Notifier{
		Config: &Config{
			MailHost: "localhost",
			MailPort: 25,
		}}

	a2 := n2.auth()
	assert.Nil(t, a2)
}

func TestNotifier_buildTemplate(t *testing.T) {
	c := &Config{
		NodeName: "test-node",
	}

	ex1 := &Execution{
		JobName:    "test",
		StartedAt:  time.Now(),
		FinishedAt: time.Now(),
		Success:    true,
		NodeName:   "test-node",
		Output:     "test-output",
	}

	exg := []*Execution{
		{
			JobName:   "test",
			StartedAt: time.Now(),
			NodeName:  "test-node2",
			Output:    "test-output",
		},
		ex1,
	}

	log := getTestLogger()
	n := Notification(c, ex1, exg, nil)
	for _, tc := range templateTestCases(n) {
		got := n.buildTemplate(tc.template, log).String()

		if tc.exp != got {
			t.Errorf("Exp: %s\nGot: %s", tc.exp, got)
		}
	}
}

type templateTestCase struct {
	desc     string
	exp      string
	template string
}

var templateTestCases = func(n *Notifier) []templateTestCase {
	return []templateTestCase{
		{
			desc:     "Report template variable",
			exp:      n.report(),
			template: "{{.Report}}",
		},
		{
			desc:     "JobName template variable",
			exp:      n.Execution.JobName,
			template: "{{.JobName}}",
		},
		{
			desc:     "ReportingNode template variable",
			exp:      n.Config.NodeName,
			template: "{{.ReportingNode}}",
		},
		{
			desc:     "StartTime template variable",
			exp:      n.Execution.StartedAt.String(),
			template: "{{.StartTime}}",
		},
		{
			desc:     "FinishedAt template variable",
			exp:      n.Execution.FinishedAt.String(),
			template: "{{.FinishedAt}}",
		},
		{
			desc:     "Success template variable",
			exp:      fmt.Sprintf("%t", n.Execution.Success),
			template: "{{.Success}}",
		},
		{
			desc:     "NodeName template variable",
			exp:      n.Execution.NodeName,
			template: "{{.NodeName}}",
		},
		{
			desc:     "Output template variable",
			exp:      n.Execution.Output,
			template: "{{.Output}}",
		},
	}
}
