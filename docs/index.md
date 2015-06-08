# Dcron - Distributed cron

Welcome to the Dcron documentation! This is the reference guide on how to use Drcon. If you want a getting started guide refer to the [getting started guide](getting-started/) of the Dcron documentation.

## What is Dcron

Dcron is a system service that runs scheduled tasks at given intervales or times, just like the cron unix service. It differs from it in the sense that it's distributed in several machines in a cluster and if one of that machines (the leader) fails, any other one can take this responsability and keep executing the sheduled tasks without human intervention.

## Dcron design

Dcron is designed to do one task well, executing commands in given intervals, following the unix philosophy of doing one thing and doing it well, like the classic and battle tested cron unix service, with the given addition of being designed for the cloud era, removing single points of failure and clusters of any size are needed to execute scheduled tasks in a decentralized fashion.
