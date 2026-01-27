# Ejemplo de uso del recurso isard_network_interface

# Crear una interfaz de red tipo bridge básica
resource "isard_network_interface" "vlan100" {
  id          = "bridge-vlan-100"
  name        = "VLAN 100 Bridge"
  description = "Bridge para VLAN 100 - Red de desarrollo"
  net         = "br-vlan100"
  kind        = "bridge"
  ifname      = "interface"
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
  ifname      = "interface"
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
  ifname      = "interface"
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
  ifname      = "interface"
  model       = "virtio"
  
  allowed {
    roles      = []
    categories = ["marketing"]  # Reemplazar con ID real
    groups     = []
    users      = []
  }
}

# Ejemplos de diferentes tipos de kind

# Interfaz tipo OVS
resource "isard_network_interface" "ovs_example" {
  id          = "ovs-vlan-2056"
  name        = "OVS VLAN 2056"
  description = "Interfaz OVS con VLAN 2056"
  kind        = "ovs"
  ifname      = "2056"
  model       = "virtio"
  qos_id      = "unlimited"
  net         = "2056"
}

# Interfaz tipo network
resource "isard_network_interface" "network_example" {
  id          = "network1-interface"
  name        = "Network 1"
  description = "Red libvirt"
  kind        = "network"
  ifname      = "network1"
  model       = "virtio"
  qos_id      = "unlimited"
  net         = "network1"
}

# Interfaz tipo personal (rango VLAN)
resource "isard_network_interface" "personal_example" {
  id          = "personal-vlan-range"
  name        = "Personal VLAN Range"
  description = "Rango VLAN personal 1000-1500"
  kind        = "personal"
  ifname      = "1000-1500"
  model       = "virtio"
  qos_id      = "unlimited"
  net         = "1000-1500"
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
    ovs = {
      id   = isard_network_interface.ovs_example.id
      name = isard_network_interface.ovs_example.name
      kind = isard_network_interface.ovs_example.kind
    }
    network = {
      id   = isard_network_interface.network_example.id
      name = isard_network_interface.network_example.name
      kind = isard_network_interface.network_example.kind
    }
    personal = {
      id   = isard_network_interface.personal_example.id
      name = isard_network_interface.personal_example.name
      kind = isard_network_interface.personal_example.kind
    }
  }
}
