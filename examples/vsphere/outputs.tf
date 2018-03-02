output "vmfs_disks" {
    value = "${data.vsphere_vmfs_disks.available.disks}"
}