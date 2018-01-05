---
layout: "solidfire"
page_title: "Provider: SolidFire"
sidebar_current: "docs-solidfire-index"
description: |-
  The SolidFire provider is used to interact with the resources supported by
  SolidFire. The provider needs to be configured with the proper credentials
  before it can be used.
---

# SolidFire Provider

The SolidFire provider is used to interact with the resources supported by
SolidFire.
The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

~> **NOTE:** The SolidFire Provider currently represents _initial support_
and therefore may undergo significant changes as the community improves it.

## Example Usage

```
# Configure the SolidFire Provider
provider "solidfire" {
  username         = "${var.solidfire_username}"
  password         = "${var.solidfire_password}"
  solidfire_server = "${var.solidfire_server}"
  api_version      = "${var.solidfire_api_version}"
}

# Create an account
resource "solidfire_account" "main-account" {
  username = "main"
}

# Create a volume tied to an account
resource "solidfire_volume" "volume1" {
  name       = "main-volume"
  accountID  = "${solidfire_account.main-account.id}"
  totalSize  = 10000000000
  enable512e = true
  minIOPS    = 50
  maxIOPS    = 10000
  burstIOPS  = 10000
}

# Create a volume access group for the volume
resource "solidfire_volume_access_group" "main-group" {
  name    = "main-volume-access-group"
  volumes = ["${solidfire_volumes.volume1.id}"]
}

# Create an initiator for the volume access group
resource "solidfire_initiator" "main-initiator" {
  name = "qn.1998-01.com.vmware:test-terraform-00000000"
  alias = "Main Initiator"
  volumeAccessGroupID = "${solidfire_volume_access_group.main-group.id}"
}
```

## Argument Reference

The following arguments are used to configure the SolidFire Provider:

* `username` - (Required) This is the username for SolidFire API operations.
* `password` - (Required) This is the password for SolidFire API operations.
* `solidfire_server` - (Required) This is the SolidFire cluster name for SolidFire 
  API operations.
* `api_version` - (Required) This is the SolidFire cluster version for SolidFire
  API operations.

## Required Privileges

In order to use the Terraform provider as non priviledged user, (TBD):


These settings were tested with [NetApp SolidFire Element OS 9.1](TBD)
For additional information on roles and permissions, please refer to official
SolidFire documentation.

