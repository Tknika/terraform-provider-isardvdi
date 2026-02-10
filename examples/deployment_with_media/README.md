# Deployments con Medios (ISOs/Floppies)

Este ejemplo demuestra cómo adjuntar medios (ISOs, discos floppy) a deployments en Isard VDI. Todos los desktops creados por el deployment tendrán los mismos medios adjuntos.

## Casos de Uso

Los medios en deployments son útiles para:
- **Formación/Training**: Todos los estudiantes tienen el mismo software/herramientas
- **Laboratorios**: Entornos de práctica con materiales pre-cargados
- **Instalación Masiva**: Desplegar múltiples VMs con el mismo ISO de instalación
- **Drivers Comunes**: Asegurar que todos tengan los drivers necesarios
- **Herramientas Compartidas**: Distribución de software a múltiples usuarios

## Ejemplos Incluidos

### 1. Deployment con ISO Específico
```hcl
resource "isard_deployment" "training_with_tools" {
  name         = "Training Environment"
  template_id  = "template-id"
  desktop_name = "Training Desktop"
  isos         = ["tools-iso-id"]
  
  allowed {
    roles = ["user"]
  }
}
```
Todos los desktops creados tendrán el ISO de herramientas adjunto.

### 2. Deployment con Múltiples ISOs
```hcl
resource "isard_deployment" "lab_deployment" {
  name         = "Lab Deployment"
  template_id  = "template-id"
  desktop_name = "Lab Desktop"
  
  isos = [
    "ubuntu-iso-id",
    "drivers-iso-id",
    "tools-iso-id"
  ]
  
  allowed {
    roles  = ["advanced", "user"]
    groups = ["students-group-id"]
  }
}
```
Cada desktop tendrá tres ISOs adjuntos.

### 3. Deployment con Data Source
```hcl
data "isard_medias" "installation_media" {
  name_filter = "Windows Server"
  kind        = "iso"
  status      = "Downloaded"
}

data "isard_medias" "virtio_drivers" {
  name_filter = "VirtIO"
  kind        = "iso"
  status      = "Downloaded"
}

resource "isard_deployment" "windows_deployment" {
  name         = "Windows Server Deployment"
  template_id  = "windows-template-id"
  desktop_name = "Windows Server"
  
  isos = concat(
    length(data.isard_medias.installation_media.medias) > 0 ? [data.isard_medias.installation_media.medias[0].id] : [],
    length(data.isard_medias.virtio_drivers.medias) > 0 ? [data.isard_medias.virtio_drivers.medias[0].id] : []
  )
  
  allowed {
    roles = ["admin", "advanced"]
  }
}
```
Busca y combina múltiples ISOs dinámicamente.

### 4. Deployment Condicional
```hcl
data "isard_medias" "required_iso" {
  name_filter = "Required Software"
  kind        = "iso"
  status      = "Downloaded"
}

resource "isard_deployment" "conditional_deployment" {
  count = length(data.isard_medias.required_iso.medias) > 0 ? 1 : 0
  
  name         = "Conditional Deployment"
  template_id  = "template-id"
  desktop_name = "Desktop"
  
  isos = [data.isard_medias.required_iso.medias[0].id]
  
  allowed {
    roles = ["user"]
  }
}
```
Solo crea el deployment si el ISO requerido está disponible.

## Propagación de Medios

Cuando adjuntas medios a un deployment:
1. **Creación**: Los ISOs se especifican en la configuración del deployment
2. **Propagación**: Cada desktop creado a partir del deployment hereda los mismos ISOs
3. **Usuarios**: Todos los usuarios autorizados ven los mismos medios en sus desktops

```
Deployment (isos = ["iso-1", "iso-2"])
    ├─> Desktop Usuario 1 (tiene iso-1, iso-2)
    ├─> Desktop Usuario 2 (tiene iso-1, iso-2)
    └─> Desktop Usuario 3 (tiene iso-1, iso-2)
```

## Diferencias con VMs Individuales

| Aspecto | Deployment | VM Individual |
|---------|-----------|---------------|
| Alcance | Múltiples usuarios | Un usuario |
| Medios | Compartidos por todos | Específicos del usuario |
| Actualización | Afecta nuevos desktops | Inmediato |
| Uso típico | Formación, labs | Desarrollo, producción |

## Gestión de Medios en Producción

