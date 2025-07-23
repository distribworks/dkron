package dkron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"net/textproto"
	"net/url"
	"strings"
	"text/template"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/jordan-wright/email"
	"github.com/sirupsen/logrus"
)

// Notifier represents a new notification to be sent by any of the available notificators
type notifier struct {
	Config         *Config
	Job            *Job
	Execution      *Execution
	ExecutionGroup []*Execution

	logger *logrus.Entry
}

type MailPasswordResponse struct {
	Code int    `json:"code"`
	Password string `json:"password"`
	Expire string `json:"expire"`
}

// NewNotifier returns a new notifier
func SendPreNotifications(config *Config, execution *Execution, exGroup []*Execution, job *Job, logger *logrus.Entry) error {
	var werr error

	n := &notifier{
		logger: logger,

		Config:         config,
		Execution:      execution,
		ExecutionGroup: exGroup,
		Job:            job,
	}

	if err := n.cronitorTelemetry("run"); err != nil {
		werr = multierror.Append(werr, fmt.Errorf("notifier: error sending cronitor telemetry %w", err))
	}

	if n.Config.PreWebhookEndpoint != "" && n.Config.PreWebhookPayload != "" {
		if err := n.callPreExecutionWebhook(); err != nil {
			werr = multierror.Append(werr, fmt.Errorf("notifier: error sending email: %w", err))
		}
	}

	return werr
}

// Send sends the notifications using any configured method
func SendPostNotifications(config *Config, execution *Execution, exGroup []*Execution, job *Job, logger *logrus.Entry) error {
	n := &notifier{
		logger: logger,

		Config:         config,
		Execution:      execution,
		ExecutionGroup: exGroup,
		Job:            job,
	}

	var werr error

	if err := n.cronitorTelemetry("complete"); err != nil {
		werr = multierror.Append(werr, fmt.Errorf("notifier: error sending cronitor telemetry %w", err))
	}

	if n.Config.MailHost != "" && n.Config.MailPort != 0 && n.Job.OwnerEmail != "" {
		if err := n.sendExecutionEmail(); err != nil {
			werr = multierror.Append(werr, fmt.Errorf("notifier: error sending email: %w", err))
		}
	}

	if n.Config.WebhookEndpoint != "" && n.Config.WebhookPayload != "" {
		if err := n.callExecutionWebhook(); err != nil {
			werr = multierror.Append(werr, fmt.Errorf("notifier: error posting notification: %w", err))
		}
	}

	if n.Config.HubbleURL != "" && n.Config.HubbleToken != "" {
		if err := n.sendHubbleNotification(); err != nil {
			werr = multierror.Append(werr, fmt.Errorf("notifier: error sending Hubble notification: %w", err))
		}
	}

	return werr
}

func (n *notifier) report() string {
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

func (n *notifier) buildTemplate(templ string) *bytes.Buffer {
	t, e := template.New("report").Parse(templ)
	if e != nil {
		n.logger.WithError(e).Error("notifier: error parsing template")
		return bytes.NewBuffer([]byte("Failed to parse template: " + e.Error()))
	}

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
		n.logger.WithError(err).Error("notifier: error executing template")
		return bytes.NewBuffer([]byte("Failed to execute template:" + err.Error()))
	}
	return out
}

func (n *notifier) sendExecutionEmail() error {
	var data *bytes.Buffer
	if n.Config.MailPayload != "" {
		data = n.buildTemplate(n.Config.MailPayload)
	} else {
		data = bytes.NewBuffer([]byte(n.Execution.Output))
	}
	e := &email.Email{
		To:      []string{n.Job.OwnerEmail},
		From:    n.Config.MailFrom,
		Subject: fmt.Sprintf("%s%s %s execution report", n.Config.MailSubjectPrefix, n.statusString(), n.Execution.JobName),
		Text:    []byte(data.Bytes()),
		Headers: textproto.MIMEHeader{},
	}

	serverAddr := fmt.Sprintf("%s:%d", n.Config.MailHost, n.Config.MailPort)
	if err := e.Send(serverAddr, n.auth()); err != nil {
		return fmt.Errorf("notifier: Error sending email %s", err)
	}

	return nil
}

