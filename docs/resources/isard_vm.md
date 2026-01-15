# Resource: isard_vm

Gestiona un desktop persistente en Isard VDI.

## Ejemplo de Uso

### Ejemplo Básico

```hcl
resource "isard_vm" "ejemplo" {
  name        = "mi-desktop"
  description = "Desktop de desarrollo"
  template_id = "60aca659-d627-4f30-a894-52da23c18212"
}
```

### Con Data Source

```hcl
data "isard_templates" "ubuntu" {
  name_filter = "Ubuntu"
}

resource "isard_vm" "desarrollo" {
  name        = "desktop-desarrollo"
  description = "Desktop Ubuntu para desarrollo"
  template_id = data.isard_templates.ubuntu.templates[0].id
}
```

### Múltiples Desktops

```hcl
resource "isard_vm" "equipo" {
  count = 3
  
  name        = "desktop-dev-${count.index + 1}"
  description = "Desktop para desarrollador ${count.index + 1}"
  template_id = data.isard_templates.ubuntu.templates[0].id
}
```

### Con Hardware Personalizado

```hcl
resource "isard_vm" "potente" {
  name        = "desktop-potente"
  description = "Desktop con recursos aumentados"
  template_id = data.isard_templates.ubuntu.templates[0].id
  vcpus       = 8
  memory      = 16.0
}
```

### Con Interfaces de Red Personalizadas

```hcl
# Buscar interfaces disponibles
data "isard_network_interfaces" "all" {}

# Crear VM con interfaces específicas
resource "isard_vm" "con_red" {
  name        = "desktop-con-red-custom"
  description = "Desktop con interfaces personalizadas"
  template_id = data.isard_templates.ubuntu.templates[0].id
  
  # Importante: Si el template tiene RDP viewers, incluir wireguard
  interfaces = [
    "wireguard",
    data.isard_network_interfaces.all.interfaces[0].id
  ]
}
```

## Argumentos

Los siguientes argumentos son soportados:

### Requeridos

- `name` - (Requerido) Nombre del desktop. Debe ser único.
- `template_id` - (Requerido) ID del template a usar como base para el desktop.

### Opcionales

- `description` - (Opcional) Descripción del desktop.
- `vcpus` - (Opcional) Número de CPUs virtuales. Si no se especifica, usa el valor del template.
- `memory` - (Opcional) Memoria RAM en GB. Si no se especifica, usa el valor del template.
- `interfaces` - (Opcional) Lista de IDs de interfaces de red a usar. Si no se especifica, usa las interfaces del template.

## Atributos Exportados

Además de los argumentos anteriores, se exportan los siguientes atributos:

- `id` - ID único del desktop en Isard VDI.
- `vcpus` - Número de CPUs virtuales asignadas al desktop (computed).
- `memory` - Memoria RAM asignada al desktop en GB (computed).

## Import

Los desktops pueden ser importados usando su ID:

```bash
terraform import isard_vm.ejemplo a1b2c3d4-e5f6-7890-abcd-ef1234567890
```

## Ciclo de Vida

### Create

Al crear un desktop:
1. Se valida que el `template_id` sea válido
2. Se crea un desktop persistente usando `POST /api/v3/persistent_desktop`
3. Se obtiene el ID del desktop creado
4. Se leen los valores de hardware asignados por el servidor

### Read

Al leer un desktop:
1. Se obtiene la información desde `GET /api/v3/domain/info/{id}`
2. Se actualizan todos los atributos computados

### Update

Actualmente no está implementado. Cualquier cambio en los atributos requerirá recrear el recurso.

### Delete

Al eliminar un desktop:
1. Se elimina permanentemente usando `DELETE /api/v3/desktop/{id}/true`
2. El parámetro `true` indica eliminación permanente (no enviar a papelera)

## Notas Importantes

### Hardware Personalizado

Actualmente, Isard VDI no permite especificar valores personalizados de `vcpus` y `memory` al crear desktops. Los valores se determinan por:
- Configuración del template base
- Cuotas del usuario/grupo
- Políticas del sistema

### Estados del Desktop

El desktop puede estar en varios estados:
- `Stopped` - Detenido
- `Started` - Iniciado
- `Creating` - En creación
- Otros estados según la configuración de Isard VDI

Terraform solo gestiona la existencia del desktop, no su estado de ejecución.

### Dependencias

Si usas un data source para obtener el `template_id`, Terraform gestionará automáticamente las dependencias:

```hcl
data "isard_templates" "ubuntu" {
  name_filter = "Ubuntu"
}

resource "isard_vm" "mi_vm" {
  name        = "mi-desktop"
  template_id = data.isard_templates.ubuntu.templates[0].id
  # Terraform esperará a que el data source se resuelva antes de crear
}
```

## Limitaciones Conocidas

1. No se puede actualizar un desktop existente (solo crear/eliminar)
2. No se puede controlar el estado de ejecución del desktop
3. No se pueden personalizar valores de hardware en la creación
4. No se puede especificar hardware adicional (discos, interfaces de red, etc.)

## Ejemplos Adicionales

### Desktop con Nombre Dinámico

```hcl
variable "entorno" {
  type    = string
  default = "desarrollo"
}

resource "isard_vm" "dinamico" {
  name        = "desktop-${var.entorno}-${formatdate("YYYYMMDD", timestamp())}"
  description = "Desktop de ${var.entorno}"
  template_id = data.isard_templates.ubuntu.templates[0].id
}
```

### Desktop con Lifecycle

```hcl
resource "isard_vm" "importante" {
  name        = "desktop-produccion"
  description = "Desktop crítico - no eliminar accidentalmente"
  template_id = data.isard_templates.ubuntu.templates[0].id
  
  lifecycle {
    prevent_destroy = true
  }
}
```
