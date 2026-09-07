package provider

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccPreCheck(t *testing.T) {
	if _, err := exec.LookPath("fnox"); err != nil {
		t.Skip("fnox binary not found on PATH, skipping acceptance test")
	}
}

func TestAccSecretDataSource(t *testing.T) {
	configPath, err := filepath.Abs("testdata/fnox_secret/fnox.toml")
	if err != nil {
		t.Fatalf("failed to resolve testdata path: %v", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
provider "fnox" {
  config_path = %q
}

data "fnox_secret" "test" {
  key = "TEST_SECRET"
}
`, configPath),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.fnox_secret.test", "key", "TEST_SECRET"),
					resource.TestCheckResourceAttr("data.fnox_secret.test", "value", "test-value-123"),
					resource.TestCheckResourceAttr("data.fnox_secret.test", "id", "TEST_SECRET"),
				),
			},
		},
	})
}
