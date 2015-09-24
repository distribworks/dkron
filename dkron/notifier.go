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
	Config    *Config
	Execution *Execution
}

func Notification(config *Config, execution *Execution) *Notifier {
	return &Notifier{
		Config:    config,
		Execution: execution,
	}
}

func (n *Notifier) Send() {
	if n.Config.MailFrom.String() != "" && n.Config.MailHost != "" && n.Config.MailPort != 0 {
		n.sendExecutionEmail()
	}
	if n.Config.WebhookURL != "" && n.Config.WebhookPayload != "" {
		n.callExecutionWebhook()
	}
}

func (n *Notifier) sendExecutionEmail() {
	e := &email.Email{
		To:      []string{n.Execution.Job.OwnerEmail},
		From:    n.Config.MailFrom.String(),
		Subject: "[Dkron] Job execution report",
		Text:    []byte("Text Body is, of course, supported!"),
		Headers: textproto.MIMEHeader{},
	}

	e.Send(fmt.Sprintf("%s:%d", n.Config.MailHost, n.Config.MailPort), smtp.PlainAuth("", n.Config.MailUsername, n.Config.MailPassword, n.Config.MailHost))
}

func (n *Notifier) callExecutionWebhook() {
	t := template.Must(template.New("report").Parse(n.Config.WebhookPayload))

	data := struct {
		Report string
	}{
		"This is the report",
	}

	out := &bytes.Buffer{}
	err := t.Execute(out, data)
	if err != nil {
		log.Error("executing template:", err)
	}

	req, err := http.NewRequest("POST", n.Config.WebhookURL, out)
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
