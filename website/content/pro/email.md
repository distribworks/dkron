---
title: Email processor
---

If you need special email notification rules for a job, use the Email processor.

Example:

```json
  "processors": {
    "email": {
      "host": "mail.google.com",
      "username": "foo",
      "password": "bar",
      "emails": "team@example.com",
      "from": "cron@example.com",
      "subject_prefix": "[Dkron]"
    }
  }
```
