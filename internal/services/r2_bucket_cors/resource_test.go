package r2_bucket_cors_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/acctest"
	"github.com/cloudflare/terraform-provider-cloudflare/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccCloudflareR2BucketCors_Basic(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	resourceName := "cloudflare_r2_bucket_cors." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccR2BucketCorsConfig(rnd, accountID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("bucket_name"), knownvalue.StringExact(rnd)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("jurisdiction"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("id"), knownvalue.StringExact("rule1")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("allowed").AtMapKey("methods"), knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("GET"), knownvalue.StringExact("POST")})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("allowed").AtMapKey("origins"), knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("https://example.com")})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("allowed").AtMapKey("headers"), knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("Content-Type")})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("expose_headers"), knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("ETag")})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("max_age_seconds"), knownvalue.Float64Exact(3600)),
				},
			},
		},
	})
}

func TestAccCloudflareR2BucketCors_Update(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	resourceName := "cloudflare_r2_bucket_cors." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccR2BucketCorsConfig(rnd, accountID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("max_age_seconds"), knownvalue.Float64Exact(3600)),
				},
			},
			{
				Config: testAccR2BucketCorsUpdateConfig(rnd, accountID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules"), knownvalue.ListSizeExact(2)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("max_age_seconds"), knownvalue.Float64Exact(7200)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("allowed").AtMapKey("methods"), knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("GET"), knownvalue.StringExact("POST"), knownvalue.StringExact("PUT")})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(1).AtMapKey("id"), knownvalue.StringExact("rule2")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(1).AtMapKey("allowed").AtMapKey("methods"), knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("DELETE"), knownvalue.StringExact("HEAD")})),
				},
			},
		},
	})
}

func TestAccCloudflareR2BucketCors_JurisdictionEU(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	resourceName := "cloudflare_r2_bucket_cors." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccR2BucketCorsJurisdictionEUConfig(rnd, accountID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("bucket_name"), knownvalue.StringExact(rnd)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("jurisdiction"), knownvalue.StringExact("eu")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("id"), knownvalue.StringExact("eu-rule")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("allowed").AtMapKey("origins"), knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("https://eu.example.com")})),
				},
			},
		},
	})
}

func TestAccCloudflareR2BucketCors_RequiredOnly(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	resourceName := "cloudflare_r2_bucket_cors." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccR2BucketCorsRequiredOnlyConfig(rnd, accountID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("bucket_name"), knownvalue.StringExact(rnd)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("jurisdiction"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules"), knownvalue.Null()),
				},
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdPrefix:                  fmt.Sprintf("accounts/%s/", accountID),
				ImportStateVerifyIdentifierAttribute: "bucket_name",
			},
		},
	})
}

func TestAccCloudflareR2BucketCors_AllHttpMethods(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	resourceName := "cloudflare_r2_bucket_cors." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccR2BucketCorsAllHttpMethodsConfig(rnd, accountID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("bucket_name"), knownvalue.StringExact(rnd)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("id"), knownvalue.StringExact("all-methods-rule")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("allowed").AtMapKey("methods"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("GET"),
						knownvalue.StringExact("PUT"),
						knownvalue.StringExact("POST"),
						knownvalue.StringExact("DELETE"),
						knownvalue.StringExact("HEAD"),
					})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("allowed").AtMapKey("origins"), knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("https://example.com")})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("allowed").AtMapKey("headers"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("Content-Type"),
						knownvalue.StringExact("Authorization"),
					})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("expose_headers"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("ETag"),
						knownvalue.StringExact("Content-Length"),
					})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("max_age_seconds"), knownvalue.Float64Exact(3600)),
				},
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdPrefix:                  fmt.Sprintf("accounts/%s/", accountID),
				ImportStateVerifyIdentifierAttribute: "bucket_name",
			},
		},
	})
}

func TestAccCloudflareR2BucketCors_MinimalRules(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	resourceName := "cloudflare_r2_bucket_cors." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccR2BucketCorsMinimalRulesConfig(rnd, accountID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("bucket_name"), knownvalue.StringExact(rnd)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("allowed").AtMapKey("methods"), knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("GET")})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("allowed").AtMapKey("origins"), knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("*")})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("allowed").AtMapKey("headers"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("id"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("expose_headers"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("max_age_seconds"), knownvalue.Null()),
				},
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdPrefix:                  fmt.Sprintf("accounts/%s/", accountID),
				ImportStateVerifyIdentifierAttribute: "bucket_name",
			},
		},
	})
}

