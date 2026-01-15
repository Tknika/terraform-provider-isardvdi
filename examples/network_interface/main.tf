# Ejemplo de uso del recurso isard_network_interface

# Crear una interfaz de red tipo bridge básica
resource "isard_network_interface" "vlan100" {
  id          = "bridge-vlan-100"
  name        = "VLAN 100 Bridge"
  description = "Bridge para VLAN 100 - Red de desarrollo"
  net         = "br-vlan100"
  kind        = "bridge"
  model       = "virtio"
  qos_id      = "unlimited"
}

# Crear interfaz visible para todos los usuarios
resource "isard_network_interface" "public_bridge" {
  id          = "bridge-public"
  name        = "Bridge Público"
  description = "Interfaz visible para todos"
  net         = "br-public"
  kind        = "bridge"
  model       = "virtio"
  qos_id      = "unlimited"
  
  # Listas vacías = visible para todos
  allowed {
    roles      = []
    categories = []
    groups     = []
    users      = []
  }
}

# Crear interfaz solo para administradores
resource "isard_network_interface" "admin_bridge" {
  id          = "bridge-admin"
  name        = "Bridge Administradores"
  description = "Solo para administradores"
  net         = "br-admin"
  kind        = "bridge"
  model       = "virtio"
  qos_id      = "standard"
  
  allowed {
    roles      = ["admin"]
    categories = []
    groups     = []
    users      = []
  }
}

# Crear interfaz para categoría específica
resource "isard_network_interface" "category_bridge" {
  id          = "bridge-marketing"
  name        = "Bridge Marketing"
  description = "Para categoría marketing"
  net         = "br-marketing"
  kind        = "bridge"
  model       = "virtio"
  
  allowed {
    roles      = []
    categories = ["marketing"]  # Reemplazar con ID real
    groups     = []
    users      = []
  }
}

# Output para ver la información de las interfaces
output "interfaces_info" {
  value = {
    vlan100 = {
      id   = isard_network_interface.vlan100.id
      name = isard_network_interface.vlan100.name
    }
    public = {
      id   = isard_network_interface.public_bridge.id
      name = isard_network_interface.public_bridge.name
    }
    admin = {
      id   = isard_network_interface.admin_bridge.id
      name = isard_network_interface.admin_bridge.name
    }
  }
}
