---
page_title: "isard_vm Resource - terraform-provider-isard"
subcategory: ""
description: |-
  Manages a persistent desktop in Isard VDI.
---

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

### Con Force Stop on Destroy

```hcl
resource "isard_vm" "produccion" {
  name        = "desktop-produccion"
  description = "Desktop de producción con stop seguro"
  template_id = data.isard_templates.ubuntu.templates[0].id
  
  # Detener la VM antes de eliminarla para evitar pérdida de datos
  force_stop_on_destroy = true
}
```

### Con Medios ISOs Adjuntos

```hcl
# Buscar ISOs disponibles
data "isard_medias" "ubuntu_iso" {
  name_filter = "Ubuntu"
  kind        = "iso"
  status      = "Downloaded"
}

resource "isard_vm" "con_iso" {
  name        = "desktop-con-iso"
  description = "Desktop con ISO adjunto"
  template_id = data.isard_templates.ubuntu.templates[0].id
  
  # Adjuntar ISO de Ubuntu
  isos = length(data.isard_medias.ubuntu_iso.medias) > 0 ? [
    data.isard_medias.ubuntu_iso.medias[0].id
  ] : []
}
```

### Con Múltiples ISOs

```hcl
resource "isard_vm" "con_multiples_isos" {
  name        = "desktop-herramientas"
  description = "Desktop con múltiples ISOs de herramientas"
  template_id = data.isard_templates.ubuntu.templates[0].id
  
  vcpus  = 4
  memory = 8
  
  # Adjuntar múltiples ISOs
  isos = [
    "iso-instalacion-id",
    "iso-drivers-id",
    "iso-herramientas-id"
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
- `isos` - (Opcional) Lista de IDs de medios ISO a adjuntar al desktop. Estos aparecerán como unidades de CD/DVD en la VM.
- `floppies` - (Opcional) Lista de IDs de medios floppy a adjuntar al desktop. Raramente usado en VMs modernas.
- `force_stop_on_destroy` - (Opcional) Si es `true`, fuerza la parada de la máquina virtual antes de eliminarla usando el endpoint de administración (parada forzada) y espera hasta 10 segundos. Por defecto: `false`. La parada forzada garantiza que la VM se detenga inmediatamente, incluso si no responde, previniendo largos tiempos de espera durante la destrucción.

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
1. Si `force_stop_on_destroy` es `true`:
   - Se fuerza la parada de la VM usando `POST /api/v3/admin/multiple_actions` con action="stopping"
   - Se verifica inmediatamente si ya está detenida (retorno instantáneo si lo está)
   - Se espera hasta 10 segundos a que la VM se detenga completamente
   - La parada forzada cambia el estado a "Stopping" (vs "Shutting-down" de soft stop)
   - Si el stop o la espera fallan, se muestra una advertencia pero se continúa
2. Se elimina permanentemente usando `DELETE /api/v3/desktop/{id}/true`
3. El parámetro `true` indica eliminación permanente (no enviar a papelera)

## Notas Importantes

### Force Stop on Destroy

Si `force_stop_on_destroy` está habilitado, Terraform:
1. Verificará inmediatamente si la VM ya está detenida (retorno instantáneo si lo está)
2. Forzará la parada de la VM usando el endpoint de administración (parada inmediata)
3. Esperará hasta 10 segundos a que la VM se detenga completamente
4. Procederá con la eliminación del desktop

La parada forzada utiliza el endpoint `/api/v3/admin/multiple_actions` que cambia el estado a "Stopping" en lugar de "Shutting-down", garantizando una detención inmediata. Esto resulta en tiempos de destrucción mucho más rápidos (típicamente 2-5 segundos) comparado con el soft stop tradicional. Si el stop o la espera fallan, Terraform mostrará una advertencia pero continuará con la eliminación.

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
