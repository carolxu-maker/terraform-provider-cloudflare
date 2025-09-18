resource "cloudflare_r2_bucket" "%[1]s" {
  account_id = "%[2]s"
  name       = "%[1]s"
}

resource "cloudflare_r2_bucket_lock" "%[1]s" {
  account_id  = "%[2]s"
  bucket_name = "%[1]s"

  rules = [{
    id      = "minimal-lock"
    enabled = false
    condition = {
      type            = "Age"
      max_age_seconds = 3600
    }
  }]
}
