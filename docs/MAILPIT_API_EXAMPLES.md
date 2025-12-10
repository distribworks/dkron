# Mailpit API Usage Examples

This document provides practical examples of using Mailpit's REST API in tests to verify email delivery.

## Overview

Mailpit provides a comprehensive REST API at `http://localhost:8025/api/v1/` that allows programmatic access to captured emails. This is essential for automated testing in CI/CD environments where the web UI is not accessible.

## Base URL

```
http://localhost:8025/api/v1/
```

## API Response Types

### Message Summary

```go
type mailpitMessage struct {
    ID      string           `json:"ID"`
    From    mailpitAddress   `json:"From"`
    To      []mailpitAddress `json:"To"`
    Subject string           `json:"Subject"`
    Snippet string           `json:"Snippet"`
    Created time.Time        `json:"Created"`
}

type mailpitAddress struct {
    Name    string `json:"Name"`
    Address string `json:"Address"`
}

type mailpitMessagesResponse struct {
    Total    int              `json:"total"`
    Messages []mailpitMessage `json:"messages"`
}
```

### Message Detail

```go
type mailpitMessageDetail struct {
    ID      string           `json:"ID"`
    From    mailpitAddress   `json:"From"`
    To      []mailpitAddress `json:"To"`
    Cc      []mailpitAddress `json:"Cc"`
    Bcc     []mailpitAddress `json:"Bcc"`
    Subject string           `json:"Subject"`
    Text    string           `json:"Text"`
    HTML    string           `json:"HTML"`
    Created time.Time        `json:"Created"`
}
```

## Common Operations

### 1. Get All Messages

Retrieve a list of all captured messages.

```go
func getMailpitMessages(t *testing.T, apiURL string) []mailpitMessage {
    t.Helper()

    resp, err := http.Get(apiURL + "/api/v1/messages")
    require.NoError(t, err, "Failed to call Mailpit API")
    defer resp.Body.Close()

    require.Equal(t, http.StatusOK, resp.StatusCode)

    var result mailpitMessagesResponse
    err = json.NewDecoder(resp.Body).Decode(&result)
    require.NoError(t, err, "Failed to decode response")

    return result.Messages
}
```

**Usage:**
```go
messages := getMailpitMessages(t, "http://localhost:8025")
assert.NotEmpty(t, messages, "Expected at least one message")
```

### 2. Get Message Details

Retrieve full details of a specific message including content.

```go
func getMailpitMessage(t *testing.T, apiURL, messageID string) *mailpitMessageDetail {
    t.Helper()

    url := fmt.Sprintf("%s/api/v1/message/%s", apiURL, messageID)
    resp, err := http.Get(url)
    require.NoError(t, err, "Failed to get message detail")
    defer resp.Body.Close()

    require.Equal(t, http.StatusOK, resp.StatusCode)

    var result mailpitMessageDetail
    err = json.NewDecoder(resp.Body).Decode(&result)
    require.NoError(t, err, "Failed to decode message detail")

    return &result
}
```

**Usage:**
```go
detail := getMailpitMessage(t, "http://localhost:8025", messageID)
assert.Contains(t, detail.Text, "expected content")
```

### 3. Search Messages

Search for messages by query string.

```go
func searchMailpitMessages(t *testing.T, apiURL, query string) []mailpitMessage {
    t.Helper()

    url := fmt.Sprintf("%s/api/v1/search?query=%s", apiURL, url.QueryEscape(query))
    resp, err := http.Get(url)
    require.NoError(t, err, "Failed to search messages")
    defer resp.Body.Close()

    require.Equal(t, http.StatusOK, resp.StatusCode)

    var result mailpitMessagesResponse
    err = json.NewDecoder(resp.Body).Decode(&result)
    require.NoError(t, err, "Failed to decode search results")

    return result.Messages
}
```

**Usage:**
```go
// Search by subject
messages := searchMailpitMessages(t, "http://localhost:8025", "subject:test")

// Search by recipient
messages := searchMailpitMessages(t, "http://localhost:8025", "to:user@example.com")
```

### 4. Find Specific Email

Find a specific email by subject or other criteria.

```go
func findEmailBySubject(t *testing.T, apiURL, subjectContains string) *mailpitMessage {
    t.Helper()

    messages := getMailpitMessages(t, apiURL)
    
    for i := range messages {
        if strings.Contains(messages[i].Subject, subjectContains) {
            return &messages[i]
        }
    }
    
    return nil
}
```

**Usage:**
```go
email := findEmailBySubject(t, "http://localhost:8025", "[Test]")
require.NotNil(t, email, "Email with subject containing '[Test]' not found")
```

### 5. Delete All Messages

Clean up all messages (useful for test isolation).

```go
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
```

