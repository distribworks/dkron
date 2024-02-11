# Elasticsearch processor

The Elasticsearch processor can fordward execution logs to an ES cluster. It need an already available Elasticsearch installation that is visible in the same network of the target node.

The output logs of the job execution will be stored in the indicated ES instace.

:::info For Dkron Pro < v3.2.3
:::
## Configuration

```json
{
  "processors": {
    "elasticsearch": {
      "url": "http://localhost:9200", //comma separated list of Elasticsearch hosts urls (default: http://localhost:9200)
      "index": "dkron_logs", //desired index name (default: dkron_logs)
      "forward": "false" //forward logs to the next processor (default: false)
    }
  }
}
```

:::info For Dkron Pro > v3.2.3
:::

## Configuration

For increased security and flexibility, configuration of the ES processor is stored in a file named `dkron-processor-elasticsearch.yml` in the same locations as `dkron.yml`, and should include a list of configurations for Elasticsearch, it can include any number of configurations.

This is an example including all available parameters:
```yaml
es1:
  index: dkron-logs
  index_date_format: '2006-01-02'
  username: elastic
  password: XXXXXXXXXXXXX
  url: https://localhost:9200
```
And for each job you only need to configure the `config` parameter in the processors configuration:

```json
{
  "processors": {
    "elasticsearch": {
      "config": "es1", // configuration to use from the config file
      "forward": "false" // forward logs to the next processor (default: false)
    }
  }
}
```
