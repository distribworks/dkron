# Cronitor Integration

:::info
This feature is available since v3.2.0
:::

Dkron includes tight integration with [Cronitor](https://cronitor.io/) for advanced execution monitoring, check their product page for more information.

To setup the integration add this parameter to your `dkron.yaml` config file, including your private cronitor API endpoint key:

:::caution
Remember: do not check your config file into your source code repository, it contains sensitive information.
:::

```
cronitor-endpoint: https://cronitor.link/p/xxxxxxxxxxxxxxxxxxxxxx
```

You can also use the `--cronitor-endpoint` CLI parameter or the `DKRON_CRONITOR_ENDPOINT` env variable.

Dkron will call Cronitor before and after running a job, sending all necessary information, no further configuration is necessary as Cronitor automagically create a monitor for each job.

![](/img/cronitor1.jpg)

You can further configure your job using the Cronitor UI or API after it has run at least once.

Cronitor can be used to notify over different channels using integrations, go to the settings page to integrate it with several popular services:

![](/img/cronitor2.jpg)
