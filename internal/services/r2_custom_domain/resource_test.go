package r2_custom_domain_test

import (
	"os"
	"testing"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/acctest"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccCloudflareR2CustomDomain_Basic(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	// zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")
	zoneID := "f8786e3be4a0d180338bb7caf27171a5"
	// domainName := os.Getenv("CLOUDFLARE_DOMAIN")
	domainName := "terraform-r2.cfapi.net"
	resourceName := "cloudflare_r2_custom_domain." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccR2CustomDomainConfig(rnd, accountID, zoneID, domainName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("bucket_name"), knownvalue.StringExact(rnd)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("domain"), knownvalue.StringExact(domainName)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("zone_id"), knownvalue.StringExact(zoneID)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("jurisdiction"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					// statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("zone_name"), knownvalue.StringExact(domainName)),
					// statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("status"), knownvalue.NotNull()),
					// statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("status").AtMapKey("ownership"), knownvalue.NotNull()),
					// statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("status").AtMapKey("ssl"), knownvalue.NotNull()),
				},
				// Check: resource.ComposeTestCheckFunc(
				// 	resource.TestCheckResourceAttr(resourceName, "zone_name", domainName),
				// 	// resource.TestCheckResourceAttr(resourceName, "notes", "this is notes"),
				// 	// resource.TestCheckResourceAttr(resourceName, "mode", "challenge"),
				// 	// resource.TestCheckResourceAttr(resourceName, "configuration.target", "asn"),
				// 	// resource.TestCheckResourceAttr(resourceName, "configuration.value", "AS112"),
				// ),
			},
		},
	})
}

func TestAccCloudflareR2CustomDomain_Update(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	// zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")
	zoneID := "f8786e3be4a0d180338bb7caf27171a5"
	// domainName := os.Getenv("CLOUDFLARE_DOMAIN")
	domainName := "terraform-r2.cfapi.net"
	resourceName := "cloudflare_r2_custom_domain." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccR2CustomDomainConfig(rnd, accountID, zoneID, domainName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("min_tls"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("ciphers"), knownvalue.Null()),
				},
			},
			{
				Config: testAccR2CustomDomainUpdateConfig(rnd, accountID, zoneID, domainName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("enabled"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("min_tls"), knownvalue.StringExact("1.2")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("ciphers"), knownvalue.ListSizeExact(2)),
				},
			},
		},
	})
}

func TestAccCloudflareR2CustomDomain_JurisdictionEU(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	// zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")
	zoneID := "f8786e3be4a0d180338bb7caf27171a5"
	// domainName := os.Getenv("CLOUDFLARE_DOMAIN")
	domainName := "terraform-r2.cfapi.net"
	resourceName := "cloudflare_r2_custom_domain." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccR2CustomDomainJurisdictionEUConfig(rnd, accountID, zoneID, domainName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("jurisdiction"), knownvalue.StringExact("eu")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("enabled"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("min_tls"), knownvalue.StringExact("1.3")),
				},
			},
		},
	})
}

func TestAccCloudflareR2CustomDomain_Minimal(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	// zoneID := os.Getenv("CLOUDFLARE_ZONE_ID")
	zoneID := "f8786e3be4a0d180338bb7caf27171a5"
	// domainName := os.Getenv("CLOUDFLARE_DOMAIN")
	domainName := "terraform-r2.cfapi.net"
	resourceName := "cloudflare_r2_custom_domain." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccR2CustomDomainMinimalConfig(rnd, accountID, zoneID, domainName),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("enabled"), knownvalue.Bool(false)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("jurisdiction"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("min_tls"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("ciphers"), knownvalue.Null()),
				},
			},
		},
	})
}

func testAccR2CustomDomainConfig(rnd, accountID, zoneID, domainName string) string {
	return acctest.LoadTestCase("basic.tf", rnd, accountID, zoneID, domainName)
}

func testAccR2CustomDomainUpdateConfig(rnd, accountID, zoneID, domainName string) string {
	return acctest.LoadTestCase("update.tf", rnd, accountID, zoneID, domainName)
}

func testAccR2CustomDomainJurisdictionEUConfig(rnd, accountID, zoneID, domainName string) string {
	return acctest.LoadTestCase("jurisdiction_eu.tf", rnd, accountID, zoneID, domainName)
}

func testAccR2CustomDomainMinimalConfig(rnd, accountID, zoneID, domainName string) string {
	return acctest.LoadTestCase("minimal.tf", rnd, accountID, zoneID, domainName)
}
