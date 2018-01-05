# Specify SolidFire resources
resource "solidfire_account" "hackathon-account" {
    username = "Hackathon Account"
}

resource "solidfire_volume" "hackathon-volume" {
    name = "Hackathon-Volume-${count.index}"
    account_id = "${solidfire_account.hackathon-account.id}"
    total_size = "${var.total_size[count.index]}"
    enable512e = true
    min_iops = 50
    max_iops = 10000
    burst_iops = 10000

    # Create 1 instances
    count = "${length(var.total_size)}"
}

resource "solidfire_volume_access_group" "hackathon-group" {
    name = "Hackathon-Group"
    volumes = ["${solidfire_volume.hackathon-volume.*.id}"]
}

resource "solidfire_initiator" "hackathon-initiator" {
    name = "iqn.1998-01.com.vmware:bdr-es65-7f17a50c"
    alias = "Hackathon EUI Cluster"
    volume_access_group_id = "${solidfire_volume_access_group.hackathon-group.id}"
    iqns = ["${solidfire_volume.hackathon-volume.*.iqn}"]
}