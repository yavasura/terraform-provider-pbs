resource "pbs_user" "backup_operator" {
  userid    = "backup-operator@ldap"
  comment   = "Managed account for backup operations"
  enable    = true
  firstname = "Backup"
  lastname  = "Operator"
  email     = "backup-operator@example.com"
}
