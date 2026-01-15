# Ejemplo de uso del recurso isard_network_interface

# Crear una interfaz de red tipo bridge
resource "isard_network_interface" "vlan100" {
  id          = "bridge-vlan-100"
  name        = "VLAN 100 Bridge"
  description = "Bridge para VLAN 100 - Red de desarrollo"
  kind        = "bridge"
  model       = "virtio"
  qos_id      = "unlimited"
}

# Crear otra interfaz para VLAN 200
resource "isard_network_interface" "vlan200" {
  id          = "bridge-vlan-200"
  name        = "VLAN 200 Bridge"
  description = "Bridge para VLAN 200 - Red de producción"
  kind        = "bridge"
  model       = "virtio"
  qos_id      = "unlimited"
}

# Output para ver la información de las interfaces
output "interfaces_info" {
  value = {
    vlan100 = {
      id   = isard_network_interface.vlan100.id
      name = isard_network_interface.vlan100.name
    }
    vlan200 = {
      id   = isard_network_interface.vlan200.id
      name = isard_network_interface.vlan200.name
    }
  }
}
