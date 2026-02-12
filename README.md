# Terraform Provider para Isard VDI

Provider de Terraform para gestionar recursos en Isard VDI a través de su API v3.

## Características

### Recursos

- ✅ **isardvdi_vm** - Creación, lectura y eliminación de desktops persistentes con soporte para:
  - Hardware personalizado (vCPUs, memoria)
  - Interfaces de red personalizadas
  - Medios ISO y floppy adjuntos
  - Parada forzada optimizada antes de destruir (force_stop_on_destroy con 10s timeout)
- ✅ **isardvdi_deployment** - Gestión de deployments para crear múltiples desktops para usuarios/grupos con:
  - Hardware personalizado
  - Medios ISO y floppy adjuntos
  - Manejo automático de error 428 (VMs en ejecución) con retry inteligente
- ✅ **isardvdi_media** - Gestión de medios (ISOs y floppies) para adjuntar a VMs
- ✅ **isardvdi_network** - Gestión de redes virtuales de usuario
- ✅ **isardvdi_network_interface** - Gestión de interfaces de red del sistema (requiere admin)
- ✅ **isardvdi_qos_net** - Gestión de perfiles QoS de red (requiere admin)

### Data Sources

- ✅ **isardvdi_templates** - Listado de templates disponibles con filtrado por nombre
- ✅ **isardvdi_network_interfaces** - Consulta de interfaces de red del sistema con filtros avanzados
- ✅ **isardvdi_groups** - Consulta de grupos del sistema con filtrado por nombre y categoría
- ✅ **isardvdi_users** - Consulta de usuarios del sistema con múltiples filtros (nombre, username, email, categoría, grupo, rol)
- ✅ **isardvdi_medias** - Consulta de medios disponibles con filtros avanzados (nombre, tipo, estado, categoría, grupo, usuario)

### Autenticación

- ✅ Soporte para autenticación mediante token JWT
- ✅ Soporte para autenticación mediante formulario (usuario/contraseña)
- ✅ Configuración SSL flexible para desarrollo y producción

## Requisitos

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.24 (para desarrollo)
- Instancia de Isard VDI con API v3 activa

## Instalación

### Desarrollo Local

1. Clonar el repositorio:
```bash
git clone https://github.com/prubioamzubiri/isard-terraform-provider.git
cd isard-terraform-provider
```

2. Compilar el provider:
```bash
go build -o terraform-provider-isard
```

3. Configurar el override local en `~/.terraformrc`:
```hcl
provider_installation {
  dev_overrides {
    "tknika/isardvdi" = "/home/tu-usuario/source/isard-terraform-provider"
  }
  direct {}
}
```

### Instalación desde Registry (futuro)

```hcl
terraform {
  required_providers {
    isardvdi = {
      source  = "registry.terraform.io/tknika/isardvdi"
      version = "~> 1.0"
    }
  }
}
```

## Uso Básico

```hcl
# Configurar el provider
provider "isardvdi" {
  endpoint     = "localhost"
  auth_method  = "form"
  cathegory_id = "default"
  username     = "admin"
  password     = "IsardVDI"
  # ssl_verification = false  # Descomenta para desarrollo con certificados autofirmados
}

# Obtener templates disponibles
data "isardvdi_templates" "ubuntu" {
  name_filter = "Ubuntu"
}

# Obtener grupos por nombre
data "isardvdi_groups" "desarrollo" {
  name_filter = "Desarrollo"
}

# Crear un desktop persistente
resource "isardvdi_vm" "mi_desktop" {
  name        = "mi-desktop-terraform"
  description = "Desktop creado con Terraform"
  template_id = data.isardvdi_templates.ubuntu.templates[0].id
}

# Crear un deployment para un equipo
resource "isardvdi_deployment" "equipo_dev" {
  name         = "Deployment Equipo Dev"
  description  = "Desktops para el equipo de desarrollo"
  template_id  = data.isardvdi_templates.ubuntu.templates[0].id
  desktop_name = "Desktop Dev"
  visible      = false
  
  vcpus  = 4
  memory = 8.0

  allowed {
    groups = [data.isardvdi_groups.desarrollo.groups[0].id]
  }
}

# Crear una red virtual
resource "isardvdi_network" "mi_red" {
  name        = "Red de Desarrollo"
  description = "Red virtual para desarrollo"
  model       = "virtio"
  qos_id      = "unlimited"
}

# Crear interfaz de red del sistema (requiere admin)
resource "isardvdi_network_interface" "bridge_custom" {
  id          = "bridge-custom"
  name        = "Bridge Personalizado"
  description = "Bridge para entorno custom"
  net         = "br-custom"
  kind        = "bridge"
  model       = "virtio"
  qos_id      = "unlimited"
  
  # Hacer visible para todos los usuarios
  allowed {
    roles      = []
    categories = []
    groups     = []
    users      = []
  }
}

# Crear VM con interfaces personalizadas
resource "isardvdi_vm" "vm_custom" {
  name        = "vm-con-red-custom"
  description = "VM con interfaces personalizadas"
  template_id = data.isardvdi_templates.ubuntu.templates[0].id
  
  interfaces = [
    "wireguard",  # Requerido para RDP
    isardvdi_network_interface.bridge_custom.id
  ]
}
```

## Documentación

### Provider

- [Configuración del Provider](docs/index.md)

### Recursos

- [Resource: isardvdi_vm](docs/resources/isardvdi_vm.md) - Gestión de VMs/desktops
- [Resource: isardvdi_deployment](docs/resources/isardvdi_deployment.md) - Gestión de deployments
- [Resource: isardvdi_media](docs/resources/isardvdi_media.md) - Gestión de medios (ISOs y floppies)
- [Resource: isardvdi_network](docs/resources/isardvdi_network.md) - Redes virtuales de usuario
- [Resource: isardvdi_network_interface](docs/resources/isardvdi_network_interface.md) - Interfaces de red del sistema
- [Resource: isardvdi_qos_net](docs/resources/isardvdi_qos_net.md) - Perfiles QoS de red

### Data Sources

- [Data Source: isardvdi_templates](docs/data-sources/isardvdi_templates.md) - Consulta de templates
- [Data Source: isardvdi_users](docs/data-sources/isardvdi_users.md) - Consulta de usuarios con filtros avanzados
- [Data Source: isardvdi_medias](docs/data-sources/isardvdi_medias.md) - Consulta de medios (ISOs y floppies)
- [Data Source: isardvdi_network_interfaces](docs/data-sources/isardvdi_network_interfaces.md) - Consulta de interfaces
- [Data Source: isardvdi_groups](docs/data-sources/isardvdi_groups.md) - Consulta de grupos

## Ejemplos

Consulta el directorio [examples/](examples/) para ver ejemplos completos de uso.

## Desarrollo

### Compilar

```bash
go build -o terraform-provider-isard
```

### Ejecutar Tests

```bash
go test ./...
```

### Depuración

Para habilitar logs detallados:

```bash
export TF_LOG=DEBUG
terraform plan
```

## Contribuir

Las contribuciones son bienvenidas. Por favor:

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/mi-feature`)
3. Commit tus cambios (`git commit -am 'Agregar nueva característica'`)
4. Push a la rama (`git push origin feature/mi-feature`)
5. Abre un Pull Request

## Licencia

[Especificar licencia]

## Soporte

Para reportar bugs o solicitar features, abre un issue en GitHub.
