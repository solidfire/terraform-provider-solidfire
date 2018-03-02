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

    # Create N instances
    count = "${length(var.total_size)}"
}

resource "solidfire_volume_access_group" "vsphere-group" {
    name = "vSphere-Group"
    volumes = ["${solidfire_volume.vsphere-volume.*.id}"]
}

resource "solidfire_initiator" "vsphere-initiator" {
    name = "iqn.1998-01.com.vmware:jg-test-esx-02-187dae17"
    alias = "vSphere Cluster"
    volume_access_group_id = "${solidfire_volume_access_group.vsphere-group.id}"
    iqns = ["${solidfire_volume.vsphere-volume.*.iqn}"]
}

# Specify vSphere resources
resource "vsphere_vmfs_datastore" "sf-datastore" {
    name = "${solidfire_volume.vsphere-volume.name}"
    host_system_id = "${data.vsphere_host.host.id}"

    disks = ["${data.vsphere_vmfs_disks.available.disks}"]
}

# The resource will be created in vSphere, but will fail because the
# network will not respond, as there is no backing OS.
resource "vsphere_virtual_machine" "vm" {
    name = "solidfire-vm"
    resource_pool_id = "${data.vsphere_resource_pool.pool.id}"
    datastore_id = "${vsphere_vmfs_datastore.sf-datastore.id}"

    num_cpus = 2
    memory = 1024
    guest_id = ""

    network_interface {
        network_id = "${data.vsphere_network.network.id}"
    }

    disk {
        label = "disk0"
        size = 20
    }

    depends_on = ["data.vsphere_vmfs_disks.available"]
}