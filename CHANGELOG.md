## Unreleased

- Add DynamoDB support
- Breaking: Remove support and deprecation message for old command, environment_variables and shell parameters

## 1.0.2

- Allow sending mail without credentials
- Fix docker tagging
- Log plugin fix and improvements
- More specific processor plugin usage logging

## 1.0.1

- Don't dockerignore dist folder, it is needed by gorelease docker builder
- fa323c2 Ignore node modules
- 2475b37 Move static assets to it's own directory inside static folder
- 987dd5d Remove hash from url on modal close
- 455495c Remove node_modules
- 0c02ce0 Reorg asset generation in subpackages

## 1.0.0

- c91852b Cookie consent in website
- 9865012 Do not install build tools
- e280d31 Docker login
- 3229dc2 Ensure building static binaries
- 01e62b6 Error on test
- 69380f5 Ignore system files
- a02a1ab Release with docker
- e795210 Remove old dockerignore entries
- c9c692c Remove unmarshalTags from dkron.go
- c5f5de0 Report errors on unmarshal config
- 62e1e15 Sums for release
- 1cf235a UnmarshalTags belongs to the agent and should be public
- 36f9318 Update readme
- 80b2ab1 mail-port flag is uint

## 1.0.0-rc4

- 913ee87 Bump mapstructure
- 5bd120f Remove legacy config loading
- f20fbe5 Update mods

## 1.0.0-rc3

- 4811e48 Fix UI run and delete
- 8695242 Redirect to dashboard

## 1,0.0-rc2

- d6dbb1a Add toggle to swagger
- ffa4feb Deep linking to job views
- fdc5344 Don't fury on make
- 236b5f4 Don't query jobs on interval in Dashboard
- ea5e60b Fixes rescheduling on boltdb store
- f55e2e3 Gardening and anchor links open modal
- b22b362 Gen
- 6887c36 Logging info
- d21cf16 Logging info and use store.Backend type instead of strings in config
- 28c130b Open modal with anchor links and gardening
- 1afb3df Several ui fixes introduced when migrating to glyphicons


## 1.0.0-rc1

- ef86e13 Bump go-plugin
- d09b942 Bump several dependencies
- f96d622 GRPC
- 8e3b4b9 Ignore dist folder
- 1b7d4bc Issue template
- caf4711 Logrus
- 5821c8c Mainly etcd
- 33a12c5 Revert "Bump several dependencies"
- fb9460d Update cron-spec.md
- 706e65d Upgrade pflag

## 0.11.3

- 723326f Add logging for pending executions response
- df76e9c Add real examples to swagger spec
- d1318a1 Add tags param to swagger spec
- 4da0b3b Big docs refactor
- 2d91a5e Break on errors
- 8fac831 Command to generate cli docs from cobra
- bdcd09c Don't use swagger2markup
- 253fe57 ECS and email pro docs
- e89b353 Expvar dep
- 187190e Fix indentation
- c8320b5 Fix testing
- 9037d65 Fix typo in getting started
- 9c60fe8 Formatting
- f11ed84 Formatting
- 20be8e5 Integrate swagger-ui for a bit better API visualization
- 2cede00 Merge branch 'master' into boltdb
- 53d8464 Only query for pending executions when there is some
- 712be35 Remove extra useless locking introduced in 88c072c
- dacb379 This should be TrimSuffix
- dec6701 Update contacts
- c21e565 Update getting-started.md
- 3fdba5f Use boltdb as default storage
- 70d9229 Wrong dash in example config file
- 9653bbc expvars are back and simple health endpoint

## 0.11.2

- 7d88742 Add code of conduct
- aed2f44 Proper serf debug logging
- 1226c93 Publish docker
- a0b6f59 Publish docker
- f1aaecc Reorg imports
- 8758bac Tests should use etcdv3
- fa3aaa4 Tests should use v3 client
- 5bcea4c Update create or update job api endpoint
- 39728d0 refactor: Methond name
- 1c64da4 refactor: Proper gin logging and mode

## 0.11.1 (2018-10-07)

 - Add support for passing payload to command STDIN (@gustavosbarreto)
 - add support for etcdv3 (@kevynhale)
 - Use etcdv3 by default
 - Jobs static URLs fixed

## 0.11.0 (2018-09-24)

- 1.11 stable not in docker hub yet
- Add builtin http plugin
- Add executor shell su option (@tengattack)
- Better dockerfile for testing
- Better flag help
- Don't depend on michellh/cli
- Filter jobs by tags (@digitalcrab)
- Fix cluster panic bug (@tengattack)
- Release with goreleaser
- Use cobra for flags
- Use go modules
- add create & update job features (@wd1900)

## 0.10.4 (2018-07-30)

