package dkron

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// checkMailpitAvailable verifies that Mailpit is running and accessible.
// This is useful for providing clear error messages in tests.
func checkMailpitAvailable(t *testing.T, host string, port int) {
	t.Helper()
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		t.Skipf("Mailpit is not available at %s. Start it with: docker run -p 8025:8025 -p 1025:1025 axllent/mailpit", address)
		return
	}
	conn.Close()
}

// mailpitMessage represents a message in the Mailpit API response
type mailpitMessage struct {
	ID      string           `json:"ID"`
	From    mailpitAddress   `json:"From"`
	To      []mailpitAddress `json:"To"`
	Subject string           `json:"Subject"`
	Snippet string           `json:"Snippet"`
	Created time.Time        `json:"Created"`
}

// mailpitAddress represents an email address in the Mailpit API
type mailpitAddress struct {
	Name    string `json:"Name"`
	Address string `json:"Address"`
}

// mailpitMessagesResponse represents the response from the Mailpit messages API
type mailpitMessagesResponse struct {
	Total    int              `json:"total"`
	Messages []mailpitMessage `json:"messages"`
}

// mailpitMessageDetail represents a detailed message from the Mailpit API
type mailpitMessageDetail struct {
	ID      string           `json:"ID"`
	From    mailpitAddress   `json:"From"`
	To      []mailpitAddress `json:"To"`
	Subject string           `json:"Subject"`
	Text    string           `json:"Text"`
	HTML    string           `json:"HTML"`
	Created time.Time        `json:"Created"`
}

// getMailpitMessages retrieves all messages from Mailpit API
func getMailpitMessages(t *testing.T, apiURL string) []mailpitMessage {
	t.Helper()

	resp, err := http.Get(apiURL + "/api/v1/messages")
	require.NoError(t, err, "Failed to call Mailpit API")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Mailpit API returned non-200 status")

	var result mailpitMessagesResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Failed to decode Mailpit API response")

	return result.Messages
}

// getMailpitMessage retrieves a specific message detail from Mailpit API
func getMailpitMessage(t *testing.T, apiURL, messageID string) *mailpitMessageDetail {
	t.Helper()

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/message/%s", apiURL, messageID))
	require.NoError(t, err, "Failed to call Mailpit API for message detail")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Mailpit API returned non-200 status for message detail")

	var result mailpitMessageDetail
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Failed to decode Mailpit message detail response")

	return &result
}

// deleteAllMailpitMessages deletes all messages from Mailpit
func deleteAllMailpitMessages(t *testing.T, apiURL string) {
	t.Helper()

	req, err := http.NewRequest("DELETE", apiURL+"/api/v1/messages", nil)
	require.NoError(t, err, "Failed to create delete request")

	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}
	// Don't fail if delete fails - it's just cleanup
}

func TestNotifier_callExecutionWebhook(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(w, r.Body)
	}))
	defer ts.Close()

	c := &Config{
		WebhookEndpoint: ts.URL,
		WebhookPayload:  `payload={"text": "{{.Report}}"}`,
		WebhookHeaders:  []string{"Content-Type: application/x-www-form-urlencoded"},
	}

	log := getTestLogger()
	err := SendPostNotifications(c, &Execution{}, []*Execution{}, &Job{}, log)
	assert.NoError(t, err)
}

func TestNotifier_callExecutionWebhookHostHeader(t *testing.T) {
	var got string
	var exp = "dkron.io"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(w, r.Body)
		got = r.Host
	}))
	defer ts.Close()

	c := &Config{
		WebhookEndpoint: ts.URL,
		WebhookPayload:  `payload={"text": "{{.Report}}"}`,
		WebhookHeaders:  []string{"Content-Type: application/x-www-form-urlencoded", fmt.Sprintf("Host: %s", exp)},
	}

	log := getTestLogger()
	err := SendPostNotifications(c, &Execution{}, []*Execution{}, &Job{}, log)
	assert.NoError(t, err)

	if exp != got {
		t.Errorf("Exp: %s\nGot: %s", exp, got)
	}
}

