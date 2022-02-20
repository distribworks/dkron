
---
title: Shell Executor
---

Shell executor runs a system command

## Configuration

Params

```
shell: Run this command using a shell environment
command: The command to run
env: Env vars separated by comma
cwd: Chdir before command run
```

Example

```json
{
  "executor": "shell",
  "executor_config": {
      "shell": "true",
      "command": "my_command",
      "env": "ENV_VAR=va1,ANOTHER_ENV_VAR=var2",
      "cwd": "/app"
  }
}
```