### Estrategia 1: Medios Base + Específicos
```hcl
# Medios comunes para todos
locals {
  common_isos = [
    "drivers-iso-id",
    "tools-iso-id"
  ]
}

# Deployment básico
resource "isard_deployment" "basic" {
  name         = "Basic Deployment"
  template_id  = "template-id"
  desktop_name = "Basic Desktop"
  isos         = local.common_isos
  
  allowed {
    roles = ["user"]
  }
}

# Deployment avanzado (más ISOs)
resource "isard_deployment" "advanced" {
  name         = "Advanced Deployment"
  template_id  = "template-id"
  desktop_name = "Advanced Desktop"
  isos = concat(
    local.common_isos,
    ["advanced-tools-iso-id"]
  )
  
  allowed {
    roles = ["advanced"]
  }
}
```

### Estrategia 2: Variables para Flexibilidad
```hcl
variable "include_dev_tools" {
  description = "Incluir herramientas de desarrollo"
  type        = bool
  default     = false
}

data "isard_medias" "dev_tools" {
  count       = var.include_dev_tools ? 1 : 0
  name_filter = "Development Tools"
  kind        = "iso"
  status      = "Downloaded"
}

resource "isard_deployment" "flexible" {
  name         = "Flexible Deployment"
  template_id  = "template-id"
  desktop_name = "Desktop"
  
  isos = concat(
    ["base-iso-id"],
    var.include_dev_tools && length(data.isard_medias.dev_tools) > 0 ? [data.isard_medias.dev_tools[0].medias[0].id] : []
  )
  
  allowed {
    roles = ["user"]
  }
}
```

### Estrategia 3: Por Perfil de Usuario
```hcl
locals {
  user_profiles = {
    student = {
      isos = ["basic-iso-id"]
    }
    developer = {
      isos = ["basic-iso-id", "ide-iso-id", "tools-iso-id"]
    }
    tester = {
      isos = ["basic-iso-id", "testing-iso-id"]
    }
  }
}

resource "isard_deployment" "by_profile" {
  for_each = local.user_profiles
  
  name         = "${each.key} Deployment"
  template_id  = "template-id"
  desktop_name = "${each.key} Desktop"
  isos         = each.value.isos
  
  allowed {
    groups = ["${each.key}-group-id"]
  }
}
```

## Verificación Pre-Deployment

Antes de crear un deployment con medios, verifica la disponibilidad:

```hcl
# Check para validar medios
check "medias_ready" {
  data "isard_medias" "required" {
    name_filter = "Required ISO"
    kind        = "iso"
    status      = "Downloaded"
  }
  
  assert {
    condition     = length(data.isard_medias.required.medias) > 0
    error_message = "Required ISO no está disponible para el deployment"
  }
}
```

## Limitaciones y Consideraciones

### Limitaciones
1. **Solo en Creación**: Los medios se especifican al crear el deployment, no se pueden modificar después
2. **Todos o Ninguno**: Todos los desktops del deployment tienen los mismos medios
3. **Estado del Medio**: Los medios deben estar en estado `Downloaded`

### Consideraciones
1. **Rendimiento**: Múltiples ISOs pueden ralentizar el arranque de los desktops
2. **Espacio**: Cada desktop no duplica los medios, se referencian
3. **Permisos**: Los usuarios deben tener permisos para usar los medios especificados
4. **Template**: El template debe ser compatible con medios adjuntos

## Mejores Prácticas

1. **Usa Data Sources**: Busca medios dinámicamente en lugar de IDs hardcoded
2. **Verifica Estado**: Siempre filtra por `status = "Downloaded"`
3. **Documentación**: Documenta qué hace cada ISO y por qué es necesario
4. **Versionado**: Mantén versiones de ISOs para rollbacks si es necesario
5. **Testing**: Prueba con un deployment pequeño antes de desplegar a producción

## Ejecución

```bash
# Ver el plan
terraform plan

# Ver solo los deployments con medios
terraform plan -target=isard_deployment.with_media

# Aplicar
terraform apply

# Ver outputs
terraform output

# Destruir
terraform destroy
```

## Troubleshooting

### Problema: "Media not found"
**Causa**: El ID del medio no existe o no está descargado  
**Solución**: Verifica el estado del medio con `data "isard_medias"`

### Problema: Desktops sin ISOs
**Causa**: Los desktops se crearon antes de que el medio estuviera disponible  
**Solución**: Los medios solo se adjuntan en la creación del desktop, no retroactivamente

### Problema: "Permission denied"
**Causa**: Los usuarios del deployment no tienen permisos para usar los medios  
**Solución**: Verifica el campo `allowed` en el recurso `isard_media`

## Siguiente Paso

- Ver [isard_medias data source](../../docs/data-sources/isard_medias.md)
- Ver [isard_media resource](../../docs/resources/isard_media.md)
- Ver [isard_deployment resource](../../docs/resources/deployment.md)
