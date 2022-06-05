---
Description: "Dkron release 3.2 public release"
Keywords:    [ "Development", "OpenSource", "Distributed systems", "cron" ]
Tags:        [ "Development", "OpenSource", "Distributed systems", "cron" ]
date:        "2021-06-01"
Topics:      [ "Development", "OpenSource", "Distributed Systems" ]
Slug:        "dkron-3-2"
---

# Dkron 3.2

## New website

Our brand new web site designed and implemented by https://github.com/Macxim, comes with a brand new look, better content structure, better documentation a new blog section and the new API navigator.

This marks the start of a new and better product design, more focused on the UX, easy of use and more documentation for Dkron.

We hope you like it as much as we do. ❤️

## New features

### Cronitor integration

Our goal is to provide a very reliable way of running your cron jobs, we share that vision with the people behind [Cronitor](https://cronitor.io/). Dkron is very reliable but sometimes a very bad event can bring your cluster down to its knees. To provide multiple options to monitor Dkron, I'm happy to introduce a new way to monitor your job executions using Cronitor service.

Cronitor is tightly integrated with Dkron, it will notify the details of every execution and Cronitor can offer multiple channels for alerting you in case something goes wrong.

Check the service https://cronitor.io/ and follow the integration guide in the docs to set up your [Dkron-Cronitor integration](/docs/usage/cronitor).

### OpenAPI

We have migrated API docs from Swagger 2 spec to OpenAPI 3 format. Check the new and better [API docs](/api) and you can download the new [OpenAPI spec](/openapi/openapi.yaml) too.
## Upcoming features

The new look of Dkron will pave the road for the upcoming v4 release. We have really interesting features almost ready for the new version, some of them are:

* Bump React Admin to v4.0
* Shell plugin will be included in the main binary
* New light image including only the main binary
* Optionally use fast-store instead of boltdb for Raft log, this will improve performance tenfold
* Farewell to the old UI

We think this will open Dkron to be used for new use cases that where not possible before.

## Wrap-up

We are very happy of giving Dkron a well deserved new face to the world and also to keep integrating with new services we love to be able to offer the best product we can for this specific -and niche- market.

We think there's a gap in Job schedulers for the rest-of-us that is currently improving, but still very needed of cost-effective and easy to operate solutions like Dkron for small-mid start-ups and for specific needs in bigger companies.

We're always open to our users feedback so feel free to contact us if you have any suggestion.
