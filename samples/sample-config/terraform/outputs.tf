
output "server_ip_address" {
  value = hcloud_server.keycloak-server.ipv4_address
}