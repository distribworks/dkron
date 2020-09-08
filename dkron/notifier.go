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
	"time"

	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
)

// Notifier represents a new notification to be sent by any of the available notificators
type Notifier struct {
	Config         *Config
	Job            *Job
	Execution      *Execution
	ExecutionGroup []*Execution
}

// Notification creates a new Notifier instance
func Notification(config *Config, execution *Execution, exGroup []*Execution, job *Job) *Notifier {
	return &Notifier{
		Config:         config,
		Execution:      execution,
		ExecutionGroup: exGroup,
		Job:            job,
	}
}

// Send sends the notifications using any configured method
func (n *Notifier) Send() error {
	if n.Config.MailHost != "" && n.Config.MailPort != 0 && n.Job.OwnerEmail != "" {
		return n.sendExecutionEmail()
	}
	if n.Config.WebhookURL != "" && n.Config.WebhookPayload != "" {
		return n.callExecutionWebhook()
	}

	return nil
}

func (n *Notifier) report() string {
	var exgStr string
	for _, ex := range n.ExecutionGroup {
		exgStr = fmt.Sprintf("%s\t[Node]: %s [Start]: %s [End]: %s [Success]: %t\n",
			exgStr,
			ex.NodeName,
			ex.StartedAt,
			ex.FinishedAt,
			ex.Success)
	}

	return fmt.Sprintf("Executed: %s\nReporting node: %s\nStart time: %s\nEnd time: %s\nSuccess: %t\nNode: %s\nOutput: %s\nExecution group: %d\n%s",
		n.Execution.JobName,
		n.Config.NodeName,
		n.Execution.StartedAt,
		n.Execution.FinishedAt,
		n.Execution.Success,
		n.Execution.NodeName,
		n.Execution.Output,
		n.Execution.Group,
		exgStr)
}

func (n *Notifier) buildTemplate(templ string) *bytes.Buffer {
	t := template.Must(template.New("report").Parse(templ))

	data := struct {
		Report        string
		JobName       string
		ReportingNode string
		StartTime     time.Time
		FinishedAt    time.Time
		Success       string
		NodeName      string
		Output        string
	}{
		n.report(),
		n.Execution.JobName,
		n.Config.NodeName,
		n.Execution.StartedAt,
		n.Execution.FinishedAt,
		fmt.Sprintf("%t", n.Execution.Success),
		n.Execution.NodeName,
		n.Execution.Output,
	}

	out := &bytes.Buffer{}
	err := t.Execute(out, data)
	if err != nil {
		log.WithError(err).Error("notifier: error executing template")
		return bytes.NewBuffer([]byte("Failed to execute template:" + err.Error()))
	}
	return out
}

func (n *Notifier) sendExecutionEmail() error {
	var data *bytes.Buffer
	if n.Config.MailPayload != "" {
		data = n.buildTemplate(n.Config.MailPayload)
	} else {
		data = bytes.NewBuffer([]byte(n.Execution.Output))
	}
	e := &email.Email{
		To:      []string{n.Job.OwnerEmail},
		From:    n.Config.MailFrom,
		Subject: fmt.Sprintf("%s%s %s execution report", n.Config.MailSubjectPrefix, n.statusString(n.Execution), n.Execution.JobName),
		Text:    []byte(data.Bytes()),
		Headers: textproto.MIMEHeader{},
	}

	serverAddr := fmt.Sprintf("%s:%d", n.Config.MailHost, n.Config.MailPort)
	if err := e.Send(serverAddr, n.auth()); err != nil {
		return fmt.Errorf("notifier: Error sending email %s", err)
	}

	return nil
}

func (n *Notifier) auth() smtp.Auth {
	var auth smtp.Auth

	if n.Config.MailUsername != "" && n.Config.MailPassword != "" {
		auth = smtp.PlainAuth("", n.Config.MailUsername, n.Config.MailPassword, n.Config.MailHost)
	}

	return auth
}

func (n *Notifier) callExecutionWebhook() error {
	out := n.buildTemplate(n.Config.WebhookPayload)
	req, err := http.NewRequest("POST", n.Config.WebhookURL, out)
	if err != nil {
		return err
	}
	for _, h := range n.Config.WebhookHeaders {
		if h != "" {
			kv := strings.Split(h, ":")
			req.Header.Set(kv[0], strings.TrimSpace(kv[1]))
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("notifier: Error posting notification: %s", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.WithFields(logrus.Fields{
		"status": resp.Status,
		"header": resp.Header,
		"body":   string(body),
	}).Debug("notifier: Webhook call response")

	return nil
}

func (n *Notifier) statusString(execution *Execution) string {
	if execution.Success {
		return "Success"
	}
	return "Failed"
}