func (n *notifier) sendHubbleNotification() error {
	decodeData := make(map[string]string, len(n.Config.HubbleTemplate))
	for k, v := range n.Config.HubbleTemplate {
		decodeData[k] = v
	}
	if n.Execution.Success {
		decodeData["status"] = "OK"
	} else {
		decodeData["status"] = "PROBLEM"
	}
	jsonByte, err := json.Marshal(decodeData)
	if err != nil {
		return fmt.Errorf("notifier: Error marshalling Hubble template: %s", err)
	}
	decodeData["value"] = `{ "job_name":"` + n.Execution.JobName + `","node_name":"` + n.Execution.NodeName + `","status":"` + n.statusString() + `","output":"` + n.Execution.Output + `"}`
	fmt.Println("Hubble decodeData: ", decodeData)
	req, err := http.NewRequest("POST", n.Config.HubbleURL, bytes.NewBuffer(jsonByte))
	if err != nil {
		return fmt.Errorf("notifier: Error creating Hubble notification request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("usertoken", n.Config.HubbleToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("notifier: Error posting Hubble notification: %s", err)
	}
	defer resp.Body.Close()
	return nil
}

func (n *notifier) updateEmailPasswd() (string, error) {
	params := url.Values{}
	params.Add("token", n.Config.MailPasswordToken)
	params.Add("domainuser", n.Config.MailUsername)
	req, err := http.NewRequest("GET", n.Config.MailPasswordUrl, nil)
	if err != nil {
		return "", fmt.Errorf("notifier: Error creating mail token request: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.URL.RawQuery = params.Encode()
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("notifier: Error get mail token notification: %s", err)
	}
	defer resp.Body.Close()
	bodeBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("notifier: Error reading mail token response body: %s", err)
	}
	// fmt.Println("Response Body:", string(bodeBytes))
	var passwordData MailPasswordResponse
	err = json.Unmarshal(bodeBytes, &passwordData);
	if err != nil {
		return "", fmt.Errorf("notifier: Error unmarshalling  mail token response body: %s", err)
	}
	// fmt.Println("in body password: ", passwordData.Password)
	return passwordData.Password, nil
}


func (n *notifier) auth() smtp.Auth {
	var auth smtp.Auth

	if n.Config.MailUsername != "" && n.Config.MailPassword != "" {
		mailPassword, err := n.updateEmailPasswd()
		if err != nil {
			mailPassword = n.Config.MailPassword
		}
		// fmt.Println("password: ", mailPassword)
		auth = smtp.PlainAuth("", n.Config.MailUsername,  mailPassword, n.Config.MailHost)
	}

	return auth
}

func (n *notifier) callPreExecutionWebhook() error {
	out := n.buildTemplate(n.Config.PreWebhookPayload)
	req, err := http.NewRequest("POST", n.Config.PreWebhookEndpoint, out)
	if err != nil {
		return err
	}
	for _, h := range n.Config.PreWebhookHeaders {
		if h != "" {
			kv := strings.Split(h, ":")
			if strings.EqualFold(kv[0], "host") {
				req.Host = strings.TrimSpace(kv[1])
			} else {
				req.Header.Set(kv[0], strings.TrimSpace(kv[1]))
			}
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("notifier: Error posting notification: %s", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	n.logger.WithFields(logrus.Fields{
		"status": resp.Status,
		"header": resp.Header,
		"body":   string(body),
	}).Debug("notifier: Pre Webhook call response")

	return nil
}

func (n *notifier) callExecutionWebhook() error {
	out := n.buildTemplate(n.Config.WebhookPayload)
	req, err := http.NewRequest("POST", n.Config.WebhookEndpoint, out)
	if err != nil {
		return err
	}
	for _, h := range n.Config.WebhookHeaders {
		if h != "" {
			kv := strings.Split(h, ":")
			if strings.EqualFold(kv[0], "host") {
				req.Host = strings.TrimSpace(kv[1])
			} else {
				req.Header.Set(kv[0], strings.TrimSpace(kv[1]))
			}
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("notifier: Error posting notification: %s", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	n.logger.WithFields(logrus.Fields{
		"status": resp.Status,
		"header": resp.Header,
		"body":   string(body),
	}).Debug("notifier: Webhook call response")

	return nil
}

func (n *notifier) statusString() string {
	if n.Execution.Success {
		return "Success"
	}
	return "Failed"
}

// cronitorTelemetry is called when a job starts to notify cronitor
func (n *notifier) cronitorTelemetry(state string) error {
	if n.Config.CronitorEndpoint != "" {
		params := url.Values{}
		params.Add("host", n.Execution.NodeName)
		params.Add("message", "Job "+state+" by Dkron")
		params.Add("series", n.Execution.Key())

		if state == "complete" && !n.Execution.Success {
			state = "fail"
		}
		params.Add("state", state)

		_, err := http.Get(n.Config.CronitorEndpoint + "/" + n.Execution.JobName + "?" + params.Encode())
		if err != nil {
			return fmt.Errorf("notifier: Error sending telemetry to cronitor: %s", err)
		}
	}

	return nil
}
