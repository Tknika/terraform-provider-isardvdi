---
page_title: "isard_deployment Resource - terraform-provider-isard"
subcategory: ""
description: |-
  Manages a deployment in Isard VDI.
---

# isard_deployment (Resource)

Gestiona un deployment en Isard VDI. Los deployments permiten crear múltiples desktops a partir de una plantilla para diferentes usuarios, grupos o categorías.

## Ejemplo de Uso

```terraform
# Crear un deployment básico
resource "isard_deployment" "example" {
  name         = "Deployment de Prueba"
  description  = "Deployment para el equipo de desarrollo"
  template_id  = "template-uuid-123"
  desktop_name = "Desktop Dev"
  visible      = false

  allowed {
    groups = ["group-uuid-1", "group-uuid-2"]
  }
}

# Deployment con hardware personalizado
resource "isard_deployment" "custom_hardware" {
  name         = "Deployment Alto Rendimiento"
  description  = "Deployment con recursos mejorados"
  template_id  = "template-uuid-456"
  desktop_name = "Desktop Performance"
  visible      = true
  
  vcpus      = 4
  memory     = 8.0
  interfaces = ["interface-uuid-1"]

  allowed {
    users = ["user-uuid-1", "user-uuid-2"]
  }
}

# Deployment para una categoría completa
resource "isard_deployment" "category_deployment" {
  name         = "Deployment Estudiantes"
  description  = "Desktops para todos los estudiantes"
  template_id  = "template-uuid-789"
  desktop_name = "Desktop Estudiante"
  visible      = true

  allowed {
    categories = ["category-students-uuid"]
  }
  
  user_permissions = ["viewer", "desktop_start"]
}

# Deployment con múltiples interfaces de red
resource "isard_deployment" "network_deployment" {
  name         = "Deployment Red Avanzada"
  description  = "Deployment con configuración de red compleja"
  template_id  = "template-uuid-abc"
  desktop_name = "Desktop Network"
  
  interfaces = [
    "interface-uuid-1",
    "interface-uuid-2",
    "interface-uuid-3"
  ]

  allowed {
    groups = ["network-team-uuid"]
    users  = ["admin-user-uuid"]
  }
}

# Deployment con viewers específicos
resource "isard_deployment" "custom_viewers" {
  name         = "Deployment con Viewers Personalizados"
  description  = "Deployment con solo viewers web"
  template_id  = "template-uuid-def"
  desktop_name = "Desktop Web Access"
  visible      = true
  
  # Solo habilitar viewers basados en navegador
  viewers = ["browser_vnc", "browser_rdp"]

  allowed {
    groups = ["remote-team-uuid"]
  }
}

# Deployment con force_stop_on_destroy habilitado
resource "isard_deployment" "safe_destroy" {
  name         = "Deployment con Stop Automático"
  description  = "Las VMs se detendrán antes de eliminar"
  template_id  = "template-uuid-xyz"
  desktop_name = "Desktop Safe"
  visible      = true
  
  # Detener todas las VMs antes de eliminar el deployment
  force_stop_on_destroy = true

  allowed {
    groups = ["production-team-uuid"]
  }
}

# Deployment con todos los viewers disponibles
resource "isard_deployment" "all_viewers" {
  name         = "Deployment Acceso Completo"
  description  = "Deployment con todos los métodos de acceso"
  template_id  = "template-uuid-ghi"
  desktop_name = "Desktop Full Access"
  
  viewers = [
    "browser_rdp",
    "browser_vnc",
    "file_rdpgw",
    "file_rdpvpn",
    "file_spice"
  ]

  allowed {
    users = ["power-user-uuid"]
  }
}
```

## Esquema de Argumentos

### Requeridos

