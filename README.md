# Travers

Monitor certain firearms exercises on a certain island.

Work in progress for personal use.

### Deploy

To hook up the cloud infra with Terraform:

- Manually zip latest src version to `function-source.zip`. To be automated.
- Populate `service-account.json` with credentials from Cloud Console.
- `terraform plan` /  `terraform apply`
  
Currently Terraform host state locally.

Presently creates:
- Cloud function to run Travers
- Event Scheduler action to trigger it daily
- Storage bucket for event data storage
- Required storage buckets & permissions

To be added:
- Messaging on new events.

### Whats in the name?

One meaning of _Travers_ is _someone who lives at a crossing place_. To get to activities which this Travers monitors, you'll have to cross body of water and some other things.
