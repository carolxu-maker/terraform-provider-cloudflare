resource "cloudflare_r2_bucket" "%[1]s" {
  account_id = "%[2]s"
  name       = "%[1]s"
}

resource "cloudflare_r2_bucket_lifecycle" "%[1]s" {
  account_id  = "%[2]s"
  bucket_name = "%[1]s"

  rules = [{
    id      = "delete-old-objects"
    enabled = true
    conditions = {
      prefix = "logs/"
    }
    delete_objects_transition = {
      condition = {
        type    = "Age"
        max_age = 30
      }
    }
  }]
}