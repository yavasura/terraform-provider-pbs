resource "pbs_namespace" "production" {
  store     = "backups"
  namespace = "production"
}

resource "pbs_namespace" "production_vms" {
  store      = "backups"
  namespace  = "production/vms"
  depends_on = [pbs_namespace.production]
}
