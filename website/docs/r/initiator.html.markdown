---
layout: "solidfire"
page_title: "SolidFire: solidfire_initiator"
sidebar_current: "docs-solidfire-resource-initiator"
description: |-
  Provides a SolidFire cluster initiator resource. This can be used to create a new initiator IQN or 
  World Wide Port Names (WWPNs).
---

# solidfire\_initiator

Provides a SolidFire cluster initiator resource. This can be used to create a new initiator IQN or 
World Wide Port Names (WWPNs).

## Example Usages

**Create SolidFire cluster initiator:**

```
resource "solidfire_initiator" "main-initiator" {
  name = "qn.1998-01.com.vmware:test-terraform-00000000"
  alias = "Terraform Main Initiator"
  volume_access_group_id = "123"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SolidFire initiator.
* `alias` - (Optional) The user-friendly alias of the SolidFire initiator.
* `volume_access_group_id` - (Optional) The ID of the SolidFire volume access group
  to use with the initiator.
  
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the initiator.