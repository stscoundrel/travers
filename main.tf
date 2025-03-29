terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
  credentials = file("service-account.json")
}

# Cloud Function
resource "google_cloudfunctions2_function" "event_fetcher" {
  name        = "event-fetcher"
  location    = var.region
  description = "Fetches and stores new events"
  depends_on = [
    google_storage_bucket_object.function_zip,
    google_storage_bucket_iam_member.cloud_function_gcs_access
  ]

  build_config {
    runtime     = "go122"
    entry_point = "FetchAndStoreEvents"
    source {
      storage_source {
        bucket = google_storage_bucket.source_bucket.name
        object = google_storage_bucket_object.function_zip.name
      }
    }
  }

  service_config {
    max_instance_count = 1
    available_memory   = "256Mi"
    timeout_seconds    = 60
    service_account_email = google_service_account.service_account.email
    environment_variables = {
      EVENT_STORAGE_BUCKET = google_storage_bucket.event_data_bucket.name
    }
  }
}

# Service account to call the function via scheduler.
resource "google_service_account" "service_account" {
  account_id   = "cloud-function-invoker"
  display_name = "Invoker service account"
}

# IAM Policies to Allow Cloud Scheduler to invoke the function
resource "google_cloudfunctions2_function_iam_member" "invoke" {
  cloud_function = google_cloudfunctions2_function.event_fetcher.name
  role           = "roles/cloudfunctions.invoker"
  member         = "serviceAccount:${google_service_account.service_account.email}"
}

resource "google_cloud_run_service_iam_member" "cloud_run_invoker" {
  service  = google_cloudfunctions2_function.event_fetcher.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.service_account.email}"
}

# Allow Cloud Scheduler to use the service account
resource "google_service_account_iam_member" "allow_scheduler" {
  service_account_id = google_service_account.service_account.id
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.service_account.email}"
}

# Cloud Scheduler - Triggers Function Daily
resource "google_cloud_scheduler_job" "daily_event_fetch" {
  name        = "daily-event-fetch"
  description = "Runs event fetcher daily"
  schedule    = "45 16 * * *" # Runs at given UTC every day
  time_zone   = "UTC"
  depends_on = [
    google_cloudfunctions2_function.event_fetcher,
    google_cloudfunctions2_function_iam_member.invoke,
    google_cloud_run_service_iam_member.cloud_run_invoker,
    google_service_account_iam_member.allow_scheduler,
  ]

  http_target {
    uri         = google_cloudfunctions2_function.event_fetcher.service_config[0].uri
    http_method = "GET"

    oidc_token {
      audience              = google_cloudfunctions2_function.event_fetcher.service_config[0].uri
      service_account_email = google_service_account.service_account.email
    }
  }
  
}

# Storage for function source code
resource "google_storage_bucket" "source_bucket" {
  name          = "${var.project_id}-function-source"
  location      = var.region
  force_destroy = true
}

resource "google_storage_bucket_object" "function_zip" {
  name   = "function-source.zip"
  bucket = google_storage_bucket.source_bucket.name
  source = "function-source.zip"
  depends_on = [google_storage_bucket.source_bucket]
}

# Bucket to store event data.
resource "google_storage_bucket" "event_data_bucket" {
  name          = "${var.project_id}-events"
  location      = var.region
  force_destroy = true
}

# Access bucket with service account.
resource "google_storage_bucket_iam_member" "cloud_function_gcs_access" {
  bucket = google_storage_bucket.event_data_bucket.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.service_account.email}"
}


# Notification condition.
resource "google_logging_metric" "event_alert_metric" {
  name   = "event_alert_metric"
  filter = "textPayload =~ \"Travers found new events\""
  metric_descriptor {
    metric_kind = "DELTA"
    value_type  = "INT64"
  }
}

# Notification channel.
resource "google_monitoring_notification_channel" "email_alert" {
  display_name = "Travers Notification"
  type         = "email"
  labels = {
    email_address = var.alert_email
  }
}

# Notification policy.
resource "google_monitoring_alert_policy" "event_alert_policy" {
  display_name = "Travers Event Alert"
  documentation {
    subject = "Travers found new events"
  }
  combiner     = "OR"
  depends_on = [google_logging_metric.event_alert_metric]

  conditions {
    display_name = "Event Log Condition"
    condition_threshold {
      # NOTE: following line does not work correctly. The resource type should be omitted, but API required it. 
      # As temporary workaround, remove the resource.type in console after deploy.
      filter          = "resource.type=\"metric\" AND metric.type=\"logging.googleapis.com/user/event_alert_metric\""
      duration        = "0s"
      comparison      = "COMPARISON_GT"
      threshold_value = 0
      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "ALIGN_COUNT"
      }
      trigger {
        count = 1
      }
    }
  }

  notification_channels = [google_monitoring_notification_channel.email_alert.id]
}

# Fetch project data dynamically to reference project_number
data "google_project" "project" {}
