package solidfire

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"solidfire": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SOLIDFIRE_USERNAME"); v == "" {
		t.Fatal("SOLIDFIRE_USERNAME must be set for acceptance tests")
	}

	if v := os.Getenv("SOLIDFIRE_PASSWORD"); v == "" {
		t.Fatal("SOLIDFIRE_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("SOLIDFIRE_SERVER"); v == "" {
		t.Fatal("SOLIDFIRE_SERVER must be set for acceptance tests")
	}

	if v := os.Getenv("SOLIDFIRE_API_VERSION"); v == "" {
		t.Fatal("SOLIDFIRE_API_VERSION must be set for acceptance tests")
	}
}