- `name` (String) Nombre del deployment. Mínimo 4 caracteres, máximo 50.
- `template_id` (String) ID de la plantilla a utilizar para crear los desktops del deployment. **Nota:** Cambiar este valor forzará la recreación del deployment.
- `desktop_name` (String) Nombre base para los desktops creados en el deployment.
- `allowed` (Block, Requerido) Configuración de usuarios, grupos, categorías y roles permitidos para acceder a este deployment. Ver [Allowed](#allowed) más abajo.

### Opcionales

- `description` (String) Descripción del deployment. Máximo 255 caracteres.
- `visible` (Boolean) Si los desktops del deployment son visibles para los usuarios. Por defecto: `false`.
- `vcpus` (Number) Número de CPUs virtuales para los desktops. Si no se especifica, usa el valor del template.
- `memory` (Number) Memoria RAM en GB para los desktops. Si no se especifica, usa el valor del template.
- `interfaces` (List of String) Lista de IDs de interfaces de red a utilizar. Si no se especifica, usa las del template.
- `viewers` (List of String) Lista de viewers habilitados para los desktops. Si no se especifica, usa los viewers del template. Valores disponibles:
  - `browser_rdp` - Visor RDP en el navegador
  - `browser_vnc` - Visor VNC en el navegador (noVNC)
  - `file_rdpgw` - Archivo RDP con gateway
  - `file_rdpvpn` - Archivo RDP con VPN
  - `file_spice` - Visor SPICE (archivo de configuración)
- `user_permissions` (List of String) Lista de permisos de usuario para el deployment.
- `force_stop_on_destroy` (Boolean) Si es `true`, detiene todas las máquinas virtuales del deployment antes de eliminarlo y espera hasta 120 segundos a que se detengan completamente. Por defecto: `false`. Esto es útil para asegurar que todas las VMs se apaguen de manera ordenada antes de la eliminación del deployment y prevenir errores de eliminación causados por VMs en ejecución.

### Atributos de Solo Lectura

- `id` (String) Identificador único del deployment.

## Nested Schema para `allowed`

El bloque `allowed` configura qué usuarios, grupos, categorías o roles tienen acceso al deployment:

### Opcionales

- `roles` (List of String) Lista de roles permitidos (por ejemplo: `["admin", "manager"]`)
- `categories` (List of String) Lista de IDs de categorías permitidas
- `groups` (List of String) Lista de IDs de grupos permitidos
- `users` (List of String) Lista de IDs de usuarios permitidos

**Nota:** Al menos uno de estos campos debe especificarse en el bloque `allowed`.

## Viewers Disponibles

El parámetro `viewers` permite controlar qué métodos de visualización están disponibles para los desktops del deployment. Si no se especifica, se utilizarán los viewers configurados en el template.

### Tipos de Viewers

- **`browser_rdp`** - Visor RDP basado en navegador web. Permite conectarse a desktops Windows sin necesidad de instalar software adicional.
- **`browser_vnc`** - Visor VNC en el navegador (noVNC). Proporciona acceso gráfico multiplataforma desde el navegador.
- **`file_rdpgw`** - Archivo de configuración RDP con gateway. Descarga un archivo .rdp configurado para conectarse a través de un gateway RDP.
- **`file_rdpvpn`** - Archivo de configuración RDP con VPN. Descarga un archivo .rdp configurado para conectarse a través de VPN.
- **`file_spice`** - Visor SPICE (Simple Protocol for Independent Computing Environments). Proporciona alto rendimiento para desktops Linux/KVM.

### Ejemplos de Configuración

```terraform
# Solo viewers web (sin instalación de cliente)
viewers = ["browser_vnc", "browser_rdp"]

# Solo SPICE para máximo rendimiento
viewers = ["file_spice"]

# RDP con múltiples opciones de conexión
viewers = ["browser_rdp", "file_rdpgw", "file_rdpvpn"]

# Todos los viewers disponibles
viewers = ["browser_rdp", "browser_vnc", "file_rdpgw", "file_rdpvpn", "file_spice"]
```

**Recomendaciones:**
- Para usuarios remotos sin VPN, usar `browser_vnc` o `browser_rdp`
- Para máximo rendimiento en red local, usar `file_spice`
- Para compatibilidad con clientes RDP nativos, incluir `file_rdpgw` o `file_rdpvpn`

## Importación

Los deployments pueden ser importados usando su ID:

```bash
terraform import isard_deployment.example deployment-uuid-123
```

## Notas Adicionales

- **Desktops Automáticos:** Al crear un deployment, Isard VDI creará automáticamente un desktop para cada usuario que coincida con los criterios especificados en `allowed`.
- **Hardware:** Si especificas `vcpus`, `memory` o `interfaces`, estos valores sobrescriben los del template para todos los desktops del deployment.
- **Visibilidad:** El atributo `visible` controla si los desktops son visibles inmediatamente para los usuarios o si están ocultos hasta que sean habilitados.
- **Eliminación:** Al eliminar un deployment, todos los desktops asociados también serán eliminados permanentemente.
- **Force Stop on Destroy:** Si `force_stop_on_destroy` está habilitado, Terraform:
  1. Detendrá todas las VMs del deployment
  2. Esperará hasta 120 segundos a que todas las VMs se detengan completamente
  3. Procederá con la eliminación del deployment
  
  Esto asegura un apagado ordenado y previene errores de eliminación causados por VMs en ejecución. Si el stop o la espera fallan, Terraform mostrará una advertencia pero continuará con la eliminación.
- **Actualización:** Algunos cambios en el deployment pueden requerir que los desktops estén detenidos. Terraform te informará si esto es necesario.
