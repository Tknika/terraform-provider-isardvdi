---
page_title: "isard_network_interfaces Data Source - terraform-provider-isard"
subcategory: ""
description: |-
  Retrieve network interfaces from Isard VDI.
---

# Data Source: isard_network_interfaces

Obtiene la lista de interfaces de red del sistema disponibles en Isard VDI. **Requiere privilegios de administrador.**

## Ejemplo de Uso

### Obtener Todas las Interfaces

```hcl
data "isard_network_interfaces" "all" {}

output "todas_las_interfaces" {
  value = data.isard_network_interfaces.all.interfaces
}
```

### Filtrar por Tipo de Interfaz

```hcl
data "isard_network_interfaces" "bridges" {
  filter = {
    kind = "bridge"
  }
}

output "interfaces_bridge" {
  value = data.isard_network_interfaces.bridges.interfaces
}
```

### Filtrar por Nombre

```hcl
data "isard_network_interfaces" "wireguard" {
  filter = {
    name = "wireguard"
  }
}

# Usar en una VM
resource "isard_vm" "con_vpn" {
  name        = "vm-con-vpn"
  template_id = data.isard_templates.ubuntu.templates[0].id
  
  interfaces = [
    data.isard_network_interfaces.wireguard.interfaces[0].id
  ]
}
```

### Filtrar por Red

```hcl
data "isard_network_interfaces" "public_bridge" {
  filter = {
    net = "br-public"
  }
}

output "interfaces_publicas" {
  value = data.isard_network_interfaces.public_bridge.interfaces
}
```

### Múltiples Filtros

```hcl
data "isard_network_interfaces" "ovs_bridges" {
  filter = {
    kind = "ovs"
    name = "prod"  # Búsqueda parcial case-insensitive
  }
}
```

### Verificar Disponibilidad

```hcl
data "isard_network_interfaces" "custom" {
  filter = {
    name = "custom-bridge"
  }
}

locals {
  interface_exists = length(data.isard_network_interfaces.custom.interfaces) > 0
  interface_id     = local.interface_exists ? data.isard_network_interfaces.custom.interfaces[0].id : null
}

resource "isard_vm" "conditional" {
  count = local.interface_exists ? 1 : 0
  
  name        = "vm-conditional"
  template_id = data.isard_templates.ubuntu.templates[0].id
  interfaces  = [local.interface_id]
}
```

### Búsqueda por Nombre (Deprecated)

```hcl
# Forma antigua (aún soportada)
data "isard_network_interfaces" "wireguard_old" {
  name = "wireguard"  # Búsqueda exacta
}

# Forma nueva (recomendada)
data "isard_network_interfaces" "wireguard_new" {
  filter = {
    name = "wireguard"  # Búsqueda parcial
  }
}
```

## Argumentos

Los siguientes argumentos son soportados:

### Opcionales

- `name` - (Opcional, Deprecated) Nombre exacto de la interfaz a buscar. Se recomienda usar `filter.name` en su lugar.
- `filter` - (Opcional) Bloque de filtros para búsqueda avanzada:
  - `name` - (Opcional) Búsqueda parcial del nombre (case-insensitive). Ejemplo: `"bridge"` encontrará `"test-bridge-1"`, `"Bridge Public"`, etc.
  - `kind` - (Opcional) Tipo de interfaz. Valores: `"bridge"`, `"network"`, `"ovs"`, `"personal"`.
  - `net` - (Opcional) Red/bridge del sistema (búsqueda exacta).

## Atributos Exportados

- `id` - ID del data source.
- `interfaces` - Lista de objetos interfaz. Cada objeto contiene:
  - `id` - ID único de la interfaz.
  - `name` - Nombre descriptivo de la interfaz.
  - `description` - Descripción de la interfaz.
  - `net` - Red/bridge del sistema asociado.
  - `kind` - Tipo de interfaz (`bridge`, `network`, `ovs`, `personal`).
  - `model` - Modelo de dispositivo de red (`virtio`, `e1000`, `rtl8139`).
  - `qos_id` - ID del perfil QoS aplicado.

## Filtrado

### Por Tipo (`kind`)

Los tipos disponibles son:
- `bridge` - Bridges Linux estándar
- `network` - Redes libvirt
- `ovs` - Open vSwitch
- `personal` - Redes personales con VLAN

```hcl
data "isard_network_interfaces" "by_type" {
  filter = {
    kind = "bridge"
  }
}
```

### Por Nombre (`name`)

La búsqueda de nombre es:
- **Parcial**: Busca coincidencias en cualquier parte del nombre
- **Case-insensitive**: No distingue mayúsculas de minúsculas

```hcl
# Encuentra: "test-bridge", "bridge-public", "My Bridge"
data "isard_network_interfaces" "with_bridge" {
  filter = {
    name = "bridge"
  }
}
```

### Por Red (`net`)

La búsqueda de red es exacta:

```hcl
data "isard_network_interfaces" "br_public" {
  filter = {
    net = "br-public"
  }
}
```

## Casos de Uso

### 1. Listar Todas las Interfaces Disponibles

```hcl
data "isard_network_interfaces" "all" {}

output "inventory" {
  value = {
    for iface in data.isard_network_interfaces.all.interfaces :
    iface.id => {
      name = iface.name
      type = iface.kind
      net  = iface.net
    }
  }
}
```

### 2. Validar Interfaz Requerida

```hcl
data "isard_network_interfaces" "required" {
  filter = {
    name = "wireguard"
  }
}

resource "null_resource" "validate" {
  count = length(data.isard_network_interfaces.required.interfaces) > 0 ? 0 : 1
  
  provisioner "local-exec" {
    command = "echo 'ERROR: wireguard interface not found' && exit 1"
  }
}
```

### 3. Asignar Dinámicamente Interfaces a VMs

```hcl
data "isard_network_interfaces" "available" {
  filter = {
    kind = "bridge"
  }
}

resource "isard_vm" "team" {
  count = 3
  
  name        = "vm-team-${count.index + 1}"
  template_id = data.isard_templates.ubuntu.templates[0].id
  
  interfaces = [
    "wireguard",
    data.isard_network_interfaces.available.interfaces[count.index % length(data.isard_network_interfaces.available.interfaces)].id
  ]
}
```

### 4. Documentar Infraestructura

```hcl
data "isard_network_interfaces" "all" {}

output "network_infrastructure" {
  description = "Inventario completo de interfaces de red"
  value = {
    total_interfaces = length(data.isard_network_interfaces.all.interfaces)
    by_type = {
      for kind in distinct([for i in data.isard_network_interfaces.all.interfaces : i.kind]) :
      kind => [
        for i in data.isard_network_interfaces.all.interfaces :
        i.name if i.kind == kind
      ]
    }
  }
}
```

## Notas Importantes

- **Solo administradores** pueden listar interfaces de red del sistema
- Los filtros se combinan con lógica AND (todos deben cumplirse)
- La búsqueda por nombre en `filter` es más flexible que el parámetro `name` directo
- Si no se especifican filtros, devuelve todas las interfaces
- Las interfaces con `allowed.roles` restrictivos también aparecen en la lista

## Relación con Otros Recursos

Este data source es útil para:
- Seleccionar interfaces para VMs (`isard_vm`)
- Validar existencia antes de crear interfaces (`isard_network_interface`)
- Auditar la configuración de red del sistema
- Generar documentación automática de la infraestructura
