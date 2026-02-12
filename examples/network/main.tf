# Ejemplo de uso del recurso isard_network

# Crear una red virtual
resource "isardvdi_network" "mi_red" {
  name        = "Red de Prueba"
  description = "Red creada desde Terraform"
  model       = "virtio"
  qos_id      = "unlimited"
}

# Output para ver la información de la red
output "network_info" {
  value = {
    id          = isardvdi_network.mi_red.id
    name        = isardvdi_network.mi_red.name
    metadata_id = isardvdi_network.mi_red.metadata_id
  }
}
