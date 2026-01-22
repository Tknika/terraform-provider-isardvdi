---
page_title: "isard_deployment Resource - terraform-provider-isard"
subcategory: ""
description: |-
  Gestiona un deployment en Isard VDI
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
- `user_permissions` (List of String) Lista de permisos de usuario para el deployment.

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
- **Actualización:** Algunos cambios en el deployment pueden requerir que los desktops estén detenidos. Terraform te informará si esto es necesario.
