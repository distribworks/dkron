---
title: Plugins
---

## Intro

Plugins in Dkron allow you to add funcionality that integrates with the workflow of the job execution in Dkron. It's a powerful system that allows you to extend and adapt Dkron to your special needs.

This page documents the basics of how the plugin system in Dkron works, and how to setup a basic development environment for plugin development if you're writing a Dkron plugin.

## How it Works

Dkron execution execution processors are provided via plugins. Each plugin exposes functionality for modifying the execution. Plugins are executed as a separate process and communicate with the main Dkron binary over an RPC interface.

The code within the binaries must adhere to certain interfaces. The network communication and RPC is handled automatically by higher-level libraries. The exact interface to implement is documented in its respective documentation section.

## Installing a Plugin

Dkron searches for plugins at startup, to install a plugin just drop the binary in one of the following locations:

1. /etc/dkron/plugins
2. Dkron executable directory

{{% children  %}}
