---
page_title: "isard_network_interface Resource - terraform-provider-isard"
subcategory: ""
description: |-
  Manages a network interface in Isard VDI.
---

# Resource: isard_network_interface

Gestiona una interfaz de red a nivel de sistema en Isard VDI. **Requiere privilegios de administrador.**

## Ejemplo de Uso

### Ejemplo Básico

```hcl
resource "isard_network_interface" "bridge_dev" {
  id          = "bridge-desarrollo"
  name        = "Bridge Desarrollo"
  description = "Interfaz bridge para desarrollo"
  net         = "br-public"
  kind        = "bridge"
  ifname      = "interface"
  model       = "virtio"
  qos_id      = "unlimited"
}
```

### Interfaz OVS con VLAN

```hcl
resource "isard_network_interface" "ovs_vlan100" {
  id          = "ovs-vlan-100"
  name        = "OVS VLAN 100"
  description = "Interfaz OVS con VLAN 100"
  net         = "100"
  kind        = "ovs"
  ifname      = "100"
  model       = "virtio"
  qos_id      = "standard"
}
```

### Interfaz de Red Personal

```hcl
resource "isard_network_interface" "personal_range" {
  id          = "personal-dev-team"
  name        = "Red Personal Equipo Dev"
  description = "Rango VLAN personal para equipo"
  net         = "200-210"
  kind        = "personal"
  ifname      = "200-210"
  model       = "virtio"
  qos_id      = "unlimited"
}
```

### Interfaz Visible para Todos

```hcl
resource "isard_network_interface" "public_bridge" {
  id          = "bridge-public"
  name        = "Bridge Público"
  description = "Interfaz visible para todos los usuarios"
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
```

### Interfaz con Permisos Específicos

```hcl
resource "isard_network_interface" "restricted_bridge" {
  id          = "bridge-restricted"
  name        = "Bridge Restringido"
  description = "Solo para roles admin y manager"
  net         = "br-restricted"
  kind        = "bridge"
  ifname      = "interface"
  model       = "virtio"
  qos_id      = "standard"
  
  allowed {
    roles      = ["admin", "manager"]
    categories = []  # Sin restricciones por categoría
    groups     = []  # Sin restricciones por grupo
    users      = []  # Sin restricciones por usuario
  }
}
```

### Interfaz para Categoría Específica

```hcl
resource "isard_network_interface" "category_bridge" {
  id          = "bridge-marketing"
  name        = "Bridge Marketing"
  description = "Solo para la categoría marketing"
  net         = "br-marketing"
  kind        = "bridge"
  ifname      = "interface"
  model       = "virtio"
  
  allowed {
    roles      = []
    categories = ["marketing-category-id"]
    groups     = []
    users      = []
  }
}
```

### Usar con VM

```hcl
# Crear la interfaz
resource "isard_network_interface" "custom_bridge" {
  id          = "custom-bridge-1"
  name        = "Custom Bridge 1"
  description = "Bridge personalizado"
  net         = "br-custom"
  kind        = "bridge"
  ifname      = "interface"
  model       = "virtio"
  qos_id      = "unlimited"
}

# Usar en una VM
resource "isard_vm" "con_interfaz" {
  name        = "vm-con-interfaz-custom"
  description = "VM con interfaz personalizada"
  template_id = data.isard_templates.ubuntu.templates[0].id
  
  interfaces = [
    "wireguard",  # Requerido para RDP viewers
    isard_network_interface.custom_bridge.id
  ]
}
```

## Argumentos

Los siguientes argumentos son soportados:

### Requeridos

- `id` - (Requerido) ID único de la interfaz. No se puede modificar después de creada.
- `name` - (Requerido) Nombre descriptivo de la interfaz.
- `net` - (Requerido) Especificación de la red/bridge del sistema:
  - Para `kind = "bridge"` o `"network"`: nombre del bridge (ej: `"br-public"`)
  - Para `kind = "ovs"`: número de VLAN o `"4095"` para wireguard
  - Para `kind = "personal"`: rango de VLANs (ej: `"200-210"`)

### Opcionales

- `description` - (Opcional) Descripción de la interfaz.
- `kind` - (Opcional, Computed) Tipo de interfaz. Por defecto según template. Valores:
  - `"bridge"` - Bridge Linux estándar
  - `"network"` - Red libvirt
  - `"ovs"` - Open vSwitch
  - `"personal"` - Red personal con rango VLAN
