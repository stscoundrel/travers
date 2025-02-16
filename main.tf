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
  depends_on = [google_storage_bucket_object.function_zip]

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
  schedule    = "20 14 * * *" # Runs at given UTC every day
  time_zone   = "UTC"

  http_target {
    uri         = google_cloudfunctions2_function.event_fetcher.service_config[0].uri
    http_method = "GET"

    oidc_token {
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

# Fetch project data dynamically to reference project_number
data "google_project" "project" {}
