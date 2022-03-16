---
title: Access Control
---

# Access Control (Preview)

{{% notice info %}}
This feature is in preview and is subject to big changes
{{% /notice %}}

Dkron provides an optional Access Control List (ACL) system which can be used to control access to data and APIs. The ACL is Capability-based, relying on policies to determine which fine grained rules can be applied. Dkron's capability based ACL system is very similar to common ACL systems you are used to.

## ACL System Overview

Dkron's ACL system is implemented with the CNCF [Open Policy Agent](https://www.openpolicyagent.org/) bringing a powerful system to suit your needs.

The ACL system is designed to be easy to use and fast to enforce while providing administrative insight. At the highest level, there are two major components to the ACL system:

* **OPA policy engine.** OPA provices policy decission making [decoupling](https://www.openpolicyagent.org/docs/latest/philosophy/#policy-decoupling) Dkron integrates OPA as a library and provides a default policy rules written in the OPA Policy language that implements a set of enforcing rules on request params to the API that are ready to use for most use cases. You don not need to learn the OPA Policy language to start using Dkron's ACL system, but you can modify the default policy rules to adapt to your use case if you need to. Read more in [OPA Docs](https://www.openpolicyagent.org/docs/latest/)

* **ACL Policies.** Dkron's ACL policies are simple JSON documents that define patterns to allow access to resources. You can find below an example ACL policy that works with the default OPA policy. The ACL JSON structure is not rigid you can adapt it to add new features in combination with the OPA Policy rules.

{{% notice note %}}
This guide is based on the usage of the default OPA Rego Policy
{{% /notice %}}

## Configuring ACLs

ACLs are not enabled by default and must be enabled. To enable ACLs simply create an ACL policy using the API. Below you can find the most basic example of an ACL policy:

Basic example policy:
```
curl localhost:8080/v1/acl/policies -d '{
    "path": {
        "/v1": {
            "capabilities": [
                "read",
            ]
        },
        "/v1/**": {
            "capabilities": [
                "create",
                "read",
                "update",
                "delete",
                "list"
            ]
        }
    }
}'
```

This policy allows any request to the API. As you can see paths uses glob patterns, and capabilities allow operations on resources.

ACLs also allows templating, providing the ability to allow or deny operations to certain resource by patterns without having to hardcode values in policies.

For example, we can for limit job actions on certain resources based on the provided token via the accepted header `X-Dkron-Token` on the request:

Example policy:
```
curl localhost:8080/v1/acl/policies -d '{
    "path": {
        "/v1/members": {
            "capabilities": ["read"]
        },
        "/v1/jobs": {
            "capabilities": [
                "list",
                "read"
            ]
        },
        "/v1/jobs/{{.Token}}-*": {
            "capabilities": [
                "create",
                "read",
                "update",
                "delete"
            ]
        }
    }
}'
```

This policy will allow all operations on jobs starting with `[Token]-job_name`, but will deny manipulation of jobs that doesn't match the pattern.

## Disable ACLs

As an administrator you will need to edit policies. Currently to be able to edit ACLs if you get locked out, you need to edit the default Rego file and disable enforcement completely. Edit the file located in `policies/main.rego` and change the `default allow` directive to `true`:

```
default allow = false -> true
```

This way the policy engine always evaluates to true, allowing any operation again. To restore ACL enforcemen, edit again the `default allow` line and set it back to `false`.
