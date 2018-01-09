# Specify SolidFire resources
resource "solidfire_account" "vsphere-account" {
    username = "vSphere Account"
}

resource "solidfire_volume" "vsphere-volume" {
    name = "vSphere-Volume-${count.index}"
    account_id = "${solidfire_account.vsphere-account.id}"
    total_size = "${var.total_size[count.index]}"
    enable512e = true
    min_iops = 50
    max_iops = 10000
    burst_iops = 10000

    # Create 1 instances
    count = "${length(var.total_size)}"
}

// resource "solidfire_volume_access_group" "vsphere-group" {
//     name = "vSphere-Group"
//     volumes = ["${solidfire_volume.vsphere-volume.*.id}"]
// }

// resource "solidfire_initiator" "vsphere-initiator" {
//     name = "iqn.1998-01.com.vmware:bdr-es65-7f17a50c"
//     alias = "vSphere Cluster"
//     volume_access_group_id = "${solidfire_volume_access_group.vsphere-group.id}"
//     iqns = ["${solidfire_volume.vsphere-volume.*.iqn}"]
// }

resource "vsphere_vmfs_datastore" "datastore" {
    name = "terraform-test"
    host_system_id = "${data.vsphere_host.esxi_host.id}"

    disks = ["${data.vsphere_vmfs_disks.available_disks.disks}"]
}