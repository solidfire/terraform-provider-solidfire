package solidfire

import (
	"strconv"
	"testing"

	"bitbucket.org/solidfire/terraform-provider-solidfire/solidfire/element"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestVolume_basic(t *testing.T) {
	var volume element.Volume
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSolidFireVolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckSolidFireVolumeConfig,
					"terraform-acceptance-test",
					"1000000000",
					"true",
					"500",
					"10000",
					"10000",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSolidFireVolumeExists("solidfire_volume.terraform-acceptance-test-1", &volume),
					resource.TestCheckResourceAttr("solidfire_volume.terraform-acceptance-test-1", "name", "terraform-acceptance-test"),
					resource.TestCheckResourceAttr("solidfire_volume.terraform-acceptance-test-1", "totalSize", "1000000000"),
					resource.TestCheckResourceAttr("solidfire_volume.terraform-acceptance-test-1", "enable512e", "true"),
					resource.TestCheckResourceAttr("solidfire_volume.terraform-acceptance-test-1", "minIOPS", "500"),
					resource.TestCheckResourceAttr("solidfire_volume.terraform-acceptance-test-1", "maxIOPS", "10000"),
					resource.TestCheckResourceAttr("solidfire_volume.terraform-acceptance-test-1", "burstIOPS", "10000"),
				),
			},
		},
	})
}

func testAccCheckSolidFireVolumeDestroy(s *terraform.State) error {
	virConn := testAccProvider.Meta().(*element.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "solidfire_volume" {
			continue
		}

		_, err := virConn.GetVolumeByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Error waiting for volume (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckSolidFireVolumeExists(n string, volume *element.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		virConn := testAccProvider.Meta().(*element.Client)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SolidFire volume key ID is set")
		}

		retrievedVol, err := virConn.GetVolumeByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		convID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		if retrievedVol.VolumeID != convID {
			return fmt.Errorf("Resource ID and volume ID do not match")
		}

		*volume = retrievedVol

		return nil
	}
}

const testAccCheckSolidFireVolumeConfig = `
resource "solidfire_volume" "terraform-acceptance-test-1" {
	name = "%s"
	accountID = "${solidfire_account.terraform-acceptance-test-1.id}"
	totalSize = "%s"
	enable512e = "%s"
	minIOPS = "%s"
	maxIOPS = "%s"
	burstIOPS = "%s"
}
resource "solidfire_account" "terraform-acceptance-test-1" {
	username = "terraform-acceptance-test-volume"
}
`