**Usage:**
```go
// Clean up before test
deleteAllMailpitMessages(t, "http://localhost:8025")

// Run your test
sendTestEmail()

// Clean up after test (optional - can use defer)
defer deleteAllMailpitMessages(t, "http://localhost:8025")
```

### 6. Wait for Email

Wait for an email to arrive (with timeout).

```go
func waitForEmail(t *testing.T, apiURL string, timeout time.Duration) *mailpitMessage {
    t.Helper()

    deadline := time.Now().Add(timeout)
    
    for time.Now().Before(deadline) {
        messages := getMailpitMessages(t, apiURL)
        if len(messages) > 0 {
            return &messages[0]
        }
        time.Sleep(100 * time.Millisecond)
    }
    
    t.Fatal("Timeout waiting for email")
    return nil
}
```

**Usage:**
```go
// Send email
sendTestEmail()

// Wait up to 5 seconds for email to arrive
email := waitForEmail(t, "http://localhost:8025", 5*time.Second)
assert.NotNil(t, email)
```

## Complete Test Example

Here's a complete example showing best practices:

```go
func TestEmailNotification(t *testing.T) {
    // Setup
    mailpitAPIURL := "http://localhost:8025"
    
    // Check if Mailpit is available
    if !isMailpitAvailable(t, mailpitAPIURL) {
        t.Skip("Mailpit not available")
    }
    
    // Clean up before test
    deleteAllMailpitMessages(t, mailpitAPIURL)
    
    // Perform action that sends email
    err := sendNotificationEmail("test@example.com", "Test Subject", "Test Body")
    require.NoError(t, err, "Failed to send email")
    
    // Give Mailpit time to process
    time.Sleep(100 * time.Millisecond)
    
    // Verify email was received
    messages := getMailpitMessages(t, mailpitAPIURL)
    require.NotEmpty(t, messages, "No emails received")
    
    // Find our test email
    var testEmail *mailpitMessage
    for i := range messages {
        if strings.Contains(messages[i].Subject, "Test Subject") {
            testEmail = &messages[i]
            break
        }
    }
    require.NotNil(t, testEmail, "Test email not found")
    
    // Verify metadata
    assert.Equal(t, "test@example.com", testEmail.To[0].Address)
    assert.Contains(t, testEmail.Subject, "Test Subject")
    
    // Get full message details
    detail := getMailpitMessage(t, mailpitAPIURL, testEmail.ID)
    require.NotNil(t, detail)
    
    // Verify content
    assert.Contains(t, detail.Text, "Test Body")
    
    // Log for debugging
    t.Logf("Email verified - ID: %s, Subject: %s", detail.ID, detail.Subject)
    
    // Cleanup (optional)
    deleteAllMailpitMessages(t, mailpitAPIURL)
}
```

## Advanced Examples

### Verify Multiple Recipients

```go
func TestMultipleRecipients(t *testing.T) {
    mailpitAPIURL := "http://localhost:8025"
    deleteAllMailpitMessages(t, mailpitAPIURL)
    
    // Send email to multiple recipients
    err := sendEmail([]string{"user1@example.com", "user2@example.com"}, "Subject", "Body")
    require.NoError(t, err)
    
    time.Sleep(100 * time.Millisecond)
    
    messages := getMailpitMessages(t, mailpitAPIURL)
    require.Len(t, messages, 1, "Expected exactly one email")
    
    detail := getMailpitMessage(t, mailpitAPIURL, messages[0].ID)
    
    // Verify all recipients
    assert.Len(t, detail.To, 2)
    recipients := []string{detail.To[0].Address, detail.To[1].Address}
    assert.Contains(t, recipients, "user1@example.com")
    assert.Contains(t, recipients, "user2@example.com")
}
```

### Verify HTML Email

```go
func TestHTMLEmail(t *testing.T) {
    mailpitAPIURL := "http://localhost:8025"
    deleteAllMailpitMessages(t, mailpitAPIURL)
    
    // Send HTML email
    err := sendHTMLEmail("user@example.com", "Subject", "<h1>Hello</h1><p>World</p>")
    require.NoError(t, err)
    
    time.Sleep(100 * time.Millisecond)
    
    messages := getMailpitMessages(t, mailpitAPIURL)
    require.NotEmpty(t, messages)
    
    detail := getMailpitMessage(t, mailpitAPIURL, messages[0].ID)
    
    // Verify HTML content
    assert.Contains(t, detail.HTML, "<h1>Hello</h1>")
    assert.Contains(t, detail.HTML, "<p>World</p>")
    
    // Verify text fallback exists
    assert.NotEmpty(t, detail.Text)
}
```

### Verify Email Ordering