func TestAccCloudflareR2BucketCors_CompleteLifecycle(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	resourceName := "cloudflare_r2_bucket_cors." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccR2BucketCorsLifecycleMinimalConfig(rnd, accountID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("allowed").AtMapKey("methods"), knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("GET")})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("allowed").AtMapKey("origins"), knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("https://minimal.com")})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("id"), knownvalue.Null()),
				},
			},
			{
				Config: testAccR2BucketCorsLifecycleFullConfig(rnd, accountID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules"), knownvalue.ListSizeExact(2)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("id"), knownvalue.StringExact("full-rule")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("allowed").AtMapKey("methods"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("GET"),
						knownvalue.StringExact("POST"),
						knownvalue.StringExact("PUT"),
					})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("allowed").AtMapKey("origins"), knownvalue.SetExact([]knownvalue.Check{
						knownvalue.StringExact("https://full.com"),
						knownvalue.StringExact("https://example.com"),
					})),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("max_age_seconds"), knownvalue.Float64Exact(7200)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(1).AtMapKey("id"), knownvalue.StringExact("admin-rule")),
				},
			},
			{
				Config: testAccR2BucketCorsLifecycleMinimalConfig(rnd, accountID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("allowed").AtMapKey("methods"), knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("GET")})),
				},
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdPrefix:                  fmt.Sprintf("accounts/%s/", accountID),
				ImportStateVerifyIdentifierAttribute: "bucket_name",
			},
		},
	})
}

func TestAccCloudflareR2BucketCors_RemoveAllRules(t *testing.T) {
	rnd := utils.GenerateRandomResourceName()
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	resourceName := "cloudflare_r2_bucket_cors." + rnd

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccR2BucketCorsConfig(rnd, accountID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules"), knownvalue.ListSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("id"), knownvalue.StringExact("rule1")),
				},
			},
			{
				Config: testAccR2BucketCorsRemoveRulesConfig(rnd, accountID),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("account_id"), knownvalue.StringExact(accountID)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("bucket_name"), knownvalue.StringExact(rnd)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("jurisdiction"), knownvalue.StringExact("default")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rules"), knownvalue.Null()),
				},
			},
			{
				ResourceName:                         resourceName,
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdPrefix:                  fmt.Sprintf("accounts/%s/", accountID),
				ImportStateVerifyIdentifierAttribute: "bucket_name",
			},
		},
	})
}

func testAccR2BucketCorsConfig(rnd, accountID string) string {
	return acctest.LoadTestCase("r2bucketcorsbasic.tf", rnd, accountID)
}

func testAccR2BucketCorsUpdateConfig(rnd, accountID string) string {
	return acctest.LoadTestCase("r2bucketcorsupdate.tf", rnd, accountID)
}

func testAccR2BucketCorsJurisdictionEUConfig(rnd, accountID string) string {
	return acctest.LoadTestCase("r2bucketcorsjurisdiction_eu.tf", rnd, accountID)
}

func testAccR2BucketCorsRequiredOnlyConfig(rnd, accountID string) string {
	return acctest.LoadTestCase("r2bucketcorsrequired_only.tf", rnd, accountID)
}

func testAccR2BucketCorsAllHttpMethodsConfig(rnd, accountID string) string {
	return acctest.LoadTestCase("r2bucketcorsall_http_methods.tf", rnd, accountID)
}

func testAccR2BucketCorsMinimalRulesConfig(rnd, accountID string) string {
	return acctest.LoadTestCase("r2bucketcorsminimal_rules.tf", rnd, accountID)
}

func testAccR2BucketCorsLifecycleMinimalConfig(rnd, accountID string) string {
	return acctest.LoadTestCase("r2bucketcorslifecycle_minimal.tf", rnd, accountID)
}

func testAccR2BucketCorsLifecycleFullConfig(rnd, accountID string) string {
	return acctest.LoadTestCase("r2bucketcorslifecycle_full.tf", rnd, accountID)
}

func testAccR2BucketCorsRemoveRulesConfig(rnd, accountID string) string {
	return acctest.LoadTestCase("r2bucketcorsremove_rules.tf", rnd, accountID)
}
