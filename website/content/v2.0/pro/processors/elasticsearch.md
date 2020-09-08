---
title: Elasticsearch processor
---

The Elasticsearch processor can fordward execution logs to an ES cluster. It need an already available Elasticsearch installation that is visible in the same network of the target node.

The output logs of the job execution will be stored in the indicated ES instace.

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
