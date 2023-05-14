---
title: Slack processor
---

The Slack processor provides slack notifications with multiple configurations and rich format.

Configuration of the slack processor is stored in a file named `dkron-processor-slack.yml` in the same locations as `dkron.yml`, and should include a list of teams, it can include any number of teams.

![](/img/slack.png)

Example:
```yaml
team1:
  webhook_url: https://hooks.slack.com/services/XXXXXXXXXXXXXXXXXXX
  bot_name: Dkron Production
```

Then configure each job with the following options:

Example:

```json
{
  "processors": {
    "slack": {
      "team": "team1",
      "channel": "#cron-production",
      "onSuccess": true
    }
  }
}
```

By default the slack procesor doesn't send notifications on job success, the `onSuccess` parameter, enables it, like in the previous example.
