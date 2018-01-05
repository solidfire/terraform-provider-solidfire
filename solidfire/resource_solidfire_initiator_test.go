package solidfire

import (
	"strconv"
	"testing"

	"bitbucket.org/solidfire/terraform-provider-solidfire/solidfire/element"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestInitiator_basic(t *testing.T) {
	var initiator element.Initiator
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSolidFireInitiatorDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckSolidFireInitiatorConfig,
					"terraform-acceptance-test",
					"terraform-acceptance-test-alias",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSolidFireInitiatorExists("solidfire_initiator.terraform-acceptance-test-1", &initiator),
					resource.TestCheckResourceAttr("solidfire_initiator.terraform-acceptance-test-1", "name", "terraform-acceptance-test"),
					resource.TestCheckResourceAttr("solidfire_initiator.terraform-acceptance-test-1", "alias", "terraform-acceptance-test-alias"),
				),
			},
		},
	})
}

func TestInitiator_update(t *testing.T) {
	var initiator element.Initiator
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSolidFireInitiatorDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckSolidFireInitiatorConfig,
					"terraform-acceptance-test",
					"terraform-acceptance-test-alias",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSolidFireInitiatorExists("solidfire_initiator.terraform-acceptance-test-1", &initiator),
					resource.TestCheckResourceAttr("solidfire_initiator.terraform-acceptance-test-1", "name", "terraform-acceptance-test"),
					resource.TestCheckResourceAttr("solidfire_initiator.terraform-acceptance-test-1", "alias", "terraform-acceptance-test-alias"),
				),
			},
			{
				Config: fmt.Sprintf(
					testAccCheckSolidFireInitiatorConfigUpdate,
					"terraform-acceptance-test",
					"terraform-acceptance-test-alias-update",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSolidFireInitiatorExists("solidfire_initiator.terraform-acceptance-test-1", &initiator),
					resource.TestCheckResourceAttr("solidfire_initiator.terraform-acceptance-test-1", "name", "terraform-acceptance-test"),
					resource.TestCheckResourceAttr("solidfire_initiator.terraform-acceptance-test-1", "alias", "terraform-acceptance-test-alias-update"),
				),
			},
		},
	})
}

func TestInitiator_removeVolumeAccessGroup(t *testing.T) {
	var initiator element.Initiator
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSolidFireInitiatorDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckSolidFireInitiatorConfig,
					"terraform-acceptance-test",
					"terraform-acceptance-test-alias",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSolidFireInitiatorExists("solidfire_initiator.terraform-acceptance-test-1", &initiator),
					resource.TestCheckResourceAttr("solidfire_initiator.terraform-acceptance-test-1", "name", "terraform-acceptance-test"),
					resource.TestCheckResourceAttr("solidfire_initiator.terraform-acceptance-test-1", "alias", "terraform-acceptance-test-alias"),
				),
			},
			{
				Config: fmt.Sprintf(
					testAccCheckSolidFireInitiatorConfigRemoveVAG,
					"terraform-acceptance-test",
					"terraform-acceptance-test-alias-update",
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSolidFireInitiatorExists("solidfire_initiator.terraform-acceptance-test-1", &initiator),
					resource.TestCheckResourceAttr("solidfire_initiator.terraform-acceptance-test-1", "name", "terraform-acceptance-test"),
					resource.TestCheckResourceAttr("solidfire_initiator.terraform-acceptance-test-1", "alias", "terraform-acceptance-test-alias-update"),
				),
			},
		},
	})
}

func testAccCheckSolidFireInitiatorDestroy(s *terraform.State) error {
	virConn := testAccProvider.Meta().(*element.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "solidfire_initiator" {
			continue
		}

		_, err := virConn.GetInitiatorByID(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Error waiting for initiator (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckSolidFireInitiatorExists(n string, initiator *element.Initiator) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		virConn := testAccProvider.Meta().(*element.Client)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SolidFire initiator key ID is set")
		}

		retrievedInit, err := virConn.GetInitiatorByID(rs.Primary.ID)
		if err != nil {
			return err
		}

		convID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		if retrievedInit.InitiatorID != convID {
			return fmt.Errorf("Resource ID and initiator ID do not match")
		}

		*initiator = retrievedInit

		return nil
	}
}

const testAccCheckSolidFireInitiatorConfig = `
resource "solidfire_initiator" "terraform-acceptance-test-1" {
	name = "%s"
	alias = "%s"
	volumeAccessGroupID = "${solidfire_volume_access_group.terraform-acceptance-test-1.id}"
}

resource "solidfire_volume_access_group" "terraform-acceptance-test-1" {
	name = "terraform-acceptance-test-group"
}
`

const testAccCheckSolidFireInitiatorConfigUpdate = `
resource "solidfire_initiator" "terraform-acceptance-test-1" {
	name = "%s"
	alias = "%s"
	volumeAccessGroupID = "${solidfire_volume_access_group.terraform-acceptance-test-2.id}"
}

resource "solidfire_volume_access_group" "terraform-acceptance-test-2" {
	name = "terraform-acceptance-test-group-2"
}
`

const testAccCheckSolidFireInitiatorConfigRemoveVAG = `
resource "solidfire_initiator" "terraform-acceptance-test-1" {
	name = "%s"
	alias = "%s"
}
`