- Replace RPC with gRPC
- Fix compose files (@kevynhale)

## 0.10.3 (2018-06-20)

- Replace goxc with makefile
- Pro docs

## 0.10.2 (2018-05-23)

### Bug fixes

- Fix status check
- Remove unnecessary updates of job finish times (@sysadmind)
- Reflect store status in API
- Fix windows plugins (@sysadmind)
- Stop Job update on JSON parse error in API (@gromo)

## 0.10.1 (2018-05-17)

### Bug fixes

- Fix dashboard job view/delete modals

## 0.10.0 (2018-05-14)

### Bug fixes

- Fix RPCconfig query missing address (#336 and related)
- Fix concurrency issue due to race condition on lock jobs (#299)
- Fix execution done missing on restart blocking concurrency forbid jobs (#349)
- Fix plugin load paths (#275)
- Fix RPC address lost on reload config (#262)

### Features and code improvements

- Slightly improve processing of last execution group (@sysadmind)
- Improve job dependencies handling (@sysadmind)
- Move dkron command to it's own package
- Milliseconds range API job create or update
- Refactor scheduler restart
- Replace bower with npm
- Executor plugins based on GRPC
- Toggle job from UI
- Search Job by name and pagination in the UI
- Add redis as storage backend (@digitalcrab)
- UI refactor with new bootstrap version and replace fontawesome with glyphicons
- Compute job status and return the value from the API providing the user with more info
- Timezone aware scheduling (@didiercrunch)


## 0.9.8 (2018-04-27)

- Fix broken release 0.9.7

## 0.9.7 (2018-02-12)

- Less verbose plugin logging
- Update broken osext dep (@ti)
- Switch from libkv to valkeyrie
- Refactor for usable core code
- Fix unsorted execution groups (@firstway)
- Fix GetLastExecutionGroup (@firstway)

## 0.9.6 (2017-11-14)

- Migrate from glide to dep
- Fix params precedence, cli params on top
- More robust test suite
- Gin logging to common logger
- Better systemd script
- Don't panic or fatal when sending notification
- Serf upgrade
- Fix templating breaking change on Go 1.9 upgrade

## 0.9.5 (2017-09-12)

Features

- New docs website using hugo

Bug fixes:

- Clean up clients upon an exit signal (@danielhan)
- Fix #280 (@didiecrunch)
- Upgrade several dependencies
- Fix static assets relative path

## 0.9.4 (2017-06-07)

- Fix mistakes in API docs
- Using "jobs", "1am" or "1pm" in the name of job leads to a dashboard bug
- Fix crash on non existent plugin name
- Embed all assets in binary, removed -ui-dir config param

### Upgrade notes

This is a breaking change; `ui-dir` configuration param has been removed, all scripting using this param should be updated.

## 0.9.3 (2017-02-12)

- Fix RPC server listen address (@firstway)
- Basic implementation of the testing infrastructure using swarm
- Basic Telemetry implementation for sending metrics to statsd and datadog
- Fix crash on backend failure
- Reverse sort executions in UI (@Eyjafjallajokull)

## 0.9.2 (2016-12-28)

Features:

- Implement concurrency policy
- Improved UI: allow delete jobs from UI, highlight JSON
- Execution Processor plugins, allows flexible routing of execution results
- Template variables for customization of notification emails (@oldmantaiter)
- Go 1.7
- Test with docker-compose, this will allow to test multiple stores easily

Bug fixes:

- Fix tests randomly failing (@oldmantaiter)
- Return empty list when no jobs instead of null
- Allow POST usage on /leave method, deprecate GET

## 0.9.1 (2016-09-24)

Bug fixes:

- Fix job stats not being updated #180
- Fix zookeeper get list of executions #184
- Fix crash when deleting a job that doesn't exists #182
- Fix Travis in forks

## 0.9.0 (2016-08-24)

Features:

- Support any size jobs
- Support chained jobs
- Schedule and other job properties validation
- New site, logo and dashboard design

Bug fixes:

- Fix execution retries
- Fix executions merge by same prefix
- Fix correct HTTP status on create/update

## 0.7.3 (2016-07-12)

- One off jobs
- Added cron spec to docs
- Execution retry on failure
- Switch JSON schema spec for it's corresponding Open API spec
- Reload config
- Fix scheduling bug
- New job Status gives more info on the job execution

## 0.7.2 (2016-06-01)

- Add some helpers and bugfixes
- Add shell property to job, reintroduced the shell execution method but now it's a choice
- Add reporting node to execution reports
- Replace server tag for dkron_server and add dkron_version

### Upgrade notes

Due to the change in the internal tags `server` to `dkron_server`, you'll need to adjust job tags if they where using that tag.

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
