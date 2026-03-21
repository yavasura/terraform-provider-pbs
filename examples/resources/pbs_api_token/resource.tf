resource "pbs_user" "automation" {
  userid  = "automation@pbs"
  comment = "Managed account for automation"
  enable  = true
}

resource "pbs_api_token" "terraform" {
  userid     = pbs_user.automation.userid
  token_name = "terraform"
  comment    = "Managed API token for Terraform"
  enable     = true
}
