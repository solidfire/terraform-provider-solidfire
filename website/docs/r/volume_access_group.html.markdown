---
layout: "solidfire"
page_title: "SolidFire: solidfire_volume_access_group"
sidebar_current: "docs-solidfire-resource-volume-access-group"
description: |-
  Provides a SolidFire cluster volume access group resource. This can be used to create a new volume access group.
  Any initiator IQN that you add to the volume access group is able to access any volume in the group without CHAP
  authentication.
---

# solidfire\_volume\_access\_group

Provides a SolidFire cluster volume access group resource. This can be used to create a new volume access group.
Any initiator IQN that you add to the volume access group is able to access any volume in the group without CHAP
authentication.

## Example Usages

**Create SolidFire cluster volume access group:**

```
resource "solidfire_volume_access_group" "main-group" {
  name = "terraform-main-group"
  volumes = ["12345", "67890"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SolidFire volume access group.
* `volumes` - (Optional) The IDs of the SolidFire volumes to add to the
  SolidFire volume access group.
  
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The unique identifier for the volume access group.
* `initiators` - Any initiators tied to the volume access group.