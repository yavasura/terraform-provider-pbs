resource "pbs_acl" "backup_admin" {
  path      = "/datastore/backups"
  ugid      = "backup-operator@pbs"
  role_id   = "DatastoreAdmin"
  propagate = true
}
