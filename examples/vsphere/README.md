# SolidFire 0.1.0 Example

This repository is designed to demonstrate the capabilities of the [Terraform
SolidFire Provider][ref-tf-solidfire] at the time of the 0.1.0 release, working
in conjunction with the [Terraform vSphere Provider][ref-tf-vsphere].

[ref-tf-solidfire]: https://www.terraform.io/docs/providers/solidfire/index.html
[ref-tf-vsphere]: https://www.terraform.io/docs/providers/vsphere/index.html

**NOTE:** This example is not completely working.

This example performs the following:

* Sets up an account. This uses the
  [`solidfire_account` resource][ref-tf-solidfire-account].
* Creates a number of volumes on the cluster tied to the account,
  using the [`solidfire_volume` resource][ref-tf-solidfire-volume].
* Sets up a volume access group for the volumes, using the
  [`solidfire_volume_access_group` resource][ref-tf-solidfire-volume-access-group].
* Creates an initiator tied to the volume access group and volumes using 
  the [`solidfire_initiator` resource][ref-tf-solidfire-initiator].

[ref-tf-solidfire-account]: https://www.terraform.io/docs/providers/solidfire/r/account.html
[ref-tf-solidfire-initiator]: https://www.terraform.io/docs/providers/solidfire/r/initiator.html
[ref-tf-solidfire-volume]: https://www.terraform.io/docs/providers/solidfire/r/volume.html
[ref-tf-solidfire-volume-access-group]: https://www.terraform.io/docs/providers/solidfire/r/volume_access_group.html

## Requirements

* A working SolidFire storage cluster.
* A VMware vSphere instance attached to the SolidFire storage cluster

## Usage Details

You can either clone the entire
[terraform-provider-solidfire][ref-tf-solidfire-github] repository, or download the
`provider.tf`, `variables.tf`, `resources.tf`, and
`terraform.tfvars.example` files into a directory of your choice. Once done,
edit the `terraform.tfvars.example` file, populating the fields with the
relevant values, and then rename it to `terraform.tfvars`. Don't forget to
configure your endpoint and credentials by either adding them to the
`provider.tf` file, or by using enviornment variables. See
[here][ref-tf-solidfire-provider-settings] for a reference on provider-level
configuration values.

[ref-tf-solidfire-github]: https://github.com/terraform-providers/terraform-provider-solidfire
[ref-tf-solidfire-provider-settings]: https://www.terraform.io/docs/providers/solidfire/index.html#argument-reference

Once done, run `terraform init`, and `terraform plan` to review the plan, then
`terraform apply` to execute. If you use Terraform 0.11.0 or higher, you can
skip `terraform plan` as `terraform apply` will now perform the plan for you and
ask you confirm the changes.