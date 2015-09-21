package dkron

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"net/textproto"
	"text/template"

	"github.com/jordan-wright/email"
)

type Notifier struct {
	Config *Config
	Job    *Job
}

func Notification(config *Config, job *Job) *Notifier {
	return &Notifier{
		Config: config,
		Job:    job,
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
		To:      []string{n.Job.OwnerEmail},
		From:    n.Config.MailFrom.String(),
		Subject: "[Dkron] Job execution report",
		Text:    []byte("Text Body is, of course, supported!"),
		Headers: textproto.MIMEHeader{},
	}

	e.Send(fmt.Sprintf("%s:%s", n.Config.MailHost, n.Config.MailPort), smtp.PlainAuth("", n.Config.MailUsername, n.Config.MailPassword, n.Config.MailHost))
}

func (n *Notifier) callExecutionWebhook() {
	var rep *bytes.Buffer
	t := template.Must(template.New("report").Parse(n.Config.WebhookPayload))

	data := struct {
		Report string
	}{
		"blahblah",
	}
	err := t.Execute(rep, data)
	if err != nil {
		log.Println("executing template:", err)
	}

	var jsonStr = []byte(n.Config.WebhookPayload)
	req, err := http.NewRequest("POST", n.Config.WebhookURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
