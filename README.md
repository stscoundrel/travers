# Travers

Monitor certain firearms exercises on a certain island. When new events are detected, sends alert to sign up.

### Why though?

For personal use to avoid manual hassle of checking when new signup become available. The events it monitors tend to be popular and available slots limited, so it helps to know when new ones are published.

### Deploy

To hook up the cloud infra with Terraform:

- Manually zip latest src version to `function-source.zip`.
- Populate `service-account.json` with credentials from Cloud Console.
- `terraform plan` /  `terraform apply`
  
Currently Terraform host state locally.

Presently creates:
- Cloud function to run Travers
- Event Scheduler action to trigger it daily
- Storage bucket for event data storage
- Required storage buckets & permissions
- Messaging events for alerts.

### Whats in the name?

One meaning of *Travers* is "_someone who lives at a crossing place_". To get to activities which this Travers monitors, you'll have to cross body of water and some other things.