func TestNotifier_sendExecutionEmail(t *testing.T) {
	// This test requires Mailpit to be running for email testing.
	// Mailpit is a local SMTP server that captures emails without sending them.
	// Start Mailpit with: docker run -p 8025:8025 -p 1025:1025 axllent/mailpit
	// View captured emails at: http://localhost:8025
	// See docs/EMAIL_TESTING.md for more information.
	// In GitHub Actions, Mailpit runs automatically as a service container.

	mailHost := "localhost"
	mailPort := 1025
	mailpitAPIURL := "http://localhost:8025"

	// Check if Mailpit is available (will skip test if not, except in CI)
	checkMailpitAvailable(t, mailHost, mailPort)

	// Clean up any existing messages before the test
	deleteAllMailpitMessages(t, mailpitAPIURL)

	c := &Config{
		MailHost:          mailHost,
		MailPort:          uint16(mailPort), // Mailpit SMTP port
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
	err := SendPostNotifications(c, ex1, exg, job, log)
	require.NoError(t, err, "Failed to send email notification")

	// Give Mailpit a moment to process the email
	time.Sleep(100 * time.Millisecond)

	// Verify email was received using Mailpit API
	messages := getMailpitMessages(t, mailpitAPIURL)
	require.NotEmpty(t, messages, "No emails were captured by Mailpit")

	// Find our test email
	var testEmail *mailpitMessage
	for i := range messages {
		if strings.Contains(messages[i].Subject, "[Test]") &&
			strings.Contains(messages[i].Subject, "test execution report") {
			testEmail = &messages[i]
			break
		}
	}
	require.NotNil(t, testEmail, "Test email not found in Mailpit")

	// Verify email metadata
	assert.Equal(t, "dkron@dkron.io", testEmail.From.Address, "From address mismatch")
	assert.Equal(t, "cron@job.com", testEmail.To[0].Address, "To address mismatch")
	assert.Contains(t, testEmail.Subject, "[Test]", "Subject should contain prefix")
	assert.Contains(t, testEmail.Subject, "Success", "Subject should indicate success")
	assert.Contains(t, testEmail.Subject, "test", "Subject should contain job name")

	// Get full message details
	messageDetail := getMailpitMessage(t, mailpitAPIURL, testEmail.ID)
	require.NotNil(t, messageDetail, "Failed to get message detail")

	// Verify email content
	// Note: Without a custom MailPayload, the email body only contains the execution output
	assert.Contains(t, messageDetail.Text, "test-output", "Email body should contain execution output")

	t.Logf("Successfully verified email via Mailpit API:")
	t.Logf("  Subject: %s", messageDetail.Subject)
	t.Logf("  From: %s", messageDetail.From.Address)
	t.Logf("  To: %s", messageDetail.To[0].Address)
	t.Logf("  Message ID: %s", messageDetail.ID)
}

func Test_auth(t *testing.T) {
	n1 := &notifier{
		Config: &Config{
			MailHost:     "localhost",
			MailPort:     25,
			MailUsername: "username",
			MailPassword: "password",
		}}

	a1 := n1.auth()
	assert.NotNil(t, a1)

	n2 := &notifier{
		Config: &Config{
			MailHost: "localhost",
			MailPort: 25,
		}}

	a2 := n2.auth()
	assert.Nil(t, a2)
}

func TestNotifier_buildTemplate(t *testing.T) {
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
			Output:    "test-output\r\ntest-node",
		},
		ex1,
	}

	log := getTestLogger()
	n := &notifier{
		Config: &Config{
			NodeName: "test-node",
		},
		Execution:      ex1,
		ExecutionGroup: exg,
		logger:         log,
	}
	for _, tc := range templateTestCases(n) {
		got := n.buildTemplate(tc.template).String()

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

var templateTestCases = func(n *notifier) []templateTestCase {
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
