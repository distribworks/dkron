#Getting started

Welcome to the intro guide to Dcron! This will explain how to setup serf, how easy is to use it, what problems could it help you to soleve, etc.

## What is Dcron

Dcron is a system service that runs scheduled tasks at given intervales or times, just like the cron unix service. It differs from it in the sense that it's distributed in several machines in a cluster and if one of that machines (the leader) fails, any other one can take this responsability and keep executing the sheduled tasks without human intervention.

## Dcron design

Dcron is designed to do one task well, executing commands in given intervals, following the unix philosophy of doing one thing and doing it well, like the classic and battle tested cron unix service, with the given addition of being designed for the cloud era, removing single points of failure and clusters of any size are needed to execute scheduled tasks in a decentralized fashion.
