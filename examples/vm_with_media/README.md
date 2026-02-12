# VMs con Medios (ISOs/Floppies)

Este ejemplo demuestra cómo adjuntar medios (ISOs, discos floppy) a máquinas virtuales en Isard VDI.

## Casos de Uso

Los medios adjuntos son útiles para:
- **Instalación de Sistemas Operativos**: Adjuntar ISOs de instalación
- **Drivers**: Añadir drivers VirtIO para Windows, drivers de hardware
- **Herramientas de Diagnóstico**: ISOs de rescate, GParted, SystemRescue
- **Software Adicional**: Herramientas que no están en el template base
- **Actualizaciones**: Paquetes de actualización en formato ISO

## Ejemplos Incluidos

### 1. VM con ISO Específico
```hcl
resource "isardvdi_vm" "windows_with_drivers" {
  name        = "Windows VM con Drivers"
  template_id = "windows-template-id"
  isos        = ["virtio-drivers-iso-id"]
}
```
Adjunta un ISO de drivers VirtIO a una VM de Windows.

### 2. VM con Múltiples ISOs
```hcl
resource "isardvdi_vm" "lab_vm" {
  name        = "VM de Laboratorio"
  template_id = "template-id"
  isos = [
    "ubuntu-22.04-iso-id",
    "gparted-live-iso-id",
    "system-rescue-iso-id"
  ]
}
```
Adjunta múltiples ISOs para diferentes propósitos.

### 3. VM con Data Source
```hcl
data "isardvdi_medias" "ubuntu_iso" {
  name_filter = "Ubuntu 22.04"
  kind        = "iso"
  status      = "Downloaded"
}

resource "isardvdi_vm" "ubuntu_vm" {
  name        = "Ubuntu Desktop"
  template_id = "template-id"
  isos = length(data.isardvdi_medias.ubuntu_iso.medias) > 0 ? [
    data.isardvdi_medias.ubuntu_iso.medias[0].id
  ] : []
}
```
Busca dinámicamente ISOs y los adjunta.

### 4. VM con Floppy
```hcl
resource "isardvdi_vm" "legacy_vm" {
  name        = "VM Legacy"
  template_id = "template-id"
  floppies    = ["driver-floppy-id"]
}
```
Para casos legacy que requieren disquetes.

## Cómo Obtener IDs de Medios

### Método 1: Data Source (Recomendado)
```hcl
data "isardvdi_medias" "my_iso" {
  name_filter = "Ubuntu"
  kind        = "iso"
  status      = "Downloaded"
}

# Usar el primer resultado
isos = [data.isardvdi_medias.my_iso.medias[0].id]
```

### Método 2: ID Directo
Si conoces el ID del medio:
```hcl
isos = ["a1b2c3d4-e5f6-7890-abcd-ef1234567890"]
```

### Método 3: Listar Todos
```hcl
data "isardvdi_medias" "all" {}

output "available_isos" {
  value = {
    for media in data.isardvdi_medias.all.medias :
    media.name => media.id
    if media.kind == "iso" && media.status == "Downloaded"
  }
}
```

## Orden de los Medios

Los ISOs se adjuntan en el orden especificado en la lista:
```hcl
isos = [
  "first-iso-id",   # Aparecerá como CD-ROM 0
  "second-iso-id",  # Aparecerá como CD-ROM 1
  "third-iso-id"    # Aparecerá como CD-ROM 2
]
```

## Limitaciones

1. **Solo al Crear**: Los medios solo se pueden especificar al crear la VM. Cambios posteriores requieren recrear la VM.

2. **Medios Descargados**: Asegúrate de que los medios tengan `status = "Downloaded"` antes de usarlos.

3. **Rendimiento**: Múltiples ISOs pueden afectar el rendimiento de arranque de la VM.

4. **Compatibilidad**: No todos los templates soportan medios adjuntos. Verifica la compatibilidad del template.

## Verificación de Medios Disponibles

Antes de crear VMs, verifica qué medios están disponibles:

```bash
terraform plan -target=data.isardvdi_medias.all
```

O usa outputs:
```hcl
data "isardvdi_medias" "all" {}

output "available_medias" {
  value = [
    for media in data.isardvdi_medias.all.medias : {
      name   = media.name
      id     = media.id
      status = media.status
      kind   = media.kind
    }
  ]
}
```

## Práctica Recomendada

1. **Buscar primero**: Usa data sources para buscar medios disponibles
2. **Verificar estado**: Asegúrate de que `status == "Downloaded"`
3. **Manejo de errores**: Usa expresiones condicionales para manejar medios no disponibles
4. **Documentación**: Comenta por qué cada ISO es necesario

### Ejemplo Completo con Buenas Prácticas

```hcl
# Buscar medios necesarios
data "isardvdi_medias" "os_iso" {
  name_filter = "Ubuntu 22.04"
  kind        = "iso"
  status      = "Downloaded"
}

data "isardvdi_medias" "drivers" {
  name_filter = "VirtIO"
  kind        = "iso"
  status      = "Downloaded"
}

# Verificar disponibilidad con check
check "medias_available" {
  assert {
    condition     = length(data.isardvdi_medias.os_iso.medias) > 0
    error_message = "Ubuntu 22.04 ISO no está disponible"
  }
  
  assert {
    condition     = length(data.isardvdi_medias.drivers.medias) > 0
    error_message = "VirtIO drivers ISO no está disponible"
  }
}

# Crear VM solo si los medios están disponibles
resource "isardvdi_vm" "production_vm" {
  name        = "Production VM"
  description = "VM con OS y drivers"
  template_id = "template-id"
  
  vcpus  = 4
  memory = 8
  
  # Adjuntar ISOs verificados
  isos = [
    data.isardvdi_medias.os_iso.medias[0].id,      # Sistema operativo
    data.isardvdi_medias.drivers.medias[0].id       # Drivers
  ]
  
  force_stop_on_destroy = true
}
```

## Ejecución

```bash
# Ver el plan
terraform plan

# Aplicar cambios
terraform apply

# Ver outputs
terraform output

# Destruir recursos
terraform destroy
```

## Troubleshooting

### Error: Media ID no encontrado
```
Error: Media not found
```
**Solución**: Verifica que el ID del medio existe y está descargado.

### Error: Template no soporta medios
```
Error: Template does not support media attachments
```
**Solución**: Usa un template diferente o consulta con el administrador de Isard VDI.

### ISOs no aparecen en la VM
**Solución**: Los ISOs solo se adjuntan en la creación. Si modificaste después, recrea la VM.

## Más Información

- Ver [isardvdi_medias data source](../../docs/data-sources/isardvdi_medias.md) para buscar medios
- Ver [isardvdi_media resource](../../docs/resources/isardvdi_media.md) para crear nuevos medios
- Ver [isardvdi_vm resource](../../docs/resources/isardvdi_vm.md) para documentación completa de VMs
