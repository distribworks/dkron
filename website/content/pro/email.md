---
title: Email processor
---

If you need special email notification rules for a job, use the Email processor.

Configuration of the email processor is stored in a file named `dkron-processor-email.yml` in the same locations as `dkron.yml`, and should include a list of providers, you could include any number of providers.

Example:
```yaml
provider1:
  host: smtp.myprovider.com
  port: 25
  username: myusername
  password: mypassword
  from: cron@mycompany.com
  subjectPrefix: '[Staging] '
```

Then you can configure each job with the following options:

Example:

```json
  "processors": {
    "email": {
      "provider": "provider1",
      "emails": "team@mycompany.com, owner@mycompany.com",
      "onSuccess": true
    }
  }
```

By default the email procesor doesn't send emails on job success, you should use the `onSuccess` parameter like in the previous example.
