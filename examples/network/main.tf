# Ejemplo de uso del recurso isard_network

# Crear una red virtual
resource "isard_network" "mi_red" {
  name        = "Red de Prueba"
  description = "Red creada desde Terraform"
  model       = "virtio"
  qos_id      = "unlimited"
}

# Output para ver la información de la red
output "network_info" {
  value = {
    id          = isard_network.mi_red.id
    name        = isard_network.mi_red.name
    metadata_id = isard_network.mi_red.metadata_id
  }
}
