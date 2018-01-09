# Specify the provider and access details
provider "solidfire" {
    username = "${var.solidfire_username}"
    password = "${var.solidfire_password}"
    solidfire_server = "${var.solidfire_cluster}"
    api_version = "${var.solidfire_api_version}"
}

provider "vsphere" {
    version = "~> 1.1"
    user = "${var.vsphere_username}"
    password = "${var.vsphere_password}"
    vsphere_server = "${var.vsphere_server}"
    allow_unverified_ssl = true
}