data "vsphere_datacenter" "datacenter" {
    name = "${var.vsphere_datacenter}"
}

data "vsphere_host" "host" {
    name = "${var.vsphere_hosts[1]}"
    datacenter_id = "${data.vsphere_datacenter.datacenter.id}"
}

data "vsphere_vmfs_disks" "available" {
    host_system_id = "${data.vsphere_host.host.id}"
    rescan = true
    filter = "naa.6f47acc1"

    depends_on = ["solidfire_initiator.vsphere-initiator"]
}

data "vsphere_resource_pool" "pool" {
    name = "NetApp-HCI-Cluster-01/Resources"
    datacenter_id = "${data.vsphere_datacenter.datacenter.id}"
}

data "vsphere_network" "network" {
    name = "vMotion"
    datacenter_id = "${data.vsphere_datacenter.datacenter.id}"
}