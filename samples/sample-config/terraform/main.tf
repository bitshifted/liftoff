
# Keycloak server with public IP address
resource "hcloud_server" "keycloak-server" {
   name = "keycloak"
  image = "ubuntu-24.04"
  server_type = "cx22"
  location = "fsn1"
  ssh_keys = [ data.hcloud_ssh_key.server_ssh_key.name ]
  user_data = "${file("cloud-config.yaml")}"
}
