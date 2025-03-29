variable "project_id" {
  description = "Project ID"
  type        = string
  default     = "travers-451115"
}

variable "alert_email" {
  description = "Email address to receive new event alerts"
  type        = string
}

variable "region" {
  description = "Google Cloud region"
  type        = string
  default     = "us-central1"
}