- `ifname` - (Opcional) Nombre de la interfaz física/virtual. Típicamente coincide con `net` o representa la VLAN/rango.
- `model` - (Opcional, Computed) Modelo de dispositivo de red. Por defecto: `"virtio"`. Valores: `"virtio"`, `"e1000"`, `"rtl8139"`.
- `qos_id` - (Opcional, Computed) ID del perfil QoS de red. Por defecto: `"unlimited"`.
- `allowed` - (Opcional) Bloque de permisos de acceso. Si se omite, la interfaz no tendrá restricciones específicas.
  - `roles` - (Opcional) Lista de IDs de roles permitidos. Lista vacía `[]` = todos los roles pueden usar la interfaz.
  - `categories` - (Opcional) Lista de IDs de categorías permitidas. Lista vacía `[]` = todas las categorías.
  - `groups` - (Opcional) Lista de IDs de grupos permitidos. Lista vacía `[]` = todos los grupos.
  - `users` - (Opcional) Lista de IDs de usuarios permitidos. Lista vacía `[]` = todos los usuarios.

**Nota sobre permisos:** Para hacer una interfaz visible a todos los usuarios, use el bloque `allowed` con todas las listas vacías (`[]`). Si omite el bloque `allowed` completamente, la interfaz seguirá siendo accesible pero sin definición explícita de permisos.

## Atributos Exportados

Los mismos que los argumentos, ya que todos son configurables y computed.

## Import

Las interfaces de red pueden ser importadas usando su ID:

```bash
terraform import isard_network_interface.bridge_dev bridge-desarrollo
```

## Ciclo de Vida

### Create

Al crear una interfaz:
1. Se valida que el usuario tenga privilegios de administrador
2. Se crea usando `POST /api/v3/admin/table/add/interfaces`
3. Se verifica la creación exitosa

### Read

Al leer una interfaz:
1. Se obtiene la información desde `GET /api/v3/admin/table/interfaces`
2. Se busca la interfaz específica por ID

### Update

Al actualizar una interfaz:
1. Se envían solo los campos modificados usando `PUT /api/v3/admin/table/update/interfaces`
2. Se releen los valores actualizados

### Delete

Al eliminar una interfaz:
1. Se elimina usando `DELETE /api/v3/admin/table/interfaces/{id}`
2. Se desasignan las VMs que usen esta interfaz

## Notas Importantes

- **Solo administradores** pueden gestionar interfaces de red del sistema
- El campo `net` es crítico y debe corresponder con la infraestructura de red real
- Al eliminar una interfaz, se desasigna de todas las VMs que la usen
- Para VMs con viewers RDP, es obligatorio incluir la interfaz `wireguard`
- Los valores `kind`, `model` y `qos_id` tienen valores por defecto del sistema si no se especifican

## Control de Acceso con `allowed`

El bloque `allowed` controla qué usuarios pueden ver y usar la interfaz de red:

### Hacer la Interfaz Visible para Todos

Para que **todos los usuarios** puedan usar la interfaz:

```hcl
allowed {
  roles      = []
  categories = []
  groups     = []
  users      = []
}
```

### Restringir por Rol

Para limitar a roles específicos:

```hcl
allowed {
  roles      = ["admin", "manager", "advanced"]
  categories = []  # Sin restricción adicional
  groups     = []
  users      = []
}
```

### Restringir por Categoría

Para limitar a categorías específicas:

```hcl
allowed {
  roles      = []
  categories = ["categoria-dev", "categoria-prod"]
  groups     = []
  users      = []
}
```

### Combinación de Restricciones

Puede combinar restricciones. El usuario debe cumplir AL MENOS UNA de las restricciones no vacías:

```hcl
allowed {
  roles      = ["admin"]               # Administradores
  categories = ["categoria-especial"]  # O usuarios de esta categoría
  groups     = []
  users      = ["user-id-especifico"]  # O este usuario específico
}
```

### Comportamiento de Permisos

- **Lista vacía (`[]`)**: Sin restricción en ese nivel (permite a todos)
- **Lista con IDs**: Solo esos IDs específicos tienen acceso
- **Campo omitido**: Sin restricción explícita en ese nivel

**Importante**: Las restricciones se evalúan con lógica OR entre tipos. Si un usuario cumple **cualquiera** de las restricciones definidas (rol, categoría, grupo o usuario), tendrá acceso.

## Tipos de Interfaz

### Bridge (`kind = "bridge"`)
Conecta VMs a bridges Linux estándar. Útil para redes locales simples.

### Network (`kind = "network"`)
Usa redes definidas en libvirt. Ofrece más opciones de configuración.

### OVS (`kind = "ovs"`)
Usa Open vSwitch para networking avanzado con soporte VLAN.

### Personal (`kind = "personal"`)
Asigna rangos de VLANs para uso personal de usuarios/grupos.

## Data Source Relacionado

Use el data source `isard_network_interfaces` para buscar interfaces existentes:

```hcl
data "isard_network_interfaces" "bridges" {
  filter = {
    kind = "bridge"
  }
}

output "bridge_list" {
  value = data.isard_network_interfaces.bridges.interfaces
}
```
