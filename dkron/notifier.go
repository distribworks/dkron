package dkron

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"net/textproto"
	"strings"
	"text/template"

	"github.com/Sirupsen/logrus"
	"github.com/jordan-wright/email"
)

type Notifier struct {
	Config         *Config
	Execution      *Execution
	ExecutionGroup []*Execution
}

func Notification(config *Config, execution *Execution, exGroup []*Execution) *Notifier {
	return &Notifier{
		Config:         config,
		Execution:      execution,
		ExecutionGroup: exGroup,
	}
}

func (n *Notifier) Send() {
	if n.Config.MailHost != "" && n.Config.MailPort != 0 && n.Execution.Job.OwnerEmail != "" {
		n.sendExecutionEmail()
	}
	if n.Config.WebhookURL != "" && n.Config.WebhookPayload != "" {
		n.callExecutionWebhook()
	}
}

func (n *Notifier) report() string {
	t := template.Must(template.New("report").Parse(n.Config.WebhookPayload))

	data := struct {
		Report string
	}{
		fmt.Sprintf("Executed: %s\nStart time: %s\nEnd time: %s\nSuccess: %t",
			n.Execution.JobName,
			n.Execution.StartedAt,
			n.Execution.FinishedAt,
			n.Execution.Success),
	}

	out := &bytes.Buffer{}
	err := t.Execute(out, data)
	if err != nil {
		log.Error("executing template:", err)
	}

	return out.String()
}

func (n *Notifier) sendExecutionEmail() {
	e := &email.Email{
		To:      []string{n.Execution.Job.OwnerEmail},
		From:    n.Config.MailFrom,
		Subject: "[Dkron] Job execution report",
		Text:    []byte(n.report()),
		Headers: textproto.MIMEHeader{},
	}

	serverAddr := fmt.Sprintf("%s:%d", n.Config.MailHost, n.Config.MailPort)
	err := e.Send(serverAddr, smtp.PlainAuth("", n.Config.MailUsername, n.Config.MailPassword, n.Config.MailHost))
	if err != nil {
		log.Fatal(err)
	}
}

func (n *Notifier) callExecutionWebhook() {
	req, err := http.NewRequest("POST", n.Config.WebhookURL, bytes.NewBufferString(n.report()))
	for _, h := range n.Config.WebhookHeaders {
		if h != "" {
			kv := strings.Split(h, ":")
			req.Header.Set(kv[0], strings.TrimSpace(kv[1]))
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.WithFields(logrus.Fields{
		"status": resp.Status,
		"header": resp.Header,
		"body":   string(body),
	}).Debug("Webhook call response")
}
