---
title: Job chaining
---

## Job chaining

You can set some jobs to run after other job is executed. To setup a job that will be executed after any other given job, just set the `parent_job` property when saving the new job.

The dependent job will be executed after the main job finished a successful execution.

Child jobs schedule property will be ignored if it's present.
