data "vsphere_datacenter" "datacenter" {
    name = "${var.datacenter}"
}

data "vsphere_host" "esxi_host" {
    name = "esxi1"
    datacenter_id = "${data.vsphere_datacenter.datacenter.id}"
}

data "vsphere_vmfs_disks" "available_disks" {
    host_system_id = "${data.vsphere_host.esxi_host.id}"
    rescan = true
    filter = ""
}