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
        max_age = 60
      }
    }
  }, {
    id      = "archive-objects"
    enabled = true
    conditions = {
      prefix = "archive/"
    }
    storage_class_transitions = {
      condition = {
        type    = "Age"
        max_age = 7
      }
      storage_class = "InfrequentAccess"
    }
  }, {
    id      = "cleanup-multipart"
    enabled = true
    conditions = {
      prefix = ""
    }
    abort_multipart_uploads_transition = {
      condition = {
        type    = "Age"
        max_age = 1
      }
    }
  }]
}