```go
func TestEmailOrdering(t *testing.T) {
    mailpitAPIURL := "http://localhost:8025"
    deleteAllMailpitMessages(t, mailpitAPIURL)
    
    // Send multiple emails
    sendEmail("user@example.com", "First Email", "Body 1")
    time.Sleep(10 * time.Millisecond)
    sendEmail("user@example.com", "Second Email", "Body 2")
    time.Sleep(100 * time.Millisecond)
    
    messages := getMailpitMessages(t, mailpitAPIURL)
    require.Len(t, messages, 2)
    
    // Messages are typically returned newest first
    assert.Contains(t, messages[0].Subject, "Second Email")
    assert.Contains(t, messages[1].Subject, "First Email")
}
```

### Test Email Not Sent

```go
func TestEmailNotSent(t *testing.T) {
    mailpitAPIURL := "http://localhost:8025"
    deleteAllMailpitMessages(t, mailpitAPIURL)
    
    // Perform action that should NOT send email
    err := performActionWithoutEmail()
    require.NoError(t, err)
    
    time.Sleep(100 * time.Millisecond)
    
    // Verify no emails were sent
    messages := getMailpitMessages(t, mailpitAPIURL)
    assert.Empty(t, messages, "Expected no emails to be sent")
}
```

## Helper Utilities

### Check Mailpit Availability

```go
func isMailpitAvailable(t *testing.T, apiURL string) bool {
    t.Helper()
    
    client := &http.Client{Timeout: 2 * time.Second}
    resp, err := client.Get(apiURL + "/api/v1/messages")
    if err != nil {
        return false
    }
    defer resp.Body.Close()
    
    return resp.StatusCode == http.StatusOK
}
```

### Count Messages

```go
func countMessages(t *testing.T, apiURL string) int {
    t.Helper()
    
    messages := getMailpitMessages(t, apiURL)
    return len(messages)
}
```

### Assert Email Count

```go
func assertEmailCount(t *testing.T, apiURL string, expected int) {
    t.Helper()
    
    messages := getMailpitMessages(t, apiURL)
    assert.Len(t, messages, expected, "Unexpected number of emails")
}
```

## Best Practices

### 1. Always Clean Up Before Tests

```go
// Clean up at the start of each test
deleteAllMailpitMessages(t, mailpitAPIURL)
```

This ensures test isolation and prevents false positives from previous test runs.

### 2. Add Small Delays

```go
// Send email
sendEmail()

// Give Mailpit time to process
time.Sleep(100 * time.Millisecond)

// Now verify
messages := getMailpitMessages(t, mailpitAPIURL)
```

Mailpit is fast, but a small delay ensures the email is fully processed.

### 3. Use Helper Functions

Create reusable helper functions for common operations:

```go
func requireEmailWithSubject(t *testing.T, apiURL, subject string) *mailpitMessageDetail {
    t.Helper()
    
    email := findEmailBySubject(t, apiURL, subject)
    require.NotNil(t, email, "Email with subject '%s' not found", subject)
    
    return getMailpitMessage(t, apiURL, email.ID)
}
```

### 4. Log Useful Information

```go
t.Logf("Email verified:")
t.Logf("  ID: %s", detail.ID)
t.Logf("  From: %s", detail.From.Address)
t.Logf("  To: %s", detail.To[0].Address)
t.Logf("  Subject: %s", detail.Subject)
```

This helps with debugging when tests fail.

### 5. Check API Availability

```go
if !isMailpitAvailable(t, mailpitAPIURL) {
    t.Skip("Mailpit not available. Start with: docker run -p 8025:8025 -p 1025:1025 axllent/mailpit")
}
```

Skip tests gracefully when Mailpit is not running locally.

## Troubleshooting

### API Returns 404

- Verify Mailpit is running: `docker ps | grep mailpit`
- Check the correct port (8025 for API, not 1025 for SMTP)
- Ensure URL includes `/api/v1/` path

### No Messages Found

- Add a delay after sending email: `time.Sleep(100 * time.Millisecond)`
- Verify email was actually sent (check for send errors)
- Check Mailpit logs: `docker logs <container-id>`
- Ensure you're not checking before the email arrives

### Message ID Invalid

- Don't reuse message IDs between tests
- Clean up messages at the start of tests
- Message IDs are only valid until messages are deleted

## Resources

- [Mailpit API Documentation](https://mailpit.axllent.org/docs/api-v1/)
- [Real Implementation Example](../dkron/notifier_test.go)
- [Email Testing Guide](EMAIL_TESTING.md)
- [GitHub Actions Integration](GITHUB_ACTIONS_MAILPIT.md)

## See Also

- `dkron/notifier_test.go` - Complete working example in the Dkron codebase
- [Mailpit GitHub](https://github.com/axllent/mailpit) - Official repository
- [REST API Reference](https://mailpit.axllent.org/docs/api-v1/) - Complete API documentation