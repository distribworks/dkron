---
title: Concurrency
toc: true
---

## Concurrency

Jobs can be configured to allow overlapping executions or forbid them. 

Concurrency property accepts two option: 

* **allow** (default): Allow concurrent job executions.
* **forbid**: If the job is already running don't send the execution, it will skip the executions until the next schedule.
