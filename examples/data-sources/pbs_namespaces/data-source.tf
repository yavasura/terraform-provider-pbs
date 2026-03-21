data "pbs_namespaces" "production" {
  store     = "backups"
  prefix    = "production"
  max_depth = 2
}

output "production_namespaces" {
  value = data.pbs_namespaces.production.namespaces
}
