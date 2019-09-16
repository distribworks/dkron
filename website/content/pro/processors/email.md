---
title: Email processor
---

The Email processor provides flexibility to job email notifications.

Configuration of the email processor is stored in a file named `dkron-processor-email.yml` in the same locations as `dkron.yml`, and should include a list of providers, it can include any number of providers.

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

Then configure each job with the following options:

Example:

```json
{
  "processors": {
    "email": {
      "provider": "provider1",
      "emails": "team@mycompany.com, owner@mycompany.com",
      "onSuccess": "true"
    }
  }
}
```

By default the email procesor doesn't send emails on job success, the `onSuccess` parameter, enables it, like in the previous example.
