## 0.7.1 (2016-05-03)

- Don't use shell call when executing commands, exploding the command line.
- Add advertise, add `advertise` option that solves joining between hosts when running docker
- Validate job size, limit to serf maximum size
- Job overwrite, now sending existing jobs doesn't overwrite non existing fields in request
- Fix for dashboard crash on non existent leader

## 0.7.0 (2016-04-10)

- Refactor leader election, the old method could lead to cases where 2 or more nodes could have the scheduler running without noticing the other master.
- Get rid of `keys`, in a serf cluster node names are unique so using it for leader keys now.
- Fix [#85](https://github.com/victorcoder/dkron/issues/85) Restart scheduler on job deletion
- Refactor logging, replace `debug` with `log-level`
- Order nodes in UI [#81](https://github.com/victorcoder/dkron/issues/81) (kudos @whizz)
- Add exposed vars to easy debugging
- Go 1.6
- Add @minutely as predefined schedule (kudos @mlafeldt)

### Upgrade from 0.6.x

To upgrade an existing installation you must first delete the pre-exiting leader key from the store. The leader key is in the form of: `[keyspace]/leader`

## 0.6.4 (2016-02-18)

- Use expvars to expose metrics
- fix https://github.com/victorcoder/dkron/issues/71
- Better example config in package and docs

## 0.6.3 (2015-12-28)

- UI: Better job view
- Logic to store only the last 100 executions

## 0.6.2 (2015-12-22)

- Fixed [#62](https://github.com/victorcoder/dkron/issues/55)

## 0.6.1 (2015-12-21)

- Fixed bugs [#55](https://github.com/victorcoder/dkron/issues/55), [#52](https://github.com/victorcoder/dkron/issues/52), etc.
- Build for linux arm

## 0.6.0 (2015-12-11)

- Some other improvements and bug fixing
- Vendoring now using Go vendor experiment + glide
- Fix: Remove executions on job delete
- Show full execution output in UI modal
- New executions results internals using RPC
- Standarized logging
- Show job tooltips with info
- Accept just "pretty" for formatting api requests
- Change how execution groups work to not use the directory concept.

## 0.5.5 (2015-11-19)

- More backend compatibility
- Accept just pretty for formatting api requests
- Show executions grouped in web UI
- Show job tooltips with all job JSON info in web UI
- Better alerts

## 0.5.4 (2015-11-17)

- Fix to web UI paths

## 0.5.3 (2015-11-16)

- Web UI works behind http proxy

## 0.5.2 (2015-11-09)

- Fix bug in join config parameter that rendered it unusable from config file.

## 0.5.1 (2015-11-06)

- Deb package
- Upgraded libkv to latest
- New config options (log level, web UI path)

## 0.5.0 (2015-09-27)

- Email and Webhook configurable notifications for job executions.
- Ability to encrypt serf network traffic between nodes.
- Pretty formating API responses
- UI now shows the execution status with color coding and partial execution.
- More API stability and predictability
- Provided API JSON schema, generated API docs based in the schema
- Tested on Travis
- Using Libkv allows to use different storage backends (etcd, consul, zookeeper)
- Add v1 versioning to the API routes

## 0.0.4 (2015-08-23)

- Compiled with Go 1.5
- Includes cluster nodes view in the UI

## 0.0.3 (2015-08-20)

- Initial release
