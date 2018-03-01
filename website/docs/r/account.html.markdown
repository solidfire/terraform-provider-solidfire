---
layout: "solidfire"
page_title: "SolidFire: solidfire_account"
sidebar_current: "docs-solidfire-resource-account"
description: |-
  Provides a SolidFire cluster account resource. This can be used to ...
---

# solidfire\_account

Provides a SolidFire cluster account resource. This can be used to ...

## Example Usages

**Create SolidFire cluster account:**

```
resource "solidfire_account" "main-account" {
  username         = "main"
  initiator_secret = "s!39naDlLa9"
  target_secret    = "2Z>D0jf3Dpa"
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required) The name of the SolidFire account.
* `initiator_secret` - (Optional) The initiator secret. If not specified, the SolidFire cluster will autogenerate
  an initiator secret.
* `target_secret` - (Optional) The target secret. If not specified, the SolidFire cluster will autogenerate
  an initiator secret.
  
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the account.