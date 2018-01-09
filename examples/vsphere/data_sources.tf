data "vsphere_datacenter" "datacenter" {
    name = "${var.datacenter}"
}

data "vsphere_host" "esxi_host" {
    name = "10.117.149.43"
    datacenter_id = "${data.vsphere_datacenter.datacenter.id}"
}

data "vsphere_vmfs_disks" "available_disks" {
    host_system_id = "${data.vsphere_host.esxi_host.id}"
    rescan = true
    filter = "naa.6f47acc1"

    depends_on = ["solidfire_volume.vsphere-volume"]
